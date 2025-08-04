package lab

import (
	"appseclabs/types"
	"time"
)

type CreateLabRequest struct {
	LabSlug string `json:"labSlug"`
}

type CreateLabResponse struct {
	IDEURL      string `json:"ideUrl"`
	LabPassword string `json:"labPassword"`
	Namespace   string `json:"namespace"`
}

type DeleteLabRequest struct {
	UserID    string `json:"userID"`
	Namespace string `param:"namespace"`
}

type GetLabStatusRequest struct {
	UserID    string `json:"userID"`
	Namespace string `param:"namespace"`
}

type GetLabStatusResponse struct {
	IDEURL     string                `json:"ideUrl"`
	Namespace  string                `json:"namespace"`
	LabSlug    string                `json:"labSlug"`
	Status     string                `json:"status"`
	Containers []ContainerStatusInfo `json:"containers,omitempty"`
	ExpiresAt  time.Time             `json:"expiresAt"`
}

type ContainerStatusInfo struct {
	Name       string    `json:"name"`
	Ready      bool      `json:"ready"`
	State      string    `json:"state"`
	Reason     string    `json:"reason"`
	ExitCode   int32     `json:"exitCode,omitempty"`
	StartedAt  time.Time `json:"startedAt,omitempty"`
	RestartCnt int32     `json:"restartCnt,omitempty"`
}
type RedeployLabRequest struct {
	Namespace string `param:"namespace"`
}

type RedeployLabResponse struct {
	Message string `json:"message"`
}

type FinishLabRequest struct {
	Namespace string `param:"namespace"`
	LabSlug   string `json:"labSlug"`
}

type CriterionContainerResult struct {
	Score   int    `json:"score"`
	Message string `json:"message"`
}

type FinishLabResponse struct {
	Message string `json:"message"`
}

type GetLabResultRequest struct {
	Namespace string `param:"namespace"`
}

type GetLabResultResponse struct {
	Namespace       string                `json:"namespace"`
	LabSlug         string                `json:"labSlug"`
	StartedAt       time.Time             `json:"startedAt"`
	FinishedAt      time.Time             `json:"finishedAt"`
	DurationSeconds int                   `json:"durationSeconds"`
	LabFinishResult types.LabFinishResult `json:"labFinishResult"`
}

type PostLabFeedbackRequest struct {
	Namespace string `param:"namespace"`
	Feedback  string `json:"feedback"`
	Rating    int    `json:"rating"`
}

type PostLabFeedbackResponse struct {
	Message string `json:"message"`
}

type GetAllLabsSessionByUserRequest struct {
	Status string `query:"status"`
}

type GetAllLabsSessionResponse struct {
	Namespace       string    `json:"namespace"`
	LabSlug         string    `json:"labSlug"`
	StartedAt       time.Time `json:"startedAt"`
	FinishedAt      time.Time `json:"finishedAt"`
	DurationSeconds int       `json:"durationSeconds"`
	Score           int       `json:"score"`
	Rating          int       `json:"rating"`
	UserFeedback    string    `json:"userFeedback"`
}
