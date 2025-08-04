package config

import (
	"errors"
	"os"
	"strings"
	"time"
)

type AuthConfig struct {
	JWTSecret       string
	PasswordEnabled bool
	JWKSURL         string
	RootPassword    string
}

type DatabaseConfig struct {
	ConString string
}

type LabSessionServiceConfig struct {
	Addr     string
	Password string
	DB       string
}

type LabClusterServiceConfig struct {
	BaseURL string
	Timeout time.Duration
}

type LabIDEServiceConfig struct {
	BaseURL string
	Timeout time.Duration
}

type TLSConfig struct {
	Enabled  bool
	CertFile string
	KeyFile  string
	Port     string
}

type Config struct {
	AuthConfig              AuthConfig
	DatabaseConfig          DatabaseConfig
	LabSessionServiceConfig LabSessionServiceConfig
	LabClusterServiceConfig LabClusterServiceConfig
	LabIDEServiceConfig     LabIDEServiceConfig
	TLSConfig               TLSConfig
	PublicConfig            PublicConfig `json:"config"`
}

func NewConfig() *Config {
	return &Config{
		AuthConfig: AuthConfig{
			JWTSecret:       os.Getenv("AUTH_JWT_SECRET"),
			PasswordEnabled: getEnvOrDefault("AUTH_PASSWORD_ENABLED", "true") == "true",
			JWKSURL:         os.Getenv("JWKS_URL"),
			RootPassword:    os.Getenv("ROOT_PASSWORD"),
		},
		DatabaseConfig: DatabaseConfig{
			ConString: os.Getenv("DB_POSTGRES_URL"),
		},
		LabSessionServiceConfig: LabSessionServiceConfig{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       os.Getenv("REDIS_DB"),
		},
		LabClusterServiceConfig: LabClusterServiceConfig{
			BaseURL: os.Getenv("LAB_CLUSTER_BASE_URL"),
			Timeout: 10 * time.Second,
		},
		LabIDEServiceConfig: LabIDEServiceConfig{
			BaseURL: os.Getenv("LAB_IDE_BASE_URL"),
			Timeout: 10 * time.Second,
		},
		TLSConfig: TLSConfig{
			Enabled:  os.Getenv("TLS_ENABLED") == "true",
			CertFile: os.Getenv("TLS_CERT_FILE"),
			KeyFile:  os.Getenv("TLS_KEY_FILE"),
			Port:     getEnvOrDefault("TLS_PORT", "443"),
		},
		PublicConfig: PublicConfig{
			PasswordAuthEnabled:  getEnvOrDefault("AUTH_PASSWORD_ENABLED", "true") == "true",
			ProviderLoginEnabled: getEnvOrDefault("AUTH_PROVIDER_LOGIN_ENABLED", "false") == "true",
			ProviderConfig: ProviderConfig{
				Name:             os.Getenv("AUTH_PROVIDER_NAME"),
				AuthorizationURL: os.Getenv("AUTH_PROVIDER_AUTHORIZATION_URL"),
				TokenURL:         os.Getenv("AUTH_PROVIDER_TOKEN_URL"),
				Audience:         os.Getenv("AUTH_PROVIDER_AUDIENCE"),
				ClientID:         os.Getenv("AUTH_PROVIDER_CLIENT_ID"),
				Scope:            getEnvOrDefault("AUTH_PROVIDER_SCOPE", "openid profile email"),
			},
		},
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func CheckEnvironmentVariables() error {
	requiredVars := []string{
		"DB_POSTGRES_URL",
		"REDIS_ADDR",
		"REDIS_DB",
		"LAB_CLUSTER_BASE_URL",
		"LAB_IDE_BASE_URL",
	}
	var missingVars []string
	for _, envVar := range requiredVars {
		if os.Getenv(envVar) == "" {
			missingVars = append(missingVars, envVar)
		}
	}

	if len(missingVars) > 0 {
		return errors.New("missing required environment variables: " + strings.Join(missingVars, ", "))
	}

	return nil
}
