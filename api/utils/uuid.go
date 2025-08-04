package utils

import (
	"crypto/rand"
	"fmt"

	"github.com/google/uuid"
)

func GenerateUUID() string {
	return uuid.New().String()
}

func GenerateRandomPassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-"

	if length <= 0 {
		return "", fmt.Errorf("password length must be greater than zero")
	}

	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	for i := 0; i < length; i++ {
		randomBytes[i] = charset[int(randomBytes[i])%len(charset)]
	}
	return string(randomBytes), nil
}
