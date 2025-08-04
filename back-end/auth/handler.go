package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	usecase *AuthUsecase
}

func NewAuthHandler(usecase *AuthUsecase) *AuthHandler {
	return &AuthHandler{
		usecase: usecase,
	}
}

func (h *AuthHandler) Login(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	token, err := h.usecase.Login(email, password)
	if err != nil {
		if err.Error() == "password authentication is disabled" {
			return echo.NewHTTPError(http.StatusForbidden, "Password authentication is disabled")
		}
		return echo.ErrUnauthorized
	}

	return c.JSON(http.StatusOK, map[string]string{
		"access_token": token,
		"id_token":     token,
	})
}
