package database

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Language struct {
	ID        uuid.UUID      `json:"id" db:"id" gorm:"type:uuid;primaryKey"`
	Name      string         `json:"name" db:"name" gorm:"unique;not null"`
	CreatedAt time.Time      `json:"-" db:"created_at"`
	UpdatedAt time.Time      `json:"-" db:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" db:"deleted_at"`

	// Relationships
	Labs []Lab `json:"-" gorm:"many2many:lab_languages;"`
}

func (db *Database) GetPossiblesLanguages() ([]Language, error) {
	var languages []Language
	if err := db.Conn.Model(&Language{}).
		Find(&languages).Error; err != nil {
		return nil, fmt.Errorf("failed to get possible languages: %w", err)
	}
	return languages, nil
}
