package dashboard

import (
	"appseclabsplataform/database"
	"time"

	"github.com/google/uuid"
)

type DashboardUsecase struct {
	database *database.Database
}

func NewDashboardUsecase(database *database.Database) *DashboardUsecase {
	return &DashboardUsecase{
		database: database,
	}
}

func (u *DashboardUsecase) GetDashboardData(userID string) (*DashboardData, error) {

	var completedLabs, onReviewLabs, averageScore int
	var weeklyActivity []WeeklyActivity
	var scoreEvolution []ScoreEvolution

	startDate := time.Now().AddDate(0, 0, -6)
	endDate := time.Now()
	weeklyLabSessions, err := u.database.GetLabAttemptsByUserFilterByDay(userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	attemptsByDay := make(map[string]int)
	averageScoreByDay := make(map[string]int)

	for i := 0; i < 7; i++ {
		day := startDate.AddDate(0, 0, i)
		dayKey := day.Format("2006-01-02")
		attemptsByDay[dayKey] = 0
		averageScoreByDay[dayKey] = 0
	}
	for _, session := range weeklyLabSessions {
		dayKey := session.CreatedAt.Format("2006-01-02")
		attemptsByDay[dayKey]++
		averageScoreByDay[dayKey] += session.Score
	}
	for i := 0; i < 7; i++ {
		day := startDate.AddDate(0, 0, i)
		dayKey := day.Format("2006-01-02")
		weeklyActivity = append(weeklyActivity, WeeklyActivity{
			Day:   day,
			Count: attemptsByDay[dayKey],
		})
		if attemptsByDay[dayKey] > 0 {
			averageScoreByDay[dayKey] /= attemptsByDay[dayKey]
		}
		scoreEvolution = append(scoreEvolution, ScoreEvolution{
			Day:   day,
			Score: averageScoreByDay[dayKey],
		})
	}

	allLabs, err := u.database.GetLabAttemptsByUser(userID)
	if err != nil {
		return nil, err
	}

	availableLabs, err := u.database.GetLabsDefinitions()
	if err != nil {
		return nil, err
	}

	var mapAvailableLabs = make(map[uuid.UUID]bool)
	for _, lab := range availableLabs {
		mapAvailableLabs[lab.ID] = false
	}

	for _, lab := range allLabs {
		switch lab.Status.Name {
		case "approved", "passed":
			mapAvailableLabs[lab.LabID] = true
		case "pending_review":
			onReviewLabs++
		}
		if lab.Score > 0 {
			averageScore += lab.Score
		}
	}

	for _, lab := range mapAvailableLabs {
		if lab {
			completedLabs++
		}
	}

	if len(allLabs) > 0 {
		averageScore /= len(allLabs)
	} else {
		averageScore = 0
	}

	return &DashboardData{
		WeeklyActivity:    weeklyActivity,
		ScoreEvolution:    scoreEvolution,
		TotalAttemptsLabs: len(allLabs),
		AvailableLabs:     len(availableLabs),
		CompletedLabs:     completedLabs,
		OnReviewLabs:      onReviewLabs,
		AverageScore:      averageScore,
	}, nil
}
