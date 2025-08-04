package database

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LabStatusLog struct {
	ID              uuid.UUID              `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	LabAttemptID    uuid.UUID              `gorm:"type:uuid;not null" json:"labAttemptId"`
	LabTransitionID uuid.UUID              `gorm:"type:uuid;not null" json:"labTransitionId"`
	ChangedBy       uuid.UUID              `gorm:"type:uuid" json:"changedBy"`
	Comment         string                 `gorm:"type:text" json:"comment"`
	Metadata        map[string]interface{} `gorm:"type:jsonb" json:"metadata"`
	CreatedAt       time.Time              `gorm:"autoCreateTime" json:"createdAt"`
	DeletedAt       gorm.DeletedAt         `gorm:"index" json:"deletedAt"`

	// Relationships
	LabAttempt    LabAttempt    `gorm:"foreignKey:LabAttemptID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	LabTransition LabTransition `gorm:"foreignKey:LabTransitionID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"labTransition"`
}

func (db *Database) AddNewLogStatus(labTransitionID uuid.UUID, labAttemptId uuid.UUID, comment string, metadata map[string]interface{}) (*LabStatusLog, error) {
	var labLog LabStatusLog

	labLog.LabTransitionID = labTransitionID
	labLog.LabAttemptID = labAttemptId
	labLog.Comment = comment
	labLog.Metadata = metadata

	if err := db.Conn.Create(&labLog).Error; err != nil {
		return nil, fmt.Errorf("failed to create status log: %w", err)
	}

	return &labLog, nil
}

func (db *Database) GetLabAttemptStatusLogs(labAttemptid uuid.UUID) ([]LabStatusLog, error) {
	var labLogs []LabStatusLog
	// Get LabStatusLogs
	if err := db.Conn.Preload("LabTransition").
		Preload("LabTransition.FromStatus").
		Preload("LabTransition.ToStatus").
		Where("lab_attempt_id = ?", labAttemptid).Find(&labLogs).Error; err != nil {
		return nil, fmt.Errorf("failed to get lab status logs: %w", err)
	}
	return labLogs, nil
}
