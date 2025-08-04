package admin

import (
	"appseclabsplataform/database"
	labcluster "appseclabsplataform/services/labCluster"
	"appseclabsplataform/utils"
	"log/slog"

	"github.com/google/uuid"
)

type AdminUsecase struct {
	Database   *database.Database
	LabCluster *labcluster.LabClusterService
}

func NewAdminUsecase(database *database.Database, labCluster *labcluster.LabClusterService) *AdminUsecase {
	return &AdminUsecase{
		Database:   database,
		LabCluster: labCluster,
	}
}

func (u *AdminUsecase) GetLabsSessions() ([]GetLabsSessionsResponse, error) {
	labsSessions, err := u.Database.GetAllLabsAttempts()
	if err != nil {
		slog.Error("Failed to get all labs attempts", "error", err)
		return nil, err
	}
	var labsSessionsResponse []GetLabsSessionsResponse
	for _, session := range labsSessions {
		labsSessionsResponse = append(labsSessionsResponse, GetLabsSessionsResponse{
			Namespace:  session.Namespace,
			LabSlug:    session.Lab.Slug,
			UserID:     session.ExternalUserID,
			StartedAt:  session.StartedAt,
			FinishedAt: session.FinishedAt,
			Score:      session.Score,
			Status:     session.Status.Name,
		})
	}
	return labsSessionsResponse, nil
}

func (u *AdminUsecase) GetLabSession(namespace string) (*GetLabSessionResponse, error) {
	labSession, err := u.Database.GetLabAttempt(namespace)
	if err != nil {
		slog.Error("Failed to get lab attempt", "namespace", namespace, "error", err)
		return nil, err
	}

	labClusterResult, err := u.LabCluster.GetLabResult(namespace)
	if err != nil {
		slog.Error("Failed to get lab cluster result", "namespace", namespace, "error", err)
		return nil, err
	}

	labSessionLogs, err := u.Database.GetLabAttemptStatusLogs(labSession.ID)
	if err != nil {
		slog.Error("Failed to get lab session logs", "namespace", namespace, "error", err)
		return nil, err
	}

	return &GetLabSessionResponse{
		Namespace:       labSession.Namespace,
		LabSlug:         labSession.Lab.Slug,
		UserID:          labSession.ExternalUserID,
		StartedAt:       labSession.StartedAt,
		FinishedAt:      labSession.FinishedAt,
		Score:           labSession.Score,
		Status:          labSession.Status.Name,
		LabFinishResult: labClusterResult.FinishResult,
		Logs:            labSessionLogs,
	}, nil
}

func (u *AdminUsecase) GetPossiblesStatus(statusFrom string) ([]GetPossiblesStatusResponse, error) {
	statuses, err := u.Database.GetPossiblesStatusChanges(statusFrom)
	if err != nil {
		slog.Error("Failed to get possible status changes", "statusFrom", statusFrom, "error", err)
		return nil, err
	}

	var response []GetPossiblesStatusResponse
	for _, status := range statuses {
		response = append(response, GetPossiblesStatusResponse{
			Name:  status.Name,
			Value: status.ID,
		})
	}
	return response, nil
}

func (u *AdminUsecase) ChangeLabStatus(namespace string, statusToID uuid.UUID, comment string) error {
	labAttempt, err := u.Database.GetLabAttempt(namespace)
	if err != nil {
		slog.Error("Failed to get lab attempt by namespace", "namespace", namespace, "error", err)
		return err
	}

	statusTo, err := u.Database.GetStatusByID(statusToID)
	if err != nil {
		slog.Error("Failed to get status by ID", "error", err)
		return err
	}

	transition, err := u.Database.GetStatusTransition(labAttempt.Status.Name, statusTo.Name)
	if err != nil {
		slog.Error("Failed to get status transition", "error", err)
		return err
	}

	if _, err := u.Database.UpdateLabAttempt(namespace, *labAttempt, statusTo.Name); err != nil {
		slog.Error("Failed to update lab attempt status", "namespace", namespace, "statusTo", statusTo.Name, "error", err)
		return err
	}

	if _, err := u.Database.AddNewLogStatus(transition.ID, labAttempt.ID, comment, nil); err != nil {
		slog.Error("Failed to add new log status", "transitionID", transition.ID, "labAttemptID", labAttempt.ID, "error", err)
		return err
	}

	return nil
}

func (u *AdminUsecase) GetLabsDefinitions() ([]database.Lab, error) {
	labs, err := u.Database.GetAllLabsDefinitions()
	if err != nil {
		slog.Error("Failed to get labs definitions", "error", err)
		return nil, err
	}

	return labs, nil
}

func (u *AdminUsecase) GetLabDefinition(slug string) (GetLabDefinitionResponse, error) {
	lab, err := u.Database.AdminGetLabDefinition(slug)
	if err != nil {
		slog.Error("Failed to get lab definition", "slug", slug, "error", err)
		return GetLabDefinitionResponse{}, err
	}

	clusterLabDefinition, err := u.LabCluster.GetLabDefinitionBySlug(slug)
	if err != nil {
		slog.Error("Failed to get cluster lab definition", "slug", slug, "error", err)
		return GetLabDefinitionResponse{}, err
	}
	return GetLabDefinitionResponse{
		Lab:           *lab,
		LabDefinition: *clusterLabDefinition,
	}, nil

}

func (u *AdminUsecase) CreateLabDefinition(labDefinition CreateLabDefinitionRequest) (database.Lab, error) {
	lab, err := u.Database.AdminCreateLabDefinition(database.Lab{
		Active:               labDefinition.Active,
		Slug:                 labDefinition.Slug,
		Title:                labDefinition.Title,
		Description:          labDefinition.Description,
		Authors:              labDefinition.Authors,
		ExternalReferences:   labDefinition.ExternalReferences,
		EstimatedTime:        labDefinition.EstimatedTime,
		Difficulty:           labDefinition.Difficulty,
		RequiresManualReview: labDefinition.RequiresManualReview,
		Readme:               labDefinition.Readme,
	}, labDefinition.Vulnerabilities, labDefinition.Languages, labDefinition.Technologies)
	if err != nil {
		slog.Error("Failed to create lab definition", "error", err)
		return database.Lab{}, err
	}

	evaluators := make([]labcluster.LabEvaluator, 0, len(labDefinition.Config.Evaluators))
	for _, evaluator := range labDefinition.Config.Evaluators {
		evaluators = append(evaluators, labcluster.LabEvaluator{
			Slug:            evaluator.Slug,
			Weight:          evaluator.Weight,
			ExploitTemplate: evaluator.ExploitTemplate,
			Config:          evaluator.Config,
		})
	}

	err = u.LabCluster.CreateLabDefinition(labcluster.LabDefinition{
		Slug: labDefinition.Slug,
		LabSpec: labcluster.LabSpec{
			Image: labDefinition.Config.SystemRequirements.Image,
			Env:   labDefinition.Config.SystemRequirements.EnvVars,
			CodeConfig: labcluster.LabCodeConfig{
				GitURL:    labDefinition.Config.SystemRequirements.CodeConfig.GitURL,
				GitBranch: labDefinition.Config.SystemRequirements.CodeConfig.GitBranch,
				GitPath:   labDefinition.Config.SystemRequirements.CodeConfig.GitPath,
			},
		},
		Evaluators: evaluators,
	})
	if err != nil {
		slog.Error("Failed to create lab definition in lab cluster", "error", err)
		// TODO Rollback the lab creation in the database if has error on lab cluster
		return database.Lab{}, err
	}
	return *lab, nil
}
func (u *AdminUsecase) UpdateLabDefinition(slug string, labDefinition UpdateLabDefinitionRequest) (database.Lab, error) {
	lab, err := u.Database.AdminUpdateLabDefinition(database.Lab{
		Slug:                 slug,
		Active:               labDefinition.Active,
		Title:                labDefinition.Title,
		Description:          labDefinition.Description,
		Authors:              labDefinition.Authors,
		ExternalReferences:   labDefinition.ExternalReferences,
		EstimatedTime:        labDefinition.EstimatedTime,
		Difficulty:           labDefinition.Difficulty,
		RequiresManualReview: labDefinition.RequiresManualReview,
		Readme:               labDefinition.Readme,
	}, labDefinition.Vulnerabilities, labDefinition.Languages, labDefinition.Technologies)
	if err != nil {
		slog.Error("Failed to update lab definition", "error", err)
		return database.Lab{}, err
	}

	evaluators := make([]labcluster.LabEvaluator, 0, len(labDefinition.Config.Evaluators))
	for _, evaluator := range labDefinition.Config.Evaluators {
		evaluators = append(evaluators, labcluster.LabEvaluator{
			Slug:            evaluator.Slug,
			Weight:          evaluator.Weight,
			ExploitTemplate: evaluator.ExploitTemplate,
			Config:          evaluator.Config,
		})
	}
	err = u.LabCluster.UpdateLabDefinition(slug, labcluster.LabDefinition{
		Slug: labDefinition.Slug,
		LabSpec: labcluster.LabSpec{
			Image: labDefinition.Config.SystemRequirements.Image,
			Env:   labDefinition.Config.SystemRequirements.EnvVars,
			CodeConfig: labcluster.LabCodeConfig{
				GitURL:    labDefinition.Config.SystemRequirements.CodeConfig.GitURL,
				GitBranch: labDefinition.Config.SystemRequirements.CodeConfig.GitBranch,
				GitPath:   labDefinition.Config.SystemRequirements.CodeConfig.GitPath,
			},
		},
		Evaluators: evaluators,
	})
	if err != nil {
		slog.Error("Failed to update lab definition in lab cluster", "error", err)
		// TODO Rollback the lab update in the database if has error on lab cluster
		return database.Lab{}, err
	}

	return *lab, nil
}
func (u *AdminUsecase) DeleteLabDefinition(slug string) error {
	err := u.Database.AdminDeleteLabDefinition(slug)
	if err != nil {
		slog.Error("Failed to delete lab definition", "slug", slug, "error", err)
		return err
	}

	return nil
}

func (u *AdminUsecase) GetPossiblesVulnerabilities() ([]database.Vulnerability, error) {
	vulnerabilities, err := u.Database.GetPossiblesVulnerabilities()
	if err != nil {
		slog.Error("Failed to get possible vulnerabilities", "error", err)
		return nil, err
	}

	return vulnerabilities, nil
}

func (u *AdminUsecase) GetPossiblesLanguages() ([]database.Language, error) {
	languages, err := u.Database.GetPossiblesLanguages()
	if err != nil {
		slog.Error("Failed to get possible languages", "error", err)
		return nil, err
	}

	return languages, nil
}

func (u *AdminUsecase) GetPossiblesTechnologies() ([]database.Technology, error) {
	technologies, err := u.Database.GetPossiblesTechnologies()
	if err != nil {
		slog.Error("Failed to get possible technologies", "error", err)
		return nil, err
	}

	return technologies, nil
}

func (u *AdminUsecase) GetPossiblesEvaluators() ([]GetPossiblesEvaluatorsResponse, error) {
	evaluators, err := u.LabCluster.GetEvaluators()
	if err != nil {
		slog.Error("Failed to get possible evaluators", "error", err)
		return nil, err
	}

	var response []GetPossiblesEvaluatorsResponse
	for _, evaluator := range evaluators {
		response = append(response, GetPossiblesEvaluatorsResponse{
			Name:  evaluator.Name,
			Value: evaluator.Slug,
		})
	}
	return response, nil
}

func (u *AdminUsecase) GetPossiblesImages() ([]GetPossiblesImagesResponse, error) {

	response := []GetPossiblesImagesResponse{
		{Name: "code-server:dev", Value: "ghcr.io/appsec-digital/code-server:dev"},
	}
	return response, nil
}

func (u *AdminUsecase) GetAllUsers() ([]database.User, error) {
	users, err := u.Database.GetAllUsers()
	if err != nil {
		slog.Error("Failed to get all users", "error", err)
		return nil, err
	}
	return users, nil
}
func (u *AdminUsecase) CreateUser(userRequest CreateUserRequest) (database.User, error) {
	passwordHash, err := utils.HashPassword(userRequest.Password)
	if err != nil {
		slog.Error("Failed to hash password", "error", err)
		return database.User{}, err
	}
	user := database.User{
		Email:        userRequest.Email,
		PasswordHash: passwordHash,
		Name:         userRequest.Name,
		Role:         userRequest.Role,
	}

	createdUser, err := u.Database.CreateUser(user)
	if err != nil {
		slog.Error("Failed to create user", "error", err)
		return database.User{}, err
	}
	return createdUser, nil
}
func (u *AdminUsecase) GetUserByID(id uuid.UUID) (database.User, error) {
	user, err := u.Database.GetUserByID(id)
	if err != nil {
		slog.Error("Failed to get user by ID", "id", id, "error", err)
		return database.User{}, err
	}
	return user, nil
}
func (u *AdminUsecase) UpdateUser(id uuid.UUID, userRequest UpdateUserRequest) (database.User, error) {
	passwordHash, err := utils.HashPassword(userRequest.Password)
	if err != nil {
		slog.Error("Failed to hash password", "error", err)
		return database.User{}, err
	}
	user := database.User{
		ID:           id,
		Email:        userRequest.Email,
		PasswordHash: passwordHash,
		Name:         userRequest.Name,
		Role:         userRequest.Role,
	}
	updatedUser, err := u.Database.UpdateUser(user)
	if err != nil {
		slog.Error("Failed to update user", "id", id, "error", err)
		return database.User{}, err
	}
	return updatedUser, nil
}
func (u *AdminUsecase) DeleteUser(id uuid.UUID) error {
	err := u.Database.DeleteUser(id)
	if err != nil {
		slog.Error("Failed to delete user", "id", id, "error", err)
		return err
	}
	return nil
}
