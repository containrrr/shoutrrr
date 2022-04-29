package slack

import (
	"fmt"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/common/webclient"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Service sends notifications to a pre-configured channel or user
type Service struct {
	standard.Standard
	webclient.ClientService
	config *Config
	pkr    format.PropKeyResolver
}

const (
	apiPostMessage = "https://slack.com/api/chat.postMessage"
)

// Send a notification message to Slack
func (service *Service) Send(message string, params *types.Params) error {
	config := service.config

	if err := service.pkr.UpdateConfigFromParams(config, params); err != nil {
		return err
	}

	payload := CreateJSONPayload(config, message)

	var err error
	if config.Token.IsAPIToken() {
		err = service.sendAPI(config, payload)
	} else {
		err = service.sendWebhook(config, payload)
	}

	if err != nil {
		return fmt.Errorf("failed to send slack notification: %v", err)
	}

	return nil
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{}
	service.pkr = format.NewPropKeyResolver(service.config)

	if err := service.config.setURL(&service.pkr, configURL); err != nil {
		return err
	}

	client := service.WebClient()
	if service.config.Token.IsAPIToken() {
		client.Headers().Set("Authorization", service.config.Token.Authorization())
	} else {
		client.SetParser(parseWebhookResponse)
	}

	return nil

}

func (service *Service) sendAPI(config *Config, payload interface{}) error {
	response := APIResponse{}

	if err := service.WebClient().Post(apiPostMessage, payload, &response); err != nil {
		return err
	}

	if !response.Ok {
		if response.Error != "" {
			return fmt.Errorf("api response: %v", response.Error)
		}
		return fmt.Errorf("unknown error")
	}

	if response.Warning != "" {
		service.Logger.Logf("Slack API warning: %q", response.Warning)
	}

	return nil
}

func (service *Service) sendWebhook(config *Config, payload interface{}) error {
	var response *string
	err := service.WebClient().Post(config.Token.WebhookURL(), payload, &response)

	if err != nil {
		return fmt.Errorf("failed to invoke webhook: %v", err)
	}

	switch *response {
	case "":
		// Treat status 200 as no error regardless of actual content
		fallthrough
	case "ok":
		return nil
	default:
		return fmt.Errorf("webhook response: %v", response)
	}

}
