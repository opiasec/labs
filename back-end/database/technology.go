package database

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Technology struct {
	ID        uuid.UUID      `json:"id" db:"id" gorm:"type:uuid;primaryKey"`
	Name      string         `json:"name" db:"name" gorm:"unique;not null"`
	CreatedAt time.Time      `json:"-" db:"created_at"`
	UpdatedAt time.Time      `json:"-" db:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" db:"deleted_at"`

	// Relationships
	Labs []Lab `json:"-" gorm:"many2many:lab_technologies;"`
}

func (db *Database) GetPossiblesTechnologies() ([]Technology, error) {
	var technologies []Technology
	if err := db.Conn.Model(&Technology{}).
		Find(&technologies).Error; err != nil {
		return nil, fmt.Errorf("failed to get possible technologies: %w", err)
	}
	return technologies, nil
}
