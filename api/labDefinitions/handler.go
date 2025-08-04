package labdefinitions

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type LabDefinitionsHandler struct {
	LabDefinitionsUsecase *LabDefinitionsUsecase
}

func NewLabDefinitionsHandler(labDefinitionsUsecase *LabDefinitionsUsecase) *LabDefinitionsHandler {
	return &LabDefinitionsHandler{LabDefinitionsUsecase: labDefinitionsUsecase}
}

func (h *LabDefinitionsHandler) GetAllLabDefinitionsHandler(c echo.Context) error {
	labDefinitions, err := h.LabDefinitionsUsecase.GetAllLabDefinitions()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to get lab definitions"})
	}
	return c.JSON(http.StatusOK, labDefinitions)
}

func (h *LabDefinitionsHandler) CreateLabDefinitionHandler(c echo.Context) error {
	var request CreateLabDefinitionRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}
	err := h.LabDefinitionsUsecase.CreateLabDefinition(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create lab definition"})
	}
	return c.JSON(http.StatusCreated, echo.Map{"message": "lab definition created successfully"})
}

func (h *LabDefinitionsHandler) GetLabDefinitionHandler(c echo.Context) error {
	slug := c.Param("slug")
	labDefinition, err := h.LabDefinitionsUsecase.GetLabDefinition(slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to get lab definition"})
	}
	if labDefinition == nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "lab definition not found"})
	}
	return c.JSON(http.StatusOK, labDefinition)
}

func (h *LabDefinitionsHandler) UpdateLabDefinitionHandler(c echo.Context) error {
	slug := c.Param("slug")
	var request UpdateLabDefinitionRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}
	err := h.LabDefinitionsUsecase.UpdateLabDefinition(slug, request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to update lab definition"})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *LabDefinitionsHandler) DeleteLabDefinitionHandler(c echo.Context) error {
	slug := c.Param("slug")
	err := h.LabDefinitionsUsecase.DeleteLabDefinition(slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to delete lab definition"})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *LabDefinitionsHandler) GetEvaluators(c echo.Context) error {

	evaluators, err := h.LabDefinitionsUsecase.GetEvaluators()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to get evaluators"})
	}
	if len(evaluators) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "no evaluators found for this lab"})
	}
	return c.JSON(http.StatusOK, evaluators)
}
