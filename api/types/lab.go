package types

import "time"

type Lab struct {
	LabSpec    LabSpec     `bson:"labSpec" json:"labSpec"`
	Slug       string      `bson:"slug" json:"slug"`
	Name       string      `bson:"name" json:"name"`
	CreatedAt  time.Time   `bson:"created_at" json:"createdAt"`
	UpdatedAt  time.Time   `bson:"updated_at" json:"updatedAt"`
	Evaluators []Evaluator `bson:"evaluators,omitempty" json:"evaluators,omitempty"`
}
type Evaluator struct {
	Slug            string            `bson:"slug" json:"slug"`
	Weight          int               `bson:"weight" json:"weight"`
	ExploitTemplate string            `bson:"exploit_template,omitempty" json:"exploitTemplate,omitempty"`
	Config          map[string]string `bson:"config,omitempty" json:"config,omitempty"`
}

type LabSpec struct {
	Image        string            `bson:"image"`
	Ports        []int32           `bson:"ports"`
	Env          map[string]string `bson:"env" json:"envVars"`
	Evaluations  []LabEvaluation   `bson:"evaluations"`
	Args         []string          `bson:"args"`
	Services     []LabService      `bson:"services"`
	Healthcheck  LabHealthcheck    `bson:"healthcheck"`
	FinishConfig LabFinishConfig   `bson:"finish_config" json:"finishConfig"`
	CodeConfig   LabCodeConfig     `bson:"code_config" json:"codeConfig"`
}

type LabEvaluation struct {
	Name   string               `bson:"name"`
	Order  int                  `bson:"order"`
	Params map[string]string    `bson:"params"`
	Envs   []EnviromentVariable `bson:"envs"`
	Weight int                  `bson:"weight"`
}

type EnviromentVariable struct {
	Name  string `bson:"name"`
	Value string `bson:"value"`
}

type LabCodeConfig struct {
	GitURL    string `bson:"git_url" json:"gitUrl"`
	GitBranch string `bson:"git_branch" json:"gitBranch"`
	GitPath   string `bson:"git_path" json:"gitPath"`
}

type LabHealthcheck struct {
	Path string `bson:"path"`
	Port int32  `bson:"port"`
}

type LabService struct {
	Name string `bson:"name"`
	Port int32  `bson:"port"`
	Path string `bson:"path"`
}

type LabFinishConfig struct {
	Criteria []Criterion `bson:"criteria" json:"criteria"`
}

type Criterion struct {
	Name     string `bson:"name" json:"name"`
	Script   string `bson:"script" json:"script"`
	Weight   int    `bson:"weight" json:"weight"`
	Order    int    `bson:"order" json:"order"`
	Required bool   `bson:"required" json:"required"`
}
