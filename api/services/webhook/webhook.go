package webhook

import (
	"appseclabs/types"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type WebhookService struct {
	BaseURL string
	Client  *http.Client
	Enabled bool
	Secret  string
}

func NewWebhookService() *WebhookService {
	return &WebhookService{
		BaseURL: os.Getenv("WEBHOOK_BASE_URL"),
		Enabled: os.Getenv("WEBHOOK_ENABLED") == "true",
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
		Secret: os.Getenv("WEBHOOK_SECRET_TOKEN"),
	}
}
func (w *WebhookService) IsEnabled() bool {
	return w.Enabled
}

func (w *WebhookService) SendFinishEvaluationResult(session types.LabSession) error {
	if !w.IsEnabled() {
		return nil
	}
	url := w.BaseURL + "/webhook/finish-evaluation"
	jsonBody, err := json.Marshal(session)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("X-Secret-Token", w.Secret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := w.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send webhook: %s", resp.Status)
	}

	return nil
}
