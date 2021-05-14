package opsgenie

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
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

func (service *Service) sendAlert(url string, apiKey string, payload AlertPayload) error {
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	jsonBuffer := bytes.NewBuffer(jsonBody)

	req, err := http.NewRequest("POST", url, jsonBuffer)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "GenieKey "+apiKey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send notification to OpsGenie: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("OpsGenie notification returned %d HTTP status code. Cannot read body: %s", resp.StatusCode, err)
		}
		return fmt.Errorf("OpsGenie notification returned %d HTTP status code: %s", resp.StatusCode, body)
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
