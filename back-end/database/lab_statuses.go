package database

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LabStatus struct {
	ID          uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"unique;not null" json:"name"`
	Description string         `gorm:"not null" json:"description"`
	IsInitial   bool           `gorm:"default:false" json:"isInitial"`
	IsFinal     bool           `gorm:"default:false" json:"isFinal"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

func (db *Database) GetPossiblesStatusChanges(statusFrom string) ([]LabStatus, error) {
	var statuses []LabStatus

	err := db.Conn.Table("lab_statuses").
		Select("lab_statuses.*").
		Joins("JOIN lab_transitions ON lab_transitions.to_status_id = lab_statuses.id").
		Joins("JOIN lab_statuses fs ON lab_transitions.from_status_id = fs.id").
		Where("fs.name = ? AND lab_transitions.deleted_at IS NULL", statusFrom).
		Find(&statuses).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get possible status changes from %s: %w", statusFrom, err)
	}

	return statuses, nil
}

func (db *Database) GetStatusByID(statusID uuid.UUID) (*LabStatus, error) {
	var result LabStatus
	if err := db.Conn.Where("id = ?", statusID).First(&result).Error; err != nil {
		return nil, fmt.Errorf("failed to get lab status by ID: %w", err)
	}
	return &result, nil
}
