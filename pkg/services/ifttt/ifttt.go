package ifttt

import (
	"fmt"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/util/jsonclient"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

const (
	apiURLFormat = "https://maker.ifttt.com/trigger/%s/with/key/%s"
)

// Service sends notifications to a IFTTT webhook
type Service struct {
	standard.Standard
	config *Config
	pkr    format.PropKeyResolver
}

// EmptyConfig returns an empty types.ServiceConfig for the service
func (service *Service) EmptyConfig() types.ServiceConfig {
	return &Config{}
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{
		UseMessageAsValue: 2,
	}
	service.pkr = format.NewPropKeyResolver(service.config)
	if err := service.config.setURL(&service.pkr, configURL); err != nil {
		return err
	}

	return nil
}

// Send a notification message to a IFTTT webhook
func (service *Service) Send(message string, params *types.Params) error {
	config := service.config
	if err := service.pkr.UpdateConfigFromParams(config, params); err != nil {
		return err
	}
	request := createJSONToSend(config, message, params)
	for _, event := range config.Events {
		apiURL := service.createAPIURLForEvent(event)
		if err := jsonclient.Post(apiURL, request, nil); err != nil {
			return fmt.Errorf("failed to send IFTTT event \"%s\": %s", event, err)
		}
	}
	return nil
}

// CreateAPIURLForEvent creates a IFTTT webhook URL for the given event
func (service *Service) createAPIURLForEvent(event string) string {
	return fmt.Sprintf(
		apiURLFormat,
		event,
		service.config.WebHookID,
	)
}
