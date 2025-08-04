package labdefinition

import (
	"appseclabsplataform/database"
)

type GetLabDefinitionBySlugResponse struct {
	database.Lab
}

type GetLabsDefinitionsResponse struct {
	database.Lab
	Readme string `json:"readme,omitempty"`
}
