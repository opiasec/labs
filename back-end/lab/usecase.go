package lab

import (
	"appseclabsplataform/database"
	labcluster "appseclabsplataform/services/labCluster"
	labide "appseclabsplataform/services/labIDE"
	labsession "appseclabsplataform/services/labSession"
	"errors"
	"log"
	"log/slog"
	"time"
)

type LabUsecase struct {
	LabCluster *labcluster.LabClusterService
	Database   *database.Database
	LabSession *labsession.LabSessionService
	LabIDE     *labide.LabIDEService
}

func NewLabUsecase(labCluster *labcluster.LabClusterService, database *database.Database, labSession *labsession.LabSessionService, labIDE *labide.LabIDEService) *LabUsecase {
	return &LabUsecase{
		LabCluster: labCluster,
		Database:   database,
		LabSession: labSession,
		LabIDE:     labIDE,
	}
}

func (u *LabUsecase) CreateLab(labSlug string, userID string) (*CreateLabResponse, error) {

	lab, err := u.Database.GetLabDefinitionBySlug(labSlug)
	if err != nil {
		slog.Error("Failed to get lab details from database", "error", err)
		return nil, err
	}
	createLabResponse, err := u.LabCluster.CreateLab(labSlug)
	if err != nil {
		slog.Error("Failed to create lab", "error", err)
		return nil, err
	}

	slog.Info("Lab created successfully", "namespace", createLabResponse.Namespace, "labID", lab.ID)

	err = u.Database.CreateLabAttempt(createLabResponse.Namespace, lab.ID, userID)
	if err != nil {
		slog.Error("Failed to create lab attempt in database", "error", err)
		return nil, err
	}

	err = u.LabSession.SetLabSession(&labsession.LabSession{
		LabAttemptID: lab.ID,
		LabSlug:      labSlug,
		Namespace:    createLabResponse.Namespace,
		UserID:       userID,
		IDEURL:       createLabResponse.IDEURL,
		Password:     createLabResponse.LabPassword,
		Status:       "running",
	}, 1*time.Hour)

	if err != nil {
		slog.Error("Failed to set lab session in Redis", "error", err)
		return nil, err
	}

	return &CreateLabResponse{
		CreateLabResponse: labcluster.CreateLabResponse(*createLabResponse),
		ExpiresAt:         time.Now().Add(1 * time.Hour).UTC(),
	}, nil
}

func (u *LabUsecase) GetLabStatus(namespace string, userID string) (*GetLabResponse, error) {
	labSession, err := u.LabSession.GetLabSession(namespace)
	if err != nil {
		slog.Error("Failed to get lab session from Redis", "error", err)
		return nil, err
	}
	if labSession == nil {
		slog.Error("Lab session not found in Redis", "namespace", namespace)
		return nil, errors.New("lab session not found")
	}

	if labSession.UserID != userID {
		slog.Error("User ID does not match lab session user ID", "expected", labSession.UserID, "got", userID)
		return nil, errors.New("unauthorized")
	}

	labStatus, err := u.LabIDE.GetStatus(namespace)
	if err != nil {
		slog.Error("Failed to get lab status from labIDE", "error", err)
		return nil, err
	}
	return &GetLabResponse{
		labide.GetLabResponse(labStatus),
	}, nil
}

func (u *LabUsecase) FinishLab(namespace string, userID string) (*FinishLabResponse, error) {
	labSession, err := u.LabSession.GetLabSession(namespace)
	if err != nil {
		slog.Error("Failed to get lab session", "error", err)
		return nil, err
	}

	if labSession.UserID != userID {
		slog.Error("User ID does not match lab session user ID", "expected", labSession.UserID, "got", userID)
		return nil, errors.New("unauthorized")
	}

	lab, err := u.LabCluster.FinishLab(namespace, labSession.LabSlug)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// TODO: Add Score logic here - for now, we just finish the lab
	// Maybe Score will come from lab cluster as webhook or something like that
	err = u.Database.FinishLabAttempt(namespace, 0)
	if err != nil {
		slog.Error("Failed to finish lab attempt in database", "error", err)
		return nil, err
	}

	err = u.LabSession.DeleteLabSession(namespace)
	if err != nil {
		slog.Error("Failed to delete lab session from Redis", "error", err)
		return nil, err
	}

	return &FinishLabResponse{
		labcluster.FinishLabResponse(*lab),
	}, nil
}

func (u *LabUsecase) LeaveLab(namespace string, userID string) (*database.LabAttempt, error) {
	labSession, err := u.LabSession.GetLabSession(namespace)
	if err != nil {
		slog.Error("Failed to get lab session", "error", err)
		return nil, err
	}

	if labSession.UserID != userID {
		slog.Error("User ID does not match lab session user ID", "expected", labSession.UserID, "got", userID)
		return nil, errors.New("unauthorized")
	}

	labAttempt, err := u.Database.LeaveLabAttempt(namespace)
	if err != nil {
		slog.Error("Failed to leave lab attempt", "error", err)
		return nil, err
	}

	err = u.LabSession.DeleteLabSession(namespace)
	if err != nil {
		slog.Error("Failed to delete lab session from Redis", "error", err)
		return nil, err
	}
	err = u.LabCluster.DeleteLabSession(namespace)
	if err != nil {
		slog.Error("Failed to delete lab from lab cluster", "error", err)
		return nil, err
	}
	return labAttempt, nil

}

func (u *LabUsecase) GetLabResult(namespace string, userID string) (*GetLabResultResponse, error) {
	labAttempt, err := u.Database.GetLabAttempt(namespace)
	if err != nil {
		slog.Error("Failed to get lab attempt from database", "error", err)
		return nil, err
	}

	if labAttempt.ExternalUserID != userID {
		slog.Error("User ID does not match lab attempt user ID", "expected", labAttempt.ExternalUserID, "got", userID)
		return nil, errors.New("unauthorized")
	}

	labClusterResult, err := u.LabCluster.GetLabResult(namespace)
	if err != nil {
		slog.Error("Failed to get lab result from lab cluster", "error", err)
		return nil, nil
	}
	slog.Info("Lab result retrieved successfully", "namespace", namespace, "labSlug", labAttempt.Lab.Slug)
	return &GetLabResultResponse{
		Namespace:       labAttempt.Namespace,
		LabSlug:         labAttempt.Lab.Slug,
		Status:          labAttempt.Status.Name,
		StartedAt:       labClusterResult.StartedAt,
		FinishedAt:      labClusterResult.FinishedAt,
		LabFinishResult: labClusterResult.FinishResult,
		Rating:          labAttempt.Rating,
		UserFeedback:    labAttempt.Feedback,
	}, nil
}

func (u *LabUsecase) GetAllLabsByUserAndStatus(status string, userID string) (*[]GetAllLabsByUserAndStatusResponse, error) {
	labsAttempts, err := u.Database.GetLabAttemptsByUser(userID)
	if err != nil {
		slog.Error("Failed to get lab attempts from database", "error", err)
		return nil, err
	}

	var response []GetAllLabsByUserAndStatusResponse
	for _, lab := range labsAttempts {
		response = append(response, GetAllLabsByUserAndStatusResponse{
			Namespace:    lab.Namespace,
			LabSlug:      lab.Lab.Slug,
			StartedAt:    lab.StartedAt,
			FinishedAt:   lab.FinishedAt,
			Score:        lab.Score,
			Status:       lab.Status.Name,
			UserFeedback: lab.Feedback,
			Rating:       lab.Rating,
		})
	}

	return &response, nil
}

func (u *LabUsecase) SendFeedback(namespace string, rating int, feedback string, userID string) error {
	labAttempt, err := u.Database.GetLabAttempt(namespace)
	if err != nil {
		slog.Error("Failed to get lab attempt from database", "error", err)
		return err
	}

	if labAttempt.ExternalUserID != userID {
		slog.Error("User ID does not match lab attempt user ID", "expected", labAttempt.ExternalUserID, "got", userID)
		return errors.New("unauthorized")
	}

	if err := u.Database.SetLabFeedback(namespace, rating, feedback); err != nil {
		slog.Error("Failed to set lab feedback in database", "error", err)
		return err
	}

	return nil
}
