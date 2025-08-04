package database

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LabAttempt struct {
	ID             uuid.UUID      `json:"id" db:"id" gorm:"type:uuid;primaryKey"`
	Namespace      uuid.UUID      `json:"namespace" db:"namespace" gorm:"unique;not null"`
	ExternalUserID string         `json:"externalUserId" db:"external_user_id" gorm:"not null"`
	LabID          uuid.UUID      `json:"labId" db:"lab_id" gorm:"type:uuid;not null"`
	StartedAt      time.Time      `json:"startedAt" db:"started_at" gorm:"not null"`
	FinishedAt     time.Time      `json:"finishedAt" db:"finished_at"`
	StatusID       uuid.UUID      `json:"statusId" db:"status_id" gorm:"type:uuid"`
	Score          int            `json:"score" db:"score" gorm:"not null;check:score BETWEEN 0 AND 100"`
	Feedback       string         `json:"feedback" db:"feedback"`
	Rating         int            `json:"rating" db:"rating" gorm:"check:rating BETWEEN 0 AND 5"`
	CreatedAt      time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time      `json:"updatedAt" db:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deletedAt" db:"deleted_at"`

	// Relationships
	Lab    Lab       `gorm:"foreignKey:LabID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Status LabStatus `gorm:"foreignKey:StatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

func (db *Database) CreateLabAttempt(namespace uuid.UUID, labID uuid.UUID, userID string) error {
	var status LabStatus
	if err := db.Conn.Where("name = ?", "running").First(&status).Error; err != nil {
		return fmt.Errorf("failed to get 'running' status: %w", err)
	}

	startedAt := time.Now().UTC()
	labAttempt := LabAttempt{
		ID:             uuid.New(),
		Namespace:      namespace,
		LabID:          labID,
		StartedAt:      startedAt,
		Status:         status,
		Score:          0,
		ExternalUserID: userID,
	}
	if err := db.Conn.Create(&labAttempt).Error; err != nil {
		return fmt.Errorf("failed to create lab attempt: %w", err)
	}

	return nil
}

func (db *Database) UpdateLabAttempt(namespace string, labAttempt LabAttempt, status string) (*LabAttempt, error) {
	var statusObj LabStatus
	if err := db.Conn.Where("name = ?", status).First(&statusObj).Error; err != nil {
		return nil, fmt.Errorf("failed to get status '%s': %w", status, err)
	}

	var existingLabAttempt LabAttempt
	if err := db.Conn.Where("namespace = ?", namespace).First(&existingLabAttempt).Error; err != nil {
		return nil, fmt.Errorf("failed to get lab attempt by namespace: %w", err)
	}

	labAttempt.Status = statusObj
	existingLabAttempt = labAttempt
	if err := db.Conn.Save(&existingLabAttempt).Error; err != nil {
		return nil, fmt.Errorf("failed to update lab attempt: %w", err)
	}

	return &existingLabAttempt, nil
}

func (db *Database) SetLabAttemptAsTimeout(namespace string) error {
	var labAttempt LabAttempt
	if err := db.Conn.Where("namespace = ?", namespace).First(&labAttempt).Error; err != nil {
		return fmt.Errorf("failed to get lab attempt by namespace: %w", err)
	}

	var status LabStatus
	if err := db.Conn.Where("name = ?", "timeout").First(&status).Error; err != nil {
		return fmt.Errorf("failed to get 'timeout' status: %w", err)
	}

	labAttempt.Status = status
	if err := db.Conn.Save(&labAttempt).Error; err != nil {
		return fmt.Errorf("failed to update lab attempt status: %w", err)
	}

	return nil
}

func (db *Database) GetLabAttemptsByUser(userID string) ([]LabAttempt, error) {
	var labAttempts []LabAttempt
	if err := db.Conn.Preload("Lab").Preload("Status").Where("external_user_id = ?", userID).Find(&labAttempts).Error; err != nil {
		return nil, fmt.Errorf("failed to get lab attempts by user: %w", err)
	}
	return labAttempts, nil
}

func (db *Database) GetLabAttempt(namespace string) (*LabAttempt, error) {
	var labAttempt LabAttempt
	if err := db.Conn.Preload("Lab").
		Preload("Status").
		Where("namespace = ?", namespace).First(&labAttempt).Error; err != nil {
		return nil, fmt.Errorf("failed to get lab attempt by namespace: %w", err)
	}
	return &labAttempt, nil
}

func (db *Database) FinishLabAttempt(namespace string, score int) error {
	var status LabStatus
	if err := db.Conn.Where("name = ?", "evaluating").First(&status).Error; err != nil {
		return fmt.Errorf("failed to get 'evaluating' status: %w", err)
	}

	var labAttempt LabAttempt
	if err := db.Conn.Where("namespace = ?", namespace).First(&labAttempt).Error; err != nil {
		return fmt.Errorf("failed to get lab attempt by namespace: %w", err)
	}

	labAttempt.Status = status
	labAttempt.Score = score
	labAttempt.FinishedAt = time.Now().UTC()
	if err := db.Conn.Save(&labAttempt).Error; err != nil {
		return fmt.Errorf("failed to update lab attempt status: %w", err)
	}

	return nil
}

func (db *Database) SetLabFeedback(namespace string, rating int, feedback string) error {
	var labAttempt LabAttempt
	if err := db.Conn.Where("namespace = ?", namespace).First(&labAttempt).Error; err != nil {
		return fmt.Errorf("failed to get lab attempt by namespace: %w", err)
	}

	labAttempt.Rating = rating
	labAttempt.Feedback = feedback
	if err := db.Conn.Save(&labAttempt).Error; err != nil {
		return fmt.Errorf("failed to update lab attempt feedback: %w", err)
	}

	return nil
}

func (db *Database) GetAllLabsAttempts() ([]LabAttempt, error) {
	var labAttempts []LabAttempt
	if err := db.Conn.Preload("Lab").Preload("Status").Find(&labAttempts).Error; err != nil {
		return nil, fmt.Errorf("failed to get all lab attempts: %w", err)
	}
	return labAttempts, nil
}

func (db *Database) GetLabAttemptsByUserFilterByDay(userID string, startDate, endDate time.Time) ([]LabAttempt, error) {
	var labAttempts []LabAttempt
	if err := db.Conn.Model(&LabAttempt{}).
		Preload("Lab").
		Preload("Status").
		Where("external_user_id = ?", userID).
		Where("created_at >= ?", startDate).
		Where("created_at <= ?", endDate).
		Find(&labAttempts).Error; err != nil {
		return nil, fmt.Errorf("failed to get lab attempts by user and date: %w", err)
	}
	return labAttempts, nil
}

func (db *Database) LeaveLabAttempt(namespace string) (*LabAttempt, error) {
	var status LabStatus
	if err := db.Conn.Where("name = ?", "abandoned").First(&status).Error; err != nil {
		return nil, fmt.Errorf("failed to get 'abandoned' status: %w", err)
	}

	var labAttempt LabAttempt
	if err := db.Conn.Where("namespace = ?", namespace).First(&labAttempt).Error; err != nil {
		return nil, fmt.Errorf("failed to get lab attempt by namespace: %w", err)
	}

	labAttempt.Status = status
	labAttempt.FinishedAt = time.Now().UTC()
	if err := db.Conn.Save(&labAttempt).Error; err != nil {
		return nil, fmt.Errorf("failed to update lab attempt status: %w", err)
	}

	return &labAttempt, nil
}
