package lab

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type LabHandler struct {
	labUsecase *LabUsecase
}

func NewLabHandler(labUsecase *LabUsecase) *LabHandler {
	return &LabHandler{
		labUsecase: labUsecase,
	}
}

func (h *LabHandler) CreateLab(c echo.Context) error {
	var req CreateLabRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	userID := c.Get("user_id").(string)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, "Not authorized")
	}

	lab, err := h.labUsecase.CreateLab(req.LabSlug, userID)
	if err != nil {
		if err.Error() == "unauthorized" {
			return c.JSON(http.StatusUnauthorized, "Not authorized")
		}
		return c.JSON(http.StatusInternalServerError, "Error creating lab")
	}

	return c.JSON(http.StatusOK, lab)
}

func (h *LabHandler) GetLabStatus(c echo.Context) error {
	var req GetLabRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Error getting lab")
	}

	userID := c.Get("user_id").(string)

	if userID == "" {
		return c.JSON(http.StatusUnauthorized, "Not authorized")
	}

	lab, err := h.labUsecase.GetLabStatus(req.Namespace, userID)
	if err != nil {
		if err.Error() == "unauthorized" {
			return c.JSON(http.StatusUnauthorized, "Not authorized")
		}
		return c.JSON(http.StatusInternalServerError, "Error getting lab")
	}

	return c.JSON(http.StatusOK, lab)
}

func (h *LabHandler) FinishLab(c echo.Context) error {
	var req FinishLabRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	userID := c.Get("user_id").(string)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, "Not authorized")
	}

	lab, err := h.labUsecase.FinishLab(req.Namespace, userID)
	if err != nil {
		if err.Error() == "unauthorized" {
			return c.JSON(http.StatusUnauthorized, "Not authorized")
		}
		return c.JSON(http.StatusInternalServerError, "Error finishing lab")
	}

	return c.JSON(http.StatusOK, lab)
}

func (h *LabHandler) LeaveLab(c echo.Context) error {
	var req LeaveLabRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	userID := c.Get("user_id").(string)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, "Not authorized")
	}

	lab, err := h.labUsecase.LeaveLab(req.Namespace, userID)
	if err != nil {
		if err.Error() == "unauthorized" {
			return c.JSON(http.StatusUnauthorized, "Not authorized")
		}
		return c.JSON(http.StatusInternalServerError, "Error leaving lab")
	}

	return c.JSON(http.StatusOK, lab)
}

func (h *LabHandler) GetLabResult(c echo.Context) error {
	var req GetLabResultRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	userID := c.Get("user_id").(string)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, "Not authorized")
	}

	labResult, err := h.labUsecase.GetLabResult(req.Namespace, userID)
	if err != nil {
		if err.Error() == "unauthorized" {
			return c.JSON(http.StatusUnauthorized, "Not authorized")
		}
		return c.JSON(http.StatusInternalServerError, "Error getting lab result")
	}

	return c.JSON(http.StatusOK, labResult)
}

func (h *LabHandler) GetAllLabsByUserAndStatus(c echo.Context) error {
	var req GetAllLabsByUserAndStatusRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	userID := c.Get("user_id").(string)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, "Not authorized")
	}

	lab, err := h.labUsecase.GetAllLabsByUserAndStatus(req.Status, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error getting all labs by user and status")
	}

	return c.JSON(http.StatusOK, lab)
}

func (h *LabHandler) SendFeedback(c echo.Context) error {
	var req SendFeedbackRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	userID := c.Get("user_id").(string)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, "Not authorized")
	}

	err := h.labUsecase.SendFeedback(req.Namespace, req.Rating, req.Feedback, userID)
	if err != nil {
		if err.Error() == "unauthorized" {
			return c.JSON(http.StatusUnauthorized, "Not authorized")
		}
		return c.JSON(http.StatusInternalServerError, "Error sending feedback")
	}

	return c.JSON(http.StatusOK, "Feedback sent")
}
