package labide

import (
	"appseclabsplataform/config"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type LabIDEService struct {
	Client  *http.Client
	BaseURL string
}

func NewLabIDEService(config *config.Config) *LabIDEService {
	return &LabIDEService{
		Client: &http.Client{
			Timeout: config.LabIDEServiceConfig.Timeout,
		},
		BaseURL: config.LabIDEServiceConfig.BaseURL,
	}
}

func (l *LabIDEService) GetStatus(ideLabURL string) (GetLabResponse, error) {
	slog.Info("Checking IDE lab status", "url", ideLabURL)
	resp, err := l.Client.Get(fmt.Sprintf("%shealthz", ideLabURL))
	if err != nil {
		return GetLabResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GetLabResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var getLabResponse GetLabResponse
	if err := json.NewDecoder(resp.Body).Decode(&getLabResponse); err != nil {
		return GetLabResponse{}, fmt.Errorf("failed to decode response: %v", err)
	}

	return getLabResponse, nil
}
