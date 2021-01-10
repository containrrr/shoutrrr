package mattermost

import (
	"bytes"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"log"
	"net/http"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/types"
)

// Service sends notifications to a pre-configured channel or user
type Service struct {
	standard.Standard
	config *Config
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger *log.Logger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{}
	if err := service.config.SetURL(configURL); err != nil {
		return err
	}

	return nil
}

// Send a notification message to Mattermost
func (service *Service) Send(message string, params *types.Params) error {
	config := service.config
	apiURL := buildURL(config)
	json, _ := CreateJSONPayload(config, message, params)
	res, err := http.Post(apiURL, "application/json", bytes.NewReader(json))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send notification to service, response status code %s", res.Status)
	}
	return err
}

// SendItems concatenates the items and sends them using Send
func (service *Service) SendItems(items []types.MessageItem, params *types.Params) error {
	return service.Send(types.ItemsToPlain(items), params)
}

// Builds the actual URL the request should go to
func buildURL(config *Config) string {
	return fmt.Sprintf("https://%s/hooks/%s", config.Host, config.Token)
}
