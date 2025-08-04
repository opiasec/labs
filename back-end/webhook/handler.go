package webhook

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type WebHookHandler struct {
	webhookUsecase *WebHookUsecase
}

func NewWebHookHandler(webhookUsecase *WebHookUsecase) *WebHookHandler {
	return &WebHookHandler{
		webhookUsecase: webhookUsecase,
	}
}

func (h *WebHookHandler) FinishEvaluationResult(c echo.Context) error {

	var sessionWebHook LabSession
	if err := c.Bind(&sessionWebHook); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := h.webhookUsecase.FinishEvaluationResult(sessionWebHook); err != nil {
		if err.Error() == "unauthorized" {
			return c.JSON(http.StatusUnauthorized, "Unauthorized")
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "Evaluation result processed successfully")
}
