package labdefinition

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type LabDefinitionHandler struct {
	labDefinitionUsecase *LabDefinitionUsecase
}

func NewLabDefinitionHandler(labDefinitionUsecase *LabDefinitionUsecase) *LabDefinitionHandler {
	return &LabDefinitionHandler{labDefinitionUsecase: labDefinitionUsecase}
}

func (h *LabDefinitionHandler) GetLabsDefinitions(c echo.Context) error {
	labDefinitions, err := h.labDefinitionUsecase.GetLabsDefinitions()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error getting lab definitions")
	}

	return c.JSON(http.StatusOK, labDefinitions)
}

func (h *LabDefinitionHandler) GetLabDefinitionBySlug(c echo.Context) error {
	slug := c.Param("slug")
	labDefinition, err := h.labDefinitionUsecase.GetLabDefinitionBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error getting lab definition")
	}

	return c.JSON(http.StatusOK, labDefinition)
}
