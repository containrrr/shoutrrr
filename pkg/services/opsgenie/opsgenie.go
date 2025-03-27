package opsgenie

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/nicholas-fedor/shoutrrr/pkg/format"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/standard"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
)

const (
	alertEndpointTemplate = "https://%s:%d/v2/alerts"
	MaxMessageLength      = 130 // Maximum length of the alert message field in OpsGenie
	httpSuccessMax        = 299 // Maximum HTTP status code for a successful response
)

// Service providing OpsGenie as a notification service.
type Service struct {
	standard.Standard
	Config *Config
	pkr    format.PropKeyResolver
}

func (service *Service) sendAlert(url string, apiKey string, payload AlertPayload) error {
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	jsonBuffer := bytes.NewBuffer(jsonBody)

	req, err := http.NewRequest(http.MethodPost, url, jsonBuffer)
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

	if resp.StatusCode < http.StatusOK || resp.StatusCode > httpSuccessMax {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("OpsGenie notification returned %d HTTP status code. Cannot read body: %s", resp.StatusCode, err)
		}

		return fmt.Errorf("OpsGenie notification returned %d HTTP status code: %s", resp.StatusCode, body)
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

// Send a notification message to OpsGenie
// See: https://docs.opsgenie.com/docs/alert-api#create-alert
func (service *Service) Send(message string, params *types.Params) error {
	config := service.Config
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
	payloadFields := *service.Config

	if err := service.pkr.UpdateConfigFromParams(&payloadFields, params); err != nil {
		return AlertPayload{}, err
	}

	// Use `Message` for the title if available, or if the message is too long
	// Use `Description` for the message in these scenarios
	title := payloadFields.Title
	description := message

	if title == "" {
		if len(message) > MaxMessageLength {
			title = message[:MaxMessageLength]
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
