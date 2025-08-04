package admin

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	adminUsecase *AdminUsecase
}

func NewAdminHandler(adminUsecase *AdminUsecase) *AdminHandler {
	return &AdminHandler{adminUsecase: adminUsecase}
}

func (h *AdminHandler) GetLabsSessions(c echo.Context) error {
	data, err := h.adminUsecase.GetLabsSessions()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, data)
}

func (h *AdminHandler) GetLabSession(c echo.Context) error {
	namespace := c.Param("namespace")
	data, err := h.adminUsecase.GetLabSession(namespace)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, data)
}

func (h *AdminHandler) GetPossiblesStatus(c echo.Context) error {
	statusFrom := c.QueryParam("from")
	data, err := h.adminUsecase.GetPossiblesStatus(statusFrom)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, data)
}

func (h *AdminHandler) ChangeLabStatus(c echo.Context) error {
	namespace := c.Param("namespace")
	var request ChangeLabStatusRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := h.adminUsecase.ChangeLabStatus(namespace, request.StatusID, request.Comment); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *AdminHandler) GetLabsDefinitions(c echo.Context) error {
	labsDefinitions, err := h.adminUsecase.GetLabsDefinitions()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, labsDefinitions)
}

func (h *AdminHandler) GetLabDefinition(c echo.Context) error {
	slug := c.Param("slug")
	labDefinition, err := h.adminUsecase.GetLabDefinition(slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, labDefinition)
}

func (h *AdminHandler) CreateLabDefinition(c echo.Context) error {
	var request CreateLabDefinitionRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	lab, err := h.adminUsecase.CreateLabDefinition(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, lab)
}

func (h *AdminHandler) UpdateLabDefinition(c echo.Context) error {
	var request UpdateLabDefinitionRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	slug := c.Param("slug")
	lab, err := h.adminUsecase.UpdateLabDefinition(slug, request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, lab)
}

func (h *AdminHandler) DeleteLabDefinition(c echo.Context) error {
	slug := c.Param("slug")
	if err := h.adminUsecase.DeleteLabDefinition(slug); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *AdminHandler) GetPossiblesVulnerabilities(c echo.Context) error {
	vulnerabilities, err := h.adminUsecase.GetPossiblesVulnerabilities()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, vulnerabilities)
}
func (h *AdminHandler) GetPossiblesLanguages(c echo.Context) error {
	languages, err := h.adminUsecase.GetPossiblesLanguages()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, languages)
}
func (h *AdminHandler) GetPossiblesTechnologies(c echo.Context) error {
	technologies, err := h.adminUsecase.GetPossiblesTechnologies()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, technologies)
}

func (h *AdminHandler) GetPossiblesEvaluators(c echo.Context) error {
	evaluators, err := h.adminUsecase.GetPossiblesEvaluators()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, evaluators)
}

func (h *AdminHandler) GetPossiblesImages(c echo.Context) error {
	images, err := h.adminUsecase.GetPossiblesImages()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, images)
}

func (h *AdminHandler) GetAllUsers(c echo.Context) error {
	users, err := h.adminUsecase.GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, users)
}

func (h *AdminHandler) CreateUser(c echo.Context) error {
	var request CreateUserRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if request.Email == "" || request.Password == "" || request.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email, password, and name are required"})
	}

	user, err := h.adminUsecase.CreateUser(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, user)
}
func (h *AdminHandler) GetUserByID(c echo.Context) error {
	id := c.Param("id")
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID format"})
	}
	user, err := h.adminUsecase.GetUserByID(idUUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}
func (h *AdminHandler) UpdateUser(c echo.Context) error {
	var request UpdateUserRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	id := c.Param("id")
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID format"})
	}
	user, err := h.adminUsecase.UpdateUser(idUUID, request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}
func (h *AdminHandler) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID format"})
	}
	if err := h.adminUsecase.DeleteUser(idUUID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
