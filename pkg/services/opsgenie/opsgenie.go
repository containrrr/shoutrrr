package opsgenie

import (
	"fmt"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util/jsonclient"
)

const (
	alertEndpointTemplate = "https://%s:%d/v2/alerts"
)

// Service providing OpsGenie as a notification service
type Service struct {
	standard.Standard
	config *Config
	pkr    format.PropKeyResolver
}

// EmptyConfig returns an empty types.ServiceConfig for the service
func (service *Service) EmptyConfig() types.ServiceConfig {
	return &Config{}
}

func (service *Service) sendAlert(url string, apiKey string, payload AlertPayload) error {
	client := jsonclient.NewClient()
	client.Headers().Add("Authorization", "GenieKey "+apiKey)
	response := &map[string]interface{}{}

	if err := client.Post(url, payload, response); err != nil {
		service.Logf("Got response: %v", response)
		return fmt.Errorf("failed to send notification to OpsGenie: %s", err)
	}

	return nil
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{}
	service.pkr = format.NewPropKeyResolver(service.config)
	return service.config.setURL(&service.pkr, configURL)
}

// Send a notification message to OpsGenie
// See: https://docs.opsgenie.com/docs/alert-api#create-alert
func (service *Service) Send(message string, params *types.Params) error {
	config := service.config
	endpointURL := fmt.Sprintf(alertEndpointTemplate, config.Host, config.Port)
	payload, err := service.newAlertPayload(message, params)
	if err != nil {
		return err
	}
	return service.sendAlert(endpointURL, config.APIKey, payload)
}

func (service *Service) newAlertPayload(message string, params *types.Params) (AlertPayload, error) {
	if params == nil {
		params = &types.Params{}
	}

	// Defensive copy
	payloadFields := *service.config

	if err := service.pkr.UpdateConfigFromParams(&payloadFields, params); err != nil {
		return AlertPayload{}, err
	}

	// Use `Message` for the title if available, or if the message is too long
	// Use `Description` for the message in these scenarios
	title := payloadFields.Title
	description := message
	if title == "" {
		if len(message) > 130 {
			title = message[:130]
		} else {
			title = message
			description = ""
		}
	}

	if payloadFields.Description != "" && description != "" {
		description = description + "\n"
	}

	result := AlertPayload{
		Message:     title,
		Alias:       payloadFields.Alias,
		Description: description + payloadFields.Description,
		Responders:  payloadFields.Responders,
		VisibleTo:   payloadFields.VisibleTo,
		Actions:     payloadFields.Actions,
		Tags:        payloadFields.Tags,
		Details:     payloadFields.Details,
		Entity:      payloadFields.Entity,
		Source:      payloadFields.Source,
		Priority:    payloadFields.Priority,
		User:        payloadFields.User,
		Note:        payloadFields.Note,
	}
	return result, nil
}
