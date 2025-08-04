package labcluster

import (
	"appseclabsplataform/config"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
)

type LabClusterService struct {
	Client  *http.Client
	BaseURL string
}

func NewLabClusterService(config *config.Config) *LabClusterService {
	return &LabClusterService{
		Client: &http.Client{
			Timeout: config.LabClusterServiceConfig.Timeout,
		},
		BaseURL: config.LabClusterServiceConfig.BaseURL,
	}
}

func (l *LabClusterService) CreateLab(labSlug string) (*CreateLabResponse, error) {
	body := CreateLabRequest{
		LabSlug: labSlug,
	}

	url := fmt.Sprintf("%s/api/labs/", l.BaseURL)
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := l.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to create lab: %s", resp.Status)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var labResponse CreateLabResponse
	err = json.Unmarshal(responseBody, &labResponse)
	if err != nil {
		return nil, err
	}

	return &labResponse, nil
}

func (l *LabClusterService) GetLabStatus(namespace string, userToken string) (*GetLabResponse, error) {

	url := fmt.Sprintf("%s/api/labs/%s/status", l.BaseURL, namespace)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", userToken)
	resp, err := l.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, errors.New("unauthorized")
	}

	var labResponse GetLabResponse
	err = json.NewDecoder(resp.Body).Decode(&labResponse)
	if err != nil {
		return nil, err
	}

	return &labResponse, nil
}

func (l *LabClusterService) FinishLab(namespace string, labSlug string) (*FinishLabResponse, error) {
	body := LabFinishRequest{
		LabSlug: labSlug,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/api/labs/%s/finish", l.BaseURL, namespace)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := l.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to finish lab: %s", resp.Status)
	}
	defer resp.Body.Close()

	var finishLabResponse FinishLabResponse
	err = json.NewDecoder(resp.Body).Decode(&finishLabResponse)
	if err != nil {
		return nil, err
	}

	return &finishLabResponse, nil
}

func (l *LabClusterService) GetLabDefinitions() ([]*LabDefinition, error) {
	url := fmt.Sprintf("%s/api/lab-definitions/", l.BaseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := l.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, errors.New("unauthorized")
	}
	var labDefinitions []*LabDefinition
	err = json.NewDecoder(resp.Body).Decode(&labDefinitions)
	if err != nil {
		return nil, err
	}

	return labDefinitions, nil
}

func (l *LabClusterService) GetLabResult(namespace string) (*GetLabResultResponse, error) {
	url := fmt.Sprintf("%s/api/labs/%s", l.BaseURL, namespace)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := l.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, errors.New("unauthorized")
	}

	var labResult GetLabResultResponse
	err = json.NewDecoder(resp.Body).Decode(&labResult)
	if err != nil {
		return nil, err
	}

	return &labResult, nil
}

func (l *LabClusterService) GetAllLabsByUserAndStatus(status string, userToken string) ([]*GetAllLabsByUserAndStatusResponse, error) {

	url := fmt.Sprintf("%s/api/labs/?status=%s", l.BaseURL, status)
	log.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", userToken)
	resp, err := l.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, errors.New("unauthorized")
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("not found")
	}
	if resp.StatusCode == http.StatusBadRequest {
		return nil, errors.New("bad request")
	}

	var labs []*GetAllLabsByUserAndStatusResponse
	err = json.NewDecoder(resp.Body).Decode(&labs)
	if err != nil {
		return nil, err
	}

	return labs, nil
}

func (l *LabClusterService) DeleteLabSession(namespace string) error {

	url := fmt.Sprintf("%s/api/labs/%s", l.BaseURL, namespace)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	resp, err := l.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete lab session: %s", resp.Status)
	}

	return nil
}

func (l *LabClusterService) SendFeedback(namespace string, userToken string, rating int, feedback string) error {

	body := SendFeedbackRequest{
		Rating:   rating,
		Feedback: feedback,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/labs/%s/feedback", l.BaseURL, namespace)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", userToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := l.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return errors.New("unauthorized")
	}

	if resp.StatusCode == http.StatusBadRequest {
		return errors.New("bad request")
	}

	if resp.StatusCode == http.StatusNotFound {
		return errors.New("not found")
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return errors.New("internal server error")
	}

	return nil
}

func (l *LabClusterService) GetLabDefinitionBySlug(slug string) (*LabDefinition, error) {
	url := fmt.Sprintf("%s/api/lab-definitions/%s", l.BaseURL, slug)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := l.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, errors.New("unauthorized")
	}

	var labDefinition LabDefinition
	err = json.NewDecoder(resp.Body).Decode(&labDefinition)
	if err != nil {
		return nil, err
	}

	return &labDefinition, nil
}

func (l *LabClusterService) CreateLabDefinition(definition LabDefinition) error {
	url := fmt.Sprintf("%s/api/lab-definitions/", l.BaseURL)
	slog.Info("Creating lab definition", "url", url, "definition", definition)
	jsonBody, err := json.Marshal(definition)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := l.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create lab definition: %s", resp.Status)
	}

	return nil
}

func (l *LabClusterService) UpdateLabDefinition(slug string, definition LabDefinition) error {
	url := fmt.Sprintf("%s/api/lab-definitions/%s", l.BaseURL, slug)
	jsonBody, err := json.Marshal(definition)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := l.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update lab definition: %s", resp.Status)
	}

	return nil
}

func (l *LabClusterService) DeleteLabDefinition(slug string) error {
	url := fmt.Sprintf("%s/api/lab-definitions/%s", l.BaseURL, slug)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := l.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete lab definition: %s", resp.Status)
	}

	return nil
}

func (l *LabClusterService) GetEvaluators() ([]Evaluator, error) {
	url := fmt.Sprintf("%s/api/evaluator/", l.BaseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := l.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, errors.New("unauthorized")
	}

	var evaluators []Evaluator
	err = json.NewDecoder(resp.Body).Decode(&evaluators)
	if err != nil {
		return nil, err
	}

	return evaluators, nil
}
