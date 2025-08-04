package auth

import (
	"appseclabsplataform/config"
	"appseclabsplataform/database"
	"appseclabsplataform/utils"
	"errors"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthUsecase struct {
	Database *database.Database
	Config   *config.Config
}

func NewAuthUsecase(database *database.Database, config *config.Config) *AuthUsecase {
	return &AuthUsecase{
		Database: database,
		Config:   config,
	}
}

func (u *AuthUsecase) Login(email, password string) (string, error) {
	if !u.Config.AuthConfig.PasswordEnabled {
		return "", errors.New("password authentication is disabled")
	}

	user, err := u.Database.GetUserByEmail(email)
	if err != nil {
		slog.Error("Failed to get user by email", "error", err)
		return "", err
	}

	if err := utils.CheckPassword(password, user.PasswordHash); err != nil {
		slog.Error("Failed to check password", "error", err)
		return "", err
	}

	token, err := generateJWT(user, u.Config.AuthConfig.JWTSecret)
	if err != nil {
		slog.Error("Failed to generate JWT", "error", err)
		return "", err
	}

	return token, nil
}

func generateJWT(user database.User, secret string) (string, error) {
	permissions := []string{user.Role}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat":         time.Now().Unix(),
		"sub":         user.ID,
		"name":        user.Name,
		"picture":     user.ImageURL,
		"email":       user.Email,
		"exp":         time.Now().Add(24 * time.Hour).Unix(),
		"iss":         "opiaseclabsAPI",
		"permissions": permissions,
	})

	byteSecret := []byte(secret)

	tokenStr, err := token.SignedString(byteSecret)
	if err != nil {
		slog.Error("failed to sign token", "error", err)
		return "", err
	}

	return tokenStr, nil
}
