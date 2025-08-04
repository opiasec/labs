package admin

import (
	"appseclabsplataform/database"
	labcluster "appseclabsplataform/services/labCluster"
	"time"

	"github.com/google/uuid"
)

type GetLabsSessionsResponse struct {
	Namespace  uuid.UUID `json:"namespace"`
	LabSlug    string    `json:"labSlug"`
	UserID     string    `json:"userId"`
	StartedAt  time.Time `json:"startedAt"`
	FinishedAt time.Time `json:"finishedAt"`
	Score      int       `json:"score"`
	Status     string    `json:"status"`
}

type GetLabSessionResponse struct {
	Namespace       uuid.UUID               `json:"namespace"`
	LabSlug         string                  `json:"labSlug"`
	UserID          string                  `json:"userId"`
	StartedAt       time.Time               `json:"startedAt"`
	FinishedAt      time.Time               `json:"finishedAt"`
	Score           int                     `json:"score"`
	Status          string                  `json:"status"`
	LabFinishResult labcluster.FinishResult `json:"finishResult"`
	Logs            []database.LabStatusLog `json:"logs"`
}

type GetLabSessionResponseLogs struct {
	database.LabStatusLog
	StatusFrom string `json:"statusFrom"`
	StatusTo   string `json:"statusTo"`
}

type GetPossibleResponse struct {
	Label string    `json:"name"`
	Value uuid.UUID `json:"value"`
}

type GetPossiblesEvaluatorsResponse struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type GetPossiblesImagesResponse struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type GetPossiblesStatusResponse struct {
	Name  string    `json:"name"`
	Value uuid.UUID `json:"value"`
}

type ChangeLabStatusRequest struct {
	StatusID uuid.UUID `json:"statusId"`
	Comment  string    `json:"comment,omitempty"`
}

type GetLabDefinitionResponse struct {
	database.Lab
	labcluster.LabDefinition `json:"config"`
}

type CreateLabDefinitionRequest struct {
	Slug                 string                           `json:"slug"`
	Title                string                           `json:"title"`
	Description          string                           `json:"description"`
	Authors              []string                         `json:"authors"`
	ExternalReferences   []string                         `json:"externalReferences"`
	Vulnerabilities      []database.Vulnerability         `json:"vulnerabilities"`
	Languages            []database.Language              `json:"languages"`
	Technologies         []database.Technology            `json:"technologies"`
	EstimatedTime        int                              `json:"estimatedTime"`
	RequiresManualReview bool                             `json:"requiresManualReview"`
	Readme               string                           `json:"readme"`
	Difficulty           string                           `json:"difficulty"`
	Config               CreateLabDefinitionRequestConfig `json:"config"`
	Active               bool                             `json:"active"`
}

type CreateLabDefinitionRequestConfig struct {
	Evaluators         []CreateLabDefinitionRequestEvaluator        `json:"evaluators"`
	SystemRequirements CreateLabDefinitionRequestSystemRequirements `json:"systemRequirements"`
}

type CreateLabDefinitionRequestEvaluator struct {
	Slug            string            `json:"slug"`
	Weight          int               `json:"weight"`
	ExploitTemplate string            `json:"exploitTemplate,omitempty"`
	Config          map[string]string `json:"config,omitempty"`
}

type CreateLabDefinitionRequestSystemRequirements struct {
	Image      string                                                 `json:"image"`
	EnvVars    map[string]string                                      `json:"envVars"`
	CodeConfig CreateLabDefinitionRequestSystemRequirementsCodeConfig `json:"codeConfig"`
}
type CreateLabDefinitionRequestSystemRequirementsCodeConfig struct {
	GitURL    string `json:"gitUrl"`
	GitBranch string `json:"gitBranch"`
	GitPath   string `json:"gitPath,omitempty"`
}

type UpdateLabDefinitionRequest struct {
	CreateLabDefinitionRequest
}

type CreateUserRequest struct {
	ID            uuid.UUID `json:"id"`
	Email         string    `json:"email"`
	Password      string    `json:"password"`
	Name          string    `json:"name"`
	ImageURL      string    `json:"image_url"`
	Role          string    `json:"role"`
	IsActive      bool      `json:"is_active"`
	EmailVerified bool      `json:"email_verified"`
}
type UpdateUserRequest struct {
	ID            uuid.UUID `json:"id"`
	Email         string    `json:"email"`
	Password      string    `json:"password"`
	Name          string    `json:"name"`
	ImageURL      string    `json:"image_url"`
	Role          string    `json:"role"`
	IsActive      bool      `json:"is_active"`
	EmailVerified bool      `json:"email_verified"`
}
