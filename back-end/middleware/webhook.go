package middleware

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func WebhookAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		secretToken := c.Request().Header.Get("X-Secret-Token")
		if secretToken != os.Getenv("WEBHOOK_SECRET_TOKEN") {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid secret token")
		}
		return next(c)
	}
}
