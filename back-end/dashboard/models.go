package dashboard

import "time"

type DashboardData struct {
	WeeklyActivity    []WeeklyActivity `json:"weeklyActivity"`
	ScoreEvolution    []ScoreEvolution `json:"scoreEvolution"`
	AvailableLabs     int              `json:"availableLabs"`
	CompletedLabs     int              `json:"completedLabs"`
	TotalAttemptsLabs int              `json:"totalAttemptsLabs"`
	OnReviewLabs      int              `json:"onReviewLabs"`
	AverageScore      int              `json:"averageScore"`
}

type WeeklyActivity struct {
	Day   time.Time `json:"day"`
	Count int       `json:"count"`
}

type ScoreEvolution struct {
	Day   time.Time `json:"day"`
	Score int       `json:"score"`
}
