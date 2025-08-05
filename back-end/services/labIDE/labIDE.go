package labide

import (
	"appseclabsplataform/config"
	"encoding/json"
	"fmt"
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

func (l *LabIDEService) GetStatus(namespace string) (GetLabResponse, error) {
	resp, err := l.Client.Get(fmt.Sprintf("%s/%s/healthz", l.BaseURL, namespace))
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
