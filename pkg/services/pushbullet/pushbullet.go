package pushbullet

import (
	"fmt"
	"net/url"

	"github.com/nicholas-fedor/shoutrrr/pkg/format"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/standard"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
	"github.com/nicholas-fedor/shoutrrr/pkg/util/jsonclient"
)

const (
	pushesEndpoint = "https://api.pushbullet.com/v2/pushes"
)

// Service providing Pushbullet as a notification service.
type Service struct {
	standard.Standard
	client jsonclient.Client
	Config *Config
	pkr    format.PropKeyResolver
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service.
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)

	service.Config = &Config{
		Title: "Shoutrrr notification", // Explicitly set default
	}
	service.pkr = format.NewPropKeyResolver(service.Config)

	if err := service.Config.setURL(&service.pkr, configURL); err != nil {
		return err
	}

	service.client = jsonclient.NewClient()
	service.client.Headers().Set("Access-Token", service.Config.Token)

	return nil
}

// GetID returns the service identifier.
func (service *Service) GetID() string {
	return Scheme
}

// Send a push notification via Pushbullet.
func (service *Service) Send(message string, params *types.Params) error {
	config := *service.Config
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

	// Validate response fields
	if response.Type != "note" {
		return fmt.Errorf("unexpected response type: got %s, expected note", response.Type)
	}

	if response.Body != message {
		return fmt.Errorf("response body mismatch: got %s, expected %s", response.Body, message)
	}

	if response.Title != config.Title {
		return fmt.Errorf("response title mismatch: got %s, expected %s", response.Title, config.Title)
	}

	if !response.Active {
		return fmt.Errorf("push notification is not active")
	}

	return nil
}
