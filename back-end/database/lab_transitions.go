package database

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LabTransition struct {
	ID               uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	FromStatusID     uuid.UUID      `gorm:"type:uuid;not null" json:"fromStatusId"`
	ToStatusID       uuid.UUID      `gorm:"type:uuid;not null" json:"toStatusId"`
	Name             string         `gorm:"type:varchar(100);not null" json:"name"`
	IsAutomatic      bool           `gorm:"default:false" json:"isAutomatic"`
	RequiresApproval bool           `gorm:"default:false" json:"requiresApproval"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deletedAt"`

	// Relationships
	FromStatus LabStatus `gorm:"foreignKey:FromStatusID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"fromStatus"`
	ToStatus   LabStatus `gorm:"foreignKey:ToStatusID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"toStatus"`
}

func (db *Database) GetStatusTransition(fromStatusName string, toStatusName string) (*LabTransition, error) {
	var labTransition LabTransition

	err := db.Conn.
		Joins("JOIN lab_statuses fs ON lab_transitions.from_status_id = fs.id").
		Joins("JOIN lab_statuses ts ON lab_transitions.to_status_id = ts.id").
		Where("fs.name = ? AND ts.name = ?", fromStatusName, toStatusName).
		First(&labTransition).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get lab transition from %s to %s: %w", fromStatusName, toStatusName, err)
	}
	return &labTransition, nil
}
