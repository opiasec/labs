package database

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Lab struct {
	ID                   uuid.UUID      `json:"id" db:"id" gorm:"type:uuid;primaryKey"`
	Slug                 string         `json:"slug" db:"slug" gorm:"unique;not null"`
	Title                string         `json:"title" db:"title" gorm:"not null"`
	Authors              pq.StringArray `json:"authors" db:"authors" gorm:"type:text[]"`
	ExternalReferences   pq.StringArray `json:"externalReferences" db:"external_references" gorm:"type:text[]"`
	EstimatedTime        int            `json:"estimatedTime" db:"estimated_time" gorm:"not null"`
	Difficulty           string         `json:"difficulty" db:"difficulty" gorm:"not null;check:difficulty IN ('easy', 'medium', 'hard')"`
	Description          string         `json:"description" db:"description" gorm:"not null"`
	RequiresManualReview bool           `json:"requiresManualReview" db:"requires_manual_review" gorm:"default:false"`
	Readme               string         `json:"readme" db:"readme" gorm:"type:text"`
	Tags                 pq.StringArray `json:"tags" db:"tags" gorm:"type:text[]"`
	Active               bool           `json:"active" db:"active" gorm:"default:true"`
	CreatedAt            time.Time      `json:"-" db:"created_at"`
	UpdatedAt            time.Time      `json:"-" db:"updated_at"`
	DeletedAt            gorm.DeletedAt `json:"-" db:"deleted_at"`

	// Relationships
	Vulnerabilities []Vulnerability `json:"vulnerabilities" gorm:"many2many:lab_vulnerabilities;"`
	Languages       []Language      `json:"languages" gorm:"many2many:lab_languages;"`
	Technologies    []Technology    `json:"technologies" gorm:"many2many:lab_technologies;"`
}

type LabVulnerability struct {
	LabID           uuid.UUID      `json:"labId" db:"lab_id" gorm:"type:uuid;primaryKey"`
	VulnerabilityID uuid.UUID      `json:"vulnerabilityId" db:"vulnerability_id" gorm:"type:uuid;primaryKey"`
	CreatedAt       time.Time      `json:"-" db:"created_at"`
	UpdatedAt       time.Time      `json:"-" db:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" db:"deleted_at"`
}

type LabLanguage struct {
	LabID      uuid.UUID      `json:"labId" db:"lab_id" gorm:"type:uuid;primaryKey"`
	LanguageID uuid.UUID      `json:"languageId" db:"language_id" gorm:"type:uuid;primaryKey"`
	CreatedAt  time.Time      `json:"-" db:"created_at"`
	UpdatedAt  time.Time      `json:"-" db:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" db:"deleted_at"`
}

type LabTechnology struct {
	LabID        uuid.UUID      `json:"labId" db:"lab_id" gorm:"type:uuid;primaryKey"`
	TechnologyID uuid.UUID      `json:"technologyId" db:"technology_id" gorm:"type:uuid;primaryKey"`
	CreatedAt    time.Time      `json:"-" db:"created_at"`
	UpdatedAt    time.Time      `json:"-" db:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" db:"deleted_at"`
}

func (db *Database) GetLabDefinitionBySlug(labSlug string) (*Lab, error) {
	var lab Lab
	if err := db.Conn.Where("slug = ?", labSlug).First(&lab).Error; err != nil {
		return nil, fmt.Errorf("failed to get lab by slug: %w", err)
	}
	return &lab, nil
}

func (db *Database) GetLabsDefinitions() ([]Lab, error) {
	var labs []Lab
	if err := db.Conn.Preload("Vulnerabilities").
		Preload("Languages").
		Preload("Technologies").
		Where("active = ?", true).
		Find(&labs).Error; err != nil {
		return nil, fmt.Errorf("failed to get labs definitions: %w", err)
	}
	return labs, nil
}

func (db *Database) GetAllLabsDefinitions() ([]Lab, error) {
	var labs []Lab
	if err := db.Conn.
		Preload("Vulnerabilities").
		Preload("Languages").
		Preload("Technologies").
		Find(&labs).Error; err != nil {
		return nil, fmt.Errorf("failed to get all labs definitions: %w", err)
	}
	return labs, nil
}

func (db *Database) AdminGetLabDefinition(slug string) (*Lab, error) {
	var lab Lab
	if err := db.Conn.Preload("Vulnerabilities").Preload("Languages").Preload("Technologies").Where("slug = ?", slug).First(&lab).Error; err != nil {
		return nil, fmt.Errorf("failed to get lab definition by slug: %w", err)
	}
	return &lab, nil
}
func (db *Database) AdminCreateLabDefinition(lab Lab, vulnerabilties []Vulnerability, languages []Language, technologies []Technology) (*Lab, error) {

	lab.Vulnerabilities = vulnerabilties
	lab.Languages = languages
	lab.Technologies = technologies
	lab.ID = uuid.New()
	if err := db.Conn.Create(&lab).Error; err != nil {
		return nil, fmt.Errorf("failed to save lab with relations: %w", err)
	}
	return &lab, nil
}

func (db *Database) AdminUpdateLabDefinition(lab Lab, vulnerabilities []Vulnerability, languages []Language, technologies []Technology) (*Lab, error) {
	existingLab, err := db.GetLabDefinitionBySlug(lab.Slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing lab definition: %w", err)
	}

	lab.ID = existingLab.ID
	lab.Vulnerabilities = vulnerabilities
	lab.Languages = languages
	lab.Technologies = technologies

	if err := db.Conn.Save(&lab).Error; err != nil {
		return nil, fmt.Errorf("failed to update lab definition: %w", err)
	}

	return &lab, nil
}

func (db *Database) AdminDeleteLabDefinition(labSlug string) error {
	existingLab, err := db.GetLabDefinitionBySlug(labSlug)
	if err != nil {
		return fmt.Errorf("failed to get existing lab definition: %w", err)
	}
	if err := db.Conn.Where("lab_id = ?", existingLab.ID).
		Delete(&LabAttempt{}).Error; err != nil {
		return fmt.Errorf("failed to delete related lab attempts: %w", err)
	}

	if err := db.Conn.Delete(&existingLab).Error; err != nil {
		return fmt.Errorf("failed to delete lab definition: %w", err)
	}

	return nil
}
