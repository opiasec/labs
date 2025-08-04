package labdefinitions

type CreateLabDefinitionRequest struct {
	Slug       string                                `json:"slug" bson:"slug"`
	LabSpec    CreateLabDefinitionLabSpecRequest     `json:"labSpec" bson:"labSpec"`
	Evaluators []CreateLabDefinitionRequestEvaluator `json:"evaluators,omitempty" bson:"evaluators,omitempty"`
}

type CreateLabDefinitionLabSpecRequest struct {
	Image      string                                      `json:"image" bson:"image"`
	Env        map[string]string                           `json:"envVars" bson:"env"`
	CodeConfig CreateLabDefinitionLabSpecCodeConfigRequest `json:"codeConfig" bson:"code_config"`
}

type CreateLabDefinitionLabSpecCodeConfigRequest struct {
	GitUrl    string `json:"gitUrl" bson:"git_url"`
	GitBranch string `json:"gitBranch" bson:"git_branch"`
	GitPath   string `json:"gitPath" bson:"git_path"`
}

type CreateLabDefinitionRequestEvaluator struct {
	Slug            string            `json:"slug" bson:"slug"`
	Weight          int               `json:"weight" bson:"weight"`
	ExploitTemplate string            `json:"exploitTemplate,omitempty" bson:"exploit_template, omitempty"`
	Config          map[string]string `json:"config,omitempty" bson:"config, omitempty"`
}

type UpdateLabDefinitionRequest struct {
	CreateLabDefinitionRequest
}
