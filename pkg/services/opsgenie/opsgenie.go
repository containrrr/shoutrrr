package opsgenie

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

const (
	alertEndpointTemplate = "https://%s/v2/alerts"
)

// Service providing OpsGenie as a notification service
type Service struct {
	standard.Standard
	config *Config
}

func (service *Service) sendAlert(url string, apiKey string, payload AlertPayload) error {
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	jsonBuffer := bytes.NewBuffer(jsonBody)

	req, err := http.NewRequest("POST", url, jsonBuffer)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "GenieKey "+apiKey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send notification to OpsGenie: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("OpsGenie notification returned %d HTTP status code. Cannot read body: %s", resp.StatusCode, err)
		}
		return fmt.Errorf("OpsGenie notification returned %d HTTP status code: %s", resp.StatusCode, body)
	}

	return nil
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger *log.Logger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{}
	err := service.config.SetURL(configURL)
	return err
}

// Send a notification message to OpsGenie
// See: https://docs.opsgenie.com/docs/alert-api#create-alert
func (service *Service) Send(message string, params *types.Params) error {
	config := service.config
	url := fmt.Sprintf(alertEndpointTemplate, config.Host)
	payload := NewAlertPayload(message, config, params)
	return service.sendAlert(url, config.ApiKey, payload)
}
