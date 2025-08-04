package labdefinitions

import (
	"appseclabs/database"
	"appseclabs/types"
	"log/slog"
)

type LabDefinitionsUsecase struct {
	Database *database.Database
}

func NewLabDefinitionsUsecase(database *database.Database) *LabDefinitionsUsecase {
	return &LabDefinitionsUsecase{Database: database}
}

func (u *LabDefinitionsUsecase) GetAllLabDefinitions() ([]types.Lab, error) {
	labDefinitions, err := u.Database.GetAllActiveLabsDefinitions()
	if err != nil {
		slog.Error("failed to get all lab definitions", "error", err)
		return nil, err
	}
	return labDefinitions, nil
}

func (u *LabDefinitionsUsecase) GetLabDefinition(slug string) (*types.Lab, error) {
	labDefinition, err := u.Database.GetLabDefinitionBySlug(slug)
	if err != nil {
		slog.Error("failed to get lab definition by slug", "slug", slug, "error", err)
		return nil, err
	}
	return labDefinition, nil
}
func (u *LabDefinitionsUsecase) CreateLabDefinition(request CreateLabDefinitionRequest) error {
	lab := &types.Lab{
		Slug: request.Slug,
		LabSpec: types.LabSpec{
			Image: request.LabSpec.Image,
			Env:   request.LabSpec.Env,
			CodeConfig: types.LabCodeConfig{
				GitURL:    request.LabSpec.CodeConfig.GitUrl,
				GitBranch: request.LabSpec.CodeConfig.GitBranch,
				GitPath:   request.LabSpec.CodeConfig.GitPath,
			},
		},
		Evaluators: make([]types.Evaluator, len(request.Evaluators)),
	}
	for i, evaluator := range request.Evaluators {
		lab.Evaluators[i] = types.Evaluator{
			Slug:            evaluator.Slug,
			Weight:          evaluator.Weight,
			ExploitTemplate: evaluator.ExploitTemplate,
			Config:          evaluator.Config,
		}
	}
	err := u.Database.CreateLabDefinition(lab)
	if err != nil {
		slog.Error("failed to create lab definition", "slug", request.Slug, "error", err)
		return err
	}
	return nil
}

func (u *LabDefinitionsUsecase) UpdateLabDefinition(slug string, request UpdateLabDefinitionRequest) error {
	lab := &types.Lab{
		Slug: slug,
		LabSpec: types.LabSpec{
			Image: request.LabSpec.Image,
			Env:   request.LabSpec.Env,
			CodeConfig: types.LabCodeConfig{
				GitURL:    request.LabSpec.CodeConfig.GitUrl,
				GitBranch: request.LabSpec.CodeConfig.GitBranch,
				GitPath:   request.LabSpec.CodeConfig.GitPath,
			},
		},
		Evaluators: make([]types.Evaluator, len(request.Evaluators)),
	}
	for i, evaluator := range request.Evaluators {
		lab.Evaluators[i] = types.Evaluator{
			Slug:            evaluator.Slug,
			Weight:          evaluator.Weight,
			ExploitTemplate: evaluator.ExploitTemplate,
			Config:          evaluator.Config,
		}
	}
	err := u.Database.UpdateLabDefinition(slug, lab)
	if err != nil {
		slog.Error("failed to update lab definition", "slug", slug, "error", err)
		return err
	}
	return nil
}

func (u *LabDefinitionsUsecase) DeleteLabDefinition(slug string) error {
	err := u.Database.DeleteLabDefinition(slug)
	if err != nil {
		slog.Error("failed to delete lab definition", "slug", slug, "error", err)
		return err
	}
	return nil
}

func (u *LabDefinitionsUsecase) GetEvaluators() ([]types.Evaluation, error) {
	evaluators, err := u.Database.GetEvaluators()
	if err != nil {
		slog.Error("failed to get evaluators", "error", err)
		return nil, err
	}
	if len(evaluators) == 0 {
		slog.Info("no evaluators found")
		return nil, nil
	}
	return evaluators, nil
}
