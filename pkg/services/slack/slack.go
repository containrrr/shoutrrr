package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/nicholas-fedor/shoutrrr/pkg/format"
	"github.com/nicholas-fedor/shoutrrr/pkg/util/jsonclient"

	"github.com/nicholas-fedor/shoutrrr/pkg/services/standard"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
)

// Service sends notifications to a pre-configured channel or user.
type Service struct {
	standard.Standard
	Config *Config
	pkr    format.PropKeyResolver
}

const (
	apiPostMessage = "https://slack.com/api/chat.postMessage"
)

// Send a notification message to Slack.
func (service *Service) Send(message string, params *types.Params) error {
	config := service.Config

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
		return fmt.Errorf("failed to send slack notification: %w", err)
	}

	return nil
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service.
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)
	service.Config = &Config{}
	service.pkr = format.NewPropKeyResolver(service.Config)

	return service.Config.setURL(&service.pkr, configURL)
}

// GetID returns the service identifier.
func (service *Service) GetID() string {
	return Scheme
}

func (service *Service) sendAPI(config *Config, payload interface{}) error {
	response := APIResponse{}
	jsonClient := jsonclient.NewClient()
	jsonClient.Headers().Set("Authorization", config.Token.Authorization())

	if err := jsonClient.Post(apiPostMessage, payload, &response); err != nil {
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
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	res, err := http.Post(config.Token.WebhookURL(), jsonclient.ContentType, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to invoke webhook: %w", err)
	}

	defer res.Body.Close()
	resBytes, _ := io.ReadAll(res.Body)
	response := string(resBytes)

	switch response {
	case "":
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("webhook status: %v", res.Status)
		}
		// Treat status 200 as no error regardless of actual content
		fallthrough
	case "ok":
		return nil
	default:
		return fmt.Errorf("webhook response: %v", response)
	}
}
