package types

import "time"

type Evaluation struct {
	Name           string         `bson:"name"`
	Description    string         `bson:"description"`
	EvaluationSpec EvaluationSpec `bson:"evaluation_spec"`
	Slug           string         `bson:"slug"`
	CreatedAt      time.Time      `bson:"created_at"`
	UpdatedAt      time.Time      `bson:"updated_at"`
}

type EvaluationSpec struct {
	Containers    []ContainerSpec `bson:"containers"`
	InitContainer ContainerSpec   `bson:"init_container"`
	Volumes       []VolumeSpec    `bson:"volumes"`
}
type ContainerSpec struct {
	Name     string         `bson:"name"`
	Image    string         `bson:"image"`
	Env      []ContainerEnv `bson:"env"`
	Args     []string       `bson:"args"`
	Commands []string       `bson:"commands"`
	Volumes  []VolumeSpec   `bson:"volumes"`
}

type ContainerEnv struct {
	Name  string `bson:"name"`
	Value string `bson:"value"`
}
type VolumeSpec struct {
	Name string `bson:"name"`
	Path string `bson:"path"`
}
