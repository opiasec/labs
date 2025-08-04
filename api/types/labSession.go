package types

import "time"

type LabSession struct {
	LabSlug      string          `bson:"lab_slug" json:"labSlug"`
	Namespace    string          `bson:"namespace" json:"namespace"`
	FinishResult LabFinishResult `bson:"finish_result" json:"finishResult"`
	StartedAt    time.Time       `bson:"started_at" json:"startedAt"`
	FinishedAt   time.Time       `bson:"finished_at" json:"finishedAt"`
	CreatedAt    time.Time       `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time       `bson:"updated_at" json:"updatedAt"`
}

type LabFinishResult struct {
	Status         string                     `bson:"status" json:"status"`
	ErrorMessage   string                     `bson:"error_message" json:"errorMessage"`
	TotalScore     int                        `bson:"total_score" json:"totalScore"`
	CriteriaResult []LabFinishResultCriterion `bson:"criteria_result" json:"criteriaResult"`
	FilesDiff      string                     `bson:"files_diff" json:"filesDiff"`
}

type LabFinishResultCriterion struct {
	Name      string `bson:"name" json:"name"`
	Score     int    `bson:"score" json:"score"`
	Weight    int    `bson:"weight" json:"weight"`
	Required  bool   `bson:"required" json:"required"`
	Message   string `bson:"message" json:"message"`
	Status    string `bson:"status" json:"status"`
	RawOutput string `bson:"raw_output" json:"rawOutput"`
}
