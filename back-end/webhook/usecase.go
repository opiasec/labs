package webhook

import (
	"appseclabsplataform/database"
	"log/slog"
)

type WebHookUsecase struct {
	database *database.Database // Assuming you have a Database struct in your database package
}

func NewWebHookUsecase(database *database.Database) *WebHookUsecase {
	return &WebHookUsecase{
		database: database,
	}
}

func (u *WebHookUsecase) FinishEvaluationResult(sessionWebHook LabSession) error {
	var statusLog string
	changeStatus := "passed"

	session, err := u.database.GetLabAttempt(sessionWebHook.Namespace)
	if err != nil {
		slog.Error("Failed to get lab attempt", "error", err)
		return err
	}

	if sessionWebHook.FinishResult.Status == "" {
		slog.Error("Invalid status in FinishResult", "status", "empty")
	}
	if sessionWebHook.FinishResult.Status == "failed" {
		changeStatus = "failed"
	}
	if sessionWebHook.FinishResult.Status == "completed" {
		if session.Lab.RequiresManualReview {
			changeStatus = "pending_review"
		} else {
			changeStatus = "passed"
		}
	}

	transaction, err := u.database.GetStatusTransition(session.Status.Name, changeStatus)
	if err != nil {
		slog.Error("Failed to check status transition", "error", err)
		return err
	}

	statusLog = "Status changed from " + session.Status.Name + " to " + changeStatus
	if _, err := u.database.UpdateLabAttempt(sessionWebHook.Namespace, *session, changeStatus); err != nil {
		slog.Error("Failed to update lab attempt", "error", err)
		return err
	}
	if _, err := u.database.AddNewLogStatus(transaction.ID, session.ID, statusLog, nil); err != nil {
		slog.Error("Failed to create lab status log", "error", err)
		return err
	}

	return nil

}
