package database

import (
	"appseclabsplataform/config"
	"appseclabsplataform/utils"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Email         string         `gorm:"type:text;uniqueIndex;not null" json:"email"`
	PasswordHash  string         `gorm:"type:text;not null" json:"-"`
	Name          string         `gorm:"type:text;not null" json:"name"`
	ImageURL      string         `gorm:"type:text" json:"image_url"`
	Role          string         `gorm:"type:text;default:'user';check:role IN ('user', 'admin', 'moderator')" json:"role"`
	IsActive      bool           `gorm:"default:true" json:"is_active"`
	EmailVerified bool           `gorm:"default:false" json:"email_verified"`
	LastLoginAt   time.Time      `gorm:"type:timestamp" json:"last_login_at"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (db *Database) CreateRootAccount(authConfig config.AuthConfig) error {
	password := authConfig.RootPassword
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	var existingUser User
	if err := db.Conn.Where("email = ?", "root@root.com").First(&existingUser).Error; err == nil {
		return nil
	}

	user := User{
		Email:        "root@root.com",
		PasswordHash: passwordHash,
		Name:         "Root User",
		Role:         "admin",
	}

	if err := db.Conn.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (db *Database) GetAllUsers() ([]User, error) {
	var users []User
	if err := db.Conn.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	return users, nil
}

func (db *Database) GetUserByEmail(email string) (User, error) {
	var user User
	if err := db.Conn.Where("email = ?", email).First(&user).Error; err != nil {
		return User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}
func (db *Database) GetUserByID(id uuid.UUID) (User, error) {
	var user User
	if err := db.Conn.Where("id = ?", id).First(&user).Error; err != nil {
		return User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}
func (db *Database) CreateUser(user User) (User, error) {
	user.ID = uuid.New()
	if err := db.Conn.Create(&user).Error; err != nil {
		return User{}, fmt.Errorf("failed to save user: %w", err)
	}
	return user, nil

}

func (db *Database) UpdateUser(user User) (User, error) {

	var existingUser User
	if err := db.Conn.Where("id = ?", user.ID).First(&existingUser).Error; err != nil {
		return User{}, fmt.Errorf("failed to get user: %w", err)
	}

	user.PasswordHash = existingUser.PasswordHash

	if err := db.Conn.Save(&user).Error; err != nil {
		return User{}, fmt.Errorf("failed to update user: %w", err)
	}
	return user, nil
}

func (db *Database) DeleteUser(id uuid.UUID) error {
	var user User
	if err := db.Conn.Where("id = ?", id).First(&user).Error; err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if err := db.Conn.Delete(&user).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
