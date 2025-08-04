package labcluster

import (
	"time"

	"github.com/google/uuid"
)

type CreateLabRequest struct {
	LabSlug string `json:"labSlug"`
	UserID  string `json:"userID"`
}

type CreateLabResponse struct {
	IDEURL      string    `json:"ideUrl"`
	LabPassword string    `json:"labPassword"`
	Namespace   uuid.UUID `json:"namespace"`
}

type GetLabResponse struct {
	IDEURL        string                `json:"ideUrl"`
	AppURL        string                `json:"appUrl"`
	Namespace     string                `json:"namespace"`
	StatusMessage string                `json:"statusMessage"`
	Containers    []ContainerStatusInfo `json:"containers"`
	ExpiresAt     time.Time             `json:"expiresAt"`
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

type RedeployLabResponse struct {
	Message string `json:"message"`
}

type FinishLabResponse struct {
	Message string `json:"message"`
}

type CriterionResult struct {
	Name     string `json:"name"`
	Score    int    `json:"score"`
	Weight   int    `json:"weight"`
	Required bool   `json:"required"`
	Message  string `json:"message"`
}

type LabDefinition struct {
	LabSpec    LabSpec        `json:"labSpec,omitempty"`
	Slug       string         `json:"slug"`
	CreatedAt  time.Time      `json:"createdAt,omitempty"`
	UpdatedAt  time.Time      `json:"updatedAt,omitempty"`
	Evaluators []LabEvaluator `json:"evaluators"`
}

type LabSpec struct {
	Image      string            `json:"image"`
	Env        map[string]string `json:"envVars"`
	CodeConfig LabCodeConfig     `json:"codeConfig"`
}

type LabCodeConfig struct {
	GitURL    string `json:"gitUrl"`
	GitBranch string `json:"gitBranch"`
	GitPath   string `json:"gitPath"`
}

type LabEvaluator struct {
	Slug            string            `json:"slug"`
	Weight          int               `json:"weight"`
	ExploitTemplate string            `json:"exploitTemplate,omitempty"`
	Config          map[string]string `json:"config,omitempty"`
}

type GetLabResultResponse struct {
	LabSlug      string       `json:"labSlug"`
	Namespace    string       `json:"namespace"`
	FinishResult FinishResult `json:"finishResult"`
	StartedAt    time.Time    `json:"startedAt"`
	FinishedAt   time.Time    `json:"finishedAt"`
	CreatedAt    time.Time    `json:"createdAt"`
	UpdatedAt    time.Time    `json:"updatedAt"`
}

type LabFinishRequest struct {
	LabSlug string `json:"labSlug"`
}

type FinishResult struct {
	Status         string                     `bson:"status" json:"status"`
	ErrorMessage   string                     `bson:"error_message" json:"errorMessage"`
	TotalScore     int                        `bson:"total_score" json:"totalScore"`
	CriteriaResult []LabFinishResultCriterion `bson:"criteria_result" json:"criteriaResult"`
	FilesDiff      string                     `bson:"files_diff" json:"filesDiff"`
}

type LabFinishResultCriterion struct {
	Name      string `bson:"name" json:"name"`
	Score     int    `bson:"score" json:"score"`
	Weight    int    `bson:"weight" json:"weight"`
	Required  bool   `bson:"required" json:"required"`
	Message   string `bson:"message" json:"message"`
	Status    string `bson:"status" json:"status"`
	RawOutput string `bson:"raw_output" json:"rawOutput"`
}

type GetAllLabsByUserAndStatusResponse struct {
	Namespace       string    `json:"namespace"`
	LabSlug         string    `json:"labSlug"`
	StartedAt       time.Time `json:"startedAt"`
	FinishedAt      time.Time `json:"finishedAt"`
	Score           int       `json:"score"`
	DurationSeconds int       `json:"durationSeconds"`
	Rating          int       `json:"rating"`
	UserFeedback    string    `json:"userFeedback"`
}

type SendFeedbackRequest struct {
	Rating   int    `json:"rating"`
	Feedback string `json:"feedback"`
}

type Evaluator struct {
	Name           string         `bson:"name"`
	Description    string         `bson:"description"`
	EvaluationSpec EvaluationSpec `bson:"evaluation_spec"`
	Slug           string         `bson:"slug"`
	CreatedAt      time.Time      `bson:"created_at"`
	UpdatedAt      time.Time      `bson:"updated_at"`
}

type EvaluationSpec struct {
	Containers    []ContainerSpec `bson:"containers"`
	InitContainer ContainerSpec   `bson:"init_container"`
	Volumes       []VolumeSpec    `bson:"volumes"`
}
type ContainerSpec struct {
	Name     string         `bson:"name"`
	Image    string         `bson:"image"`
	Env      []ContainerEnv `bson:"env"`
	Args     []string       `bson:"args"`
	Commands []string       `bson:"commands"`
	Volumes  []VolumeSpec   `bson:"volumes"`
}

type ContainerEnv struct {
	Name  string `bson:"name"`
	Value string `bson:"value"`
}
type VolumeSpec struct {
	Name string `bson:"name"`
	Path string `bson:"path"`
}
