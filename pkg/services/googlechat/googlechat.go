package googlechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Service providing Google Chat as a notification service.
type Service struct {
	standard.Standard
	config *Config
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service.
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{}

	err := service.config.SetURL(configURL)

	return err
}

// Send a notification message to Google Chat.
func (service *Service) Send(message string, _ *types.Params) error {
	config := service.config

	jsonBody, err := json.Marshal(JSON{
		Text: message,
	})
	if err != nil {
		return err
	}

	postURL := getAPIURL(config)

	jsonBuffer := bytes.NewBuffer(jsonBody)
	resp, err := http.Post(postURL.String(), "application/json", jsonBuffer)
	if err != nil {
		return fmt.Errorf("failed to send notification to Google Chat: %s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Google Chat API notification returned %d HTTP status code", resp.StatusCode)
	}

	return nil
}

func getAPIURL(config *Config) *url.URL {
	query := url.Values{}
	query.Set("key", config.Key)
	query.Set("token", config.Token)

	return &url.URL{
		Path:     config.Path,
		Host:     config.Host,
		Scheme:   "https",
		RawQuery: query.Encode(),
	}
}
