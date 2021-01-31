package hangouts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Service providing Hangouts Chat as a notification service.
type Service struct {
	standard.Standard
	config *Config
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service.
func (service *Service) Initialize(configURL *url.URL, logger *log.Logger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{}

	err := service.config.SetURL(configURL)

	return err
}

// Send a notification message to Hangouts Chat.
func (service *Service) Send(message string, _ *types.Params) error {
	config := service.config

	jsonBody, err := json.Marshal(JSON{
		Text: message,
	})

	if err != nil {
		return err
	}

	jsonBuffer := bytes.NewBuffer(jsonBody)
	resp, err := http.Post(config.URL.String(), "application/json", jsonBuffer)

	if err != nil {
		return fmt.Errorf("failed to send notification to Hangouts Chat: %s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Hangouts Chat API notification returned %d HTTP status code", resp.StatusCode)
	}

	return nil
}
