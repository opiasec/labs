package lab

import (
	labcluster "appseclabsplataform/services/labCluster"
	labide "appseclabsplataform/services/labIDE"
	"time"

	"github.com/google/uuid"
)

type CreateLabRequest struct {
	LabSlug string `json:"labSlug"`
}

type CreateLabResponse struct {
	labcluster.CreateLabResponse
	ExpiresAt time.Time `json:"expiresAt"`
}

type GetLabRequest struct {
	Namespace string `param:"namespace"`
}

type GetLabResponse struct {
	labide.GetLabResponse
}

type DeleteLabRequest struct {
	Namespace string `param:"namespace"`
}

type RedeployLabRequest struct {
	Namespace string `param:"namespace"`
}

type RedeployLabResponse struct {
	labcluster.RedeployLabResponse
}

type FinishLabRequest struct {
	Namespace string `param:"namespace"`
}

type LeaveLabRequest struct {
	Namespace string `param:"namespace"`
}

type FinishLabResponse struct {
	labcluster.FinishLabResponse
}

type GetLabResultRequest struct {
	Namespace string `param:"namespace"`
}

type GetLabResultResponse struct {
	Namespace       uuid.UUID               `json:"namespace"`
	LabSlug         string                  `json:"labSlug"`
	Rating          int                     `json:"rating"`
	Status          string                  `json:"status"`
	UserFeedback    string                  `json:"userFeedback"`
	StartedAt       time.Time               `json:"startedAt"`
	FinishedAt      time.Time               `json:"finishedAt"`
	DurationSeconds int                     `json:"durationSeconds"`
	LabFinishResult labcluster.FinishResult `json:"finishResult"`
}

type GetAllLabsByUserAndStatusRequest struct {
	Status string `query:"status"`
}

type GetAllLabsByUserAndStatusResponse struct {
	Namespace       uuid.UUID `json:"namespace"`
	LabSlug         string    `json:"labSlug"`
	StartedAt       time.Time `json:"startedAt"`
	FinishedAt      time.Time `json:"finishedAt"`
	Score           int       `json:"score"`
	Status          string    `json:"status"`
	DurationSeconds int       `json:"durationSeconds"`
	Rating          int       `json:"rating"`
	UserFeedback    string    `json:"userFeedback"`
}

type SendFeedbackRequest struct {
	Namespace string `param:"namespace"`
	Rating    int    `json:"rating"`
	Feedback  string `json:"feedback"`
}
