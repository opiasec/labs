package dashboard

import (
	"github.com/labstack/echo/v4"
)

type DashboardHandler struct {
	dashboardUsecase *DashboardUsecase
}

func NewDashboardHandler(dashboardUsecase *DashboardUsecase) *DashboardHandler {
	return &DashboardHandler{
		dashboardUsecase: dashboardUsecase,
	}
}
func (h *DashboardHandler) GetDashboardData(c echo.Context) error {
	userID := c.Get("user_id").(string)
	if userID == "" {
		return c.JSON(400, map[string]string{"error": "userID is required"})
	}
	data, err := h.dashboardUsecase.GetDashboardData(userID)
	if err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}
	return c.JSON(200, data)
}
