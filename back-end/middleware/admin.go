package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Admin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userRole := c.Get("permissions")

		if userRole == nil {
			return echo.NewHTTPError(http.StatusForbidden, "User role not found")
		}

		roles, ok := userRole.([]string)
		if !ok {
			return echo.NewHTTPError(http.StatusForbidden, "Invalid user role format")
		}

		for _, role := range roles {
			if role == "admin" {
				return next(c)
			}
		}

		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to access this resource")
	}
}
