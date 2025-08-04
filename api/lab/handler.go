package lab

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type LabHandler struct {
	LabUsecase *LabUsecase
}

func NewLabHandler(labUsecase *LabUsecase) *LabHandler {
	return &LabHandler{LabUsecase: labUsecase}
}

func (h *LabHandler) CreateLabHandler(c echo.Context) error {
	var req CreateLabRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	lab, err := h.LabUsecase.CreateLab(req.LabSlug)
	if err != nil {
		if err.Error() == "lab not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "lab not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create lab"})
	}

	return c.JSON(http.StatusOK, lab)

}

func (h *LabHandler) DeleteLabHandler(c echo.Context) error {
	var req DeleteLabRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	err := h.LabUsecase.DeleteLab(req.Namespace)
	if err != nil {
		if err.Error() == "unauthorized" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete lab"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "lab deleted"})
}

func (h *LabHandler) FinishLabHandler(c echo.Context) error {
	var req FinishLabRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	finishLabResponse, err := h.LabUsecase.FinishLab(req.Namespace, req.LabSlug)
	if err != nil {
		if err.Error() == "unauthorized" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to finish lab"})
	}

	return c.JSON(http.StatusOK, finishLabResponse)
}

func (h *LabHandler) GetLabResultHandler(c echo.Context) error {
	var req GetLabResultRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	labResult, err := h.LabUsecase.GetLabResult(req.Namespace)
	if err != nil {
		if err.Error() == "unauthorized" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get lab result"})
	}
	return c.JSON(http.StatusOK, labResult)
}
