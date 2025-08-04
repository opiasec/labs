package evaluation

import (
	"appseclabs/types"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type SASTScorer struct {
}

var severityWeights = map[string]int{
	"ERROR":   -10,
	"WARNING": -5,
	"INFO":    -2,
}

func (s *SASTScorer) Score(result string, evaluator types.Evaluator) (types.LabFinishResultCriterion, error) {
	if severityWeightsConfig, ok := evaluator.Config["ERROR"]; ok {
		if severityWeightsConfig != "" {
			intValue, err := strconv.Atoi(severityWeightsConfig)
			if err != nil {
				return types.LabFinishResultCriterion{}, fmt.Errorf("invalid ERROR value: %v", err)
			}
			severityWeights["ERROR"] = intValue
		} else {
			return types.LabFinishResultCriterion{}, fmt.Errorf("ERROR is not defined in the evaluator config")
		}
	}
	if severityWeightsConfig, ok := evaluator.Config["WARNING"]; ok {
		if severityWeightsConfig != "" {
			intValue, err := strconv.Atoi(severityWeightsConfig)
			if err != nil {
				return types.LabFinishResultCriterion{}, fmt.Errorf("invalid WARNING value: %v", err)
			}
			severityWeights["WARNING"] = intValue
		} else {
			return types.LabFinishResultCriterion{}, fmt.Errorf("WARNING is not defined in the evaluator config")
		}
	}
	if severityWeightsConfig, ok := evaluator.Config["INFO"]; ok {
		if severityWeightsConfig != "" {
			intValue, err := strconv.Atoi(severityWeightsConfig)
			if err != nil {
				return types.LabFinishResultCriterion{}, fmt.Errorf("invalid INFO value: %v", err)
			}
			severityWeights["INFO"] = intValue
		} else {
			return types.LabFinishResultCriterion{}, fmt.Errorf("INFO is not defined in the evaluator config")
		}
	}
	var parsed map[string]interface{}
	err := json.Unmarshal([]byte(result), &parsed)
	if err != nil {
		return types.LabFinishResultCriterion{}, fmt.Errorf("error unmarshalling result: %v", err)
	}

	resultsRaw, ok := parsed["results"].([]interface{})
	if !ok {
		return types.LabFinishResultCriterion{}, fmt.Errorf("invalid or missing 'results' field")
	}

	score := 100
	findings := []map[string]interface{}{}
	for _, r := range resultsRaw {
		entry, ok := r.(map[string]interface{})
		if !ok {
			continue
		}

		extra, ok := entry["extra"].(map[string]interface{})
		if !ok {
			continue
		}

		severityRaw, ok := extra["severity"].(string)
		if !ok {
			continue
		}

		severity := strings.ToUpper(severityRaw)
		score += severityWeights[severity]
		findings = append(findings, entry)
	}

	if score < 0 {
		score = 0
	}

	weightedScore := int((float64(score) / 100.0) * float64(evaluator.Weight))

	message := "Nenhuma vulnerabilidade encontrada."
	if len(findings) > 0 {
		message = fmt.Sprintf("Foram encontradas %d vulnerabilidades.", len(findings))
	}

	return types.LabFinishResultCriterion{
		Name:      "SAST Score",
		Score:     weightedScore,
		Status:    "completed",
		Weight:    evaluator.Weight,
		Message:   message,
		RawOutput: result,
	}, nil
}
