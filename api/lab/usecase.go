package lab

import (
	"appseclabs/database"
	"appseclabs/k8s"
	"appseclabs/services/webhook"
	"appseclabs/types"
	"appseclabs/utils"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
)

type LabUsecase struct {
	KubeClient     *k8s.K8s
	Database       *database.Database
	LabBaseURL     string
	LabExpiration  time.Duration
	Logger         *zap.SugaredLogger
	WebhookService *webhook.WebhookService
}

func NewLabUsecase(kubeClient *k8s.K8s, database *database.Database, logger *zap.SugaredLogger, webhookService *webhook.WebhookService) *LabUsecase {

	labBaseURL := os.Getenv("LAB_BASE_URL")
	return &LabUsecase{
		KubeClient:     kubeClient,
		Database:       database,
		LabBaseURL:     labBaseURL,
		LabExpiration:  time.Hour * 1,
		Logger:         logger,
		WebhookService: webhookService,
	}
}

func (u *LabUsecase) CreateLab(labSlug string) (*CreateLabResponse, error) {

	namespace := utils.GenerateUUID()

	IDEURL := fmt.Sprintf("%s/%s/", u.LabBaseURL, namespace)

	lab, err := u.Database.GetLab(labSlug)

	if err != nil {
		u.Logger.Error("Error getting lab", zap.Error(err))
		return nil, err
	}

	codeServerPassword, err := u.KubeClient.CreateLab(namespace, lab)
	if err != nil {
		u.Logger.Error("Error creating lab", zap.Error(err))
		return nil, err
	}

	err = u.Database.SaveLabSession(types.LabSession{
		Namespace: namespace,
		LabSlug:   labSlug,
	})

	if err != nil {
		u.Logger.Error("Error saving lab session", zap.Error(err))
		return nil, err
	}

	return &CreateLabResponse{
		IDEURL:      IDEURL,
		LabPassword: codeServerPassword,
		Namespace:   namespace,
	}, nil

}

func (u *LabUsecase) DeleteLab(namespace string) error {

	err := u.KubeClient.DeleteLab(namespace)
	if err != nil {
		u.Logger.Error("Error deleting lab", zap.Error(err))
		return err
	}

	return nil
}

func (u *LabUsecase) FinishLab(namespace, labSlug string) (FinishLabResponse, error) {

	lab, err := u.Database.GetLab(labSlug)
	if err != nil {
		u.Logger.Error("Error getting lab", zap.Error(err))
		return FinishLabResponse{}, err
	}

	err = u.KubeClient.ScaleCodeServer(namespace, labSlug, 0)
	if err != nil {
		u.Logger.Error("Error scaling code server", zap.Error(err))
		return FinishLabResponse{}, err
	}

	go func() {
		finishResult, err := u.KubeClient.CreateEvaluation(namespace, lab)
		if err != nil {
			u.Logger.Error("Error creating evaluation", zap.Error(err))
			return
		}
		labSessionTemp := types.LabSession{
			Namespace:    namespace,
			LabSlug:      labSlug,
			FinishResult: finishResult,
			FinishedAt:   time.Now(),
		}
		err = u.Database.UpdateLabSession(namespace, labSessionTemp)
		if err != nil {
			u.Logger.Error("Error saving lab session", zap.Error(err))
			return
		}
		err = u.KubeClient.DeleteLab(namespace)
		if err != nil {
			u.Logger.Error("Error deleting lab", zap.Error(err))
			return
		}
		err = u.WebhookService.SendFinishEvaluationResult(labSessionTemp)
		if err != nil {
			u.Logger.Error("Error sending finish evaluation result to webhook", zap.Error(err))
			return
		}

	}()
	return FinishLabResponse{
		Message: "finish lab started",
	}, nil
}

func (u *LabUsecase) GetLabResult(namespace string) (*types.LabSession, error) {
	mongoLabResult, err := u.Database.GetLabSession(namespace)
	if err != nil {
		u.Logger.Error("Error getting lab result", zap.Error(err))
		return nil, err
	}
	return mongoLabResult, nil
}
