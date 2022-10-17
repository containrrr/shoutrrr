package pushbullet

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util/jsonclient"
	"net/url"
)

const (
	pushesEndpoint = "https://api.pushbullet.com/v2/pushes"
)

// Service providing Pushbullet as a notification service
type Service struct {
	standard.Standard
	client jsonclient.Client
	config *Config
	pkr    format.PropKeyResolver
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)

	service.config = &Config{}
	service.pkr = format.NewPropKeyResolver(service.config)
	if err := service.config.setURL(&service.pkr, configURL); err != nil {
		return err
	}

	service.client = jsonclient.NewClient()
	service.client.Headers().Set("Access-Token", service.config.Token)

	return nil
}

// Send a push notification via Pushbullet
func (service *Service) Send(message string, params *types.Params) error {
	config := *service.config
	if err := service.pkr.UpdateConfigFromParams(&config, params); err != nil {
		return err
	}

	for _, target := range config.Targets {
		if err := doSend(&config, target, message, service.client); err != nil {
			return err
		}
	}
	return nil
}

func doSend(config *Config, target string, message string, client jsonclient.Client) error {

	push := NewNotePush(message, config.Title)
	push.SetTarget(target)

	response := PushResponse{}
	if err := client.Post(pushesEndpoint, push, &response); err != nil {
		errorResponse := &ErrorResponse{}
		if client.ErrorResponse(err, errorResponse) {
			return fmt.Errorf("API error: %w", errorResponse)
		}
		return fmt.Errorf("failed to push: %w", err)
	}

	// TODO: Look at response fields?

	return nil
}
