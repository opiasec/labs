package labdefinition

import (
	"appseclabsplataform/database"
	labcluster "appseclabsplataform/services/labCluster"
	"log/slog"
)

type LabDefinitionUsecase struct {
	labCluster *labcluster.LabClusterService
	Database   *database.Database
}

func NewLabDefinitionUsecase(labCluster *labcluster.LabClusterService, database *database.Database) *LabDefinitionUsecase {
	return &LabDefinitionUsecase{labCluster: labCluster, Database: database}
}

func (u *LabDefinitionUsecase) GetLabsDefinitions() ([]GetLabsDefinitionsResponse, error) {
	labsDefinitions, err := u.Database.GetLabsDefinitions()
	if err != nil {
		slog.Error("failed to get lab definitions", "error", err)
		return nil, err
	}

	var responses []GetLabsDefinitionsResponse
	for _, lab := range labsDefinitions {
		response := GetLabsDefinitionsResponse{
			Lab:    lab,
			Readme: "",
		}
		responses = append(responses, response)
	}
	return responses, nil
}

func (u *LabDefinitionUsecase) GetLabDefinitionBySlug(slug string) (*GetLabDefinitionBySlugResponse, error) {
	labDefinition, err := u.Database.GetLabDefinitionBySlug(slug)
	if err != nil {
		slog.Error("failed to get lab definition by slug", "slug", slug, "error", err)
		return nil, err
	}

	response := GetLabDefinitionBySlugResponse{
		Lab: database.Lab(*labDefinition),
	}
	return &response, nil
}
