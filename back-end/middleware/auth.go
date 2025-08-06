package middleware

import (
	"appseclabsplataform/config"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	UserID       string                 `json:"user_id"`
	UserImageURL string                 `json:"user_image_url"`
	UserName     string                 `json:"user_name"`
	UserEmail    string                 `json:"user_email"`
	Metadata     map[string]interface{} `json:"metadata"`
	Permissions  []string               `json:"permissions"`
}

func Auth(config *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"error": "missing authorization header"})
			}

			if !strings.HasPrefix(token, "Bearer ") {
				return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"error": "invalid token format"})
			}

			claims, err := validateTokenFromJWKS(config.AuthConfig.JWKSURL, token)
			if err == nil {
				c.Set("user_id", claims.Subject)
				c.Set("auth_token", token)
				c.Set("permissions", claims.Permissions)

				return next(c)
			}
			if config.AuthConfig.PasswordEnabled {
				slog.Error("JWT validation failed, falling back to secret validation", "error", err)
				claims, err = validateTokenFromSecret(config.AuthConfig.JWTSecret, token)
				if err == nil {
					c.Set("user_id", claims.Subject)
					c.Set("auth_token", token)
					c.Set("permissions", claims.Permissions)

					return next(c)
				}
			}
			return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
		}
	}
}

func validateTokenFromSecret(secret string, token string) (*CustomClaims, error) {
	if secret == "" {
		slog.Error("JWT secret is not set")
		return nil, echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": "internal server error"})
	}

	// Validate token
	tokenString := strings.TrimPrefix(token, "Bearer ")
	tokenParsed, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
		}
		return []byte(secret), nil
	})
	if err != nil || !tokenParsed.Valid {
		slog.Error("error parsing token", "error", err)
		return nil, echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
	}

	if claims, ok := tokenParsed.Claims.(*CustomClaims); ok && tokenParsed.Valid {
		return claims, nil
	}

	return nil, echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
}

func validateTokenFromJWKS(jwksURL string, token string) (*CustomClaims, error) {
	if jwksURL == "" {
		slog.Error("JWKS URL is not set")
		return nil, echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": "internal server error"})
	}

	tokenString := strings.TrimPrefix(token, "Bearer ")
	tokenSimpleParsed, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if tokenSimpleParsed == nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"error": "invalid token format"})
	}

	if err != nil {
		slog.Error("error parsing token", "error", err)
		return nil, echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
	}

	kid, ok := tokenSimpleParsed.Header["kid"].(string)
	if !ok {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"error": "missing kid in token header"})
	}

	resp, err := http.Get(jwksURL)
	if err != nil {
		slog.Error("error getting JWKS", "error", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": "internal server error"})
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("error reading JWKS", "error", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": "internal server error"})
	}

	var jwks struct {
		Keys []struct {
			Kty string `json:"kty"`
			Kid string `json:"kid"`
			Use string `json:"use"`
			Alg string `json:"alg"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}

	if err := json.Unmarshal(body, &jwks); err != nil {
		slog.Error("error unmarshalling JWKS", "error", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": "internal server error"})
	}

	if len(jwks.Keys) == 0 {
		slog.Error("no keys found in JWKS")
		return nil, echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": "internal server error"})
	}

	var key *struct {
		Kty string `json:"kty"`
		Kid string `json:"kid"`
		Use string `json:"use"`
		Alg string `json:"alg"`
		N   string `json:"n"`
		E   string `json:"e"`
	}

	for _, k := range jwks.Keys {
		if k.Kid == kid {
			key = &k
			break
		}
	}

	if key == nil {
		slog.Error("key not found for kid", "kid", kid)
		return nil, echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"error": "key not found"})
	}
	nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
	if err != nil {
		slog.Error("error decoding n", "error", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": "internal server error"})
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
	if err != nil {
		slog.Error("error decoding e", "error", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": "internal server error"})
	}

	publicKey := &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: int(new(big.Int).SetBytes(eBytes).Int64()),
	}

	tokenParsed, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		slog.Error("error parsing token", "error", err)
		return nil, echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
	}

	if claims, ok := tokenParsed.Claims.(*CustomClaims); ok && tokenParsed.Valid {
		return claims, nil
	}

	return nil, echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
}
