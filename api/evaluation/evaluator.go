package evaluation

import "appseclabs/types"

type Evaluator struct {
	Scorers map[string]Scorer
}

type Scorer interface {
	Score(result string, evaluator types.Evaluator) (types.LabFinishResultCriterion, error)
}

func NewEvaluator() *Evaluator {
	scorers := make(map[string]Scorer)
	scorers["sast"] = &SASTScorer{}
	return &Evaluator{
		Scorers: scorers,
	}
}
