package rocketchat

import (
	"fmt"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util/jsonclient"
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

// Send a notification message to Rocket.chat
func (service *Service) Send(message string, params *types.Params) error {
	config := service.config
	apiURL := buildURL(config)
	request := CreateJSONPayload(config, message, params)
	if err := jsonclient.Post(apiURL, request, nil); err != nil {
		return fmt.Errorf("Error while posting to URL: %v\nHOST: %s\nPORT: %s", err, config.Host, config.Port)
	}
	return nil
}

func buildURL(config *Config) string {
	if config.Port != "" {
		return fmt.Sprintf("https://%s:%s/hooks/%s/%s", config.Host, config.Port, config.TokenA, config.TokenB)
	}

	return fmt.Sprintf("https://%s/hooks/%s/%s", config.Host, config.TokenA, config.TokenB)
}
