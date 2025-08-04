package labsession

import "github.com/google/uuid"

type LabSession struct {
	LabAttemptID uuid.UUID `json:"labAttemptId"`
	LabSlug      string    `json:"labSlug"`
	Namespace    uuid.UUID `json:"namespace"`
	UserID       string    `json:"userId"`
	Status       string    `json:"status"`
	IDEURL       string    `json:"ideUrl"`
	Password     string    `json:"password,omitempty"`
}
