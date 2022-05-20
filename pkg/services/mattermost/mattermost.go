package mattermost

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/standard"

	"github.com/containrrr/shoutrrr/pkg/types"
)

// Service sends notifications to a pre-configured channel or user
type Service struct {
	standard.Standard
	config *Config
	pkr    format.PropKeyResolver
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{}
	service.pkr = format.NewPropKeyResolver(service.config)
	return service.config.setURL(&service.pkr, configURL)
}

// Send a notification message to Mattermost
func (service *Service) Send(message string, params *types.Params) error {
	config := service.config
	apiURL := buildURL(config)

	if err := service.pkr.UpdateConfigFromParams(config, params); err != nil {
		return err
	}
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

// Builds the actual URL the request should go to
func buildURL(config *Config) string {
	return fmt.Sprintf("https://%s/hooks/%s", config.Host, config.Token)
}
