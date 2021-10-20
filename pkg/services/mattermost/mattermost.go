package mattermost

import (
	"fmt"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/util/jsonclient"

	"github.com/containrrr/shoutrrr/pkg/types"
)

// Service sends notifications to a pre-configured channel or user
type Service struct {
	standard.Standard
	config *Config
}

// EmptyConfig returns an empty types.ServiceConfig for the service
func (service *Service) EmptyConfig() types.ServiceConfig {
	return &Config{}
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
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
	request := CreateJSONPayload(config, message, params)
	if err := jsonclient.Post(apiURL, request, nil); err != nil {
		return err
	}
	return nil
}

// Builds the actual URL the request should go to
func buildURL(config *Config) string {
	return fmt.Sprintf("https://%s/hooks/%s", config.Host, config.Token)
}
