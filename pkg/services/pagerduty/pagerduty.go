package pagerduty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
	"io"
	"net/http"
	"net/url"
)

const (
	eventEndpointTemplate = "https://%s:%d/v2/enqueue"
)

// Service providing PagerDuty as a notification service
type Service struct {
	standard.Standard
	config *Config
	pkr    format.PropKeyResolver
}

func (service *Service) sendAlert(url string, payload EventPayload) error {
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	jsonBuffer := bytes.NewBuffer(jsonBody)

	req, err := http.NewRequest("POST", url, jsonBuffer)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send notification to PagerDuty: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("PagerDuty notification returned %d HTTP status code. Cannot read body: %s", resp.StatusCode, err)
		}
		return fmt.Errorf("PagerDuty notification returned %d HTTP status code: %s", resp.StatusCode, body)
	}

	return nil
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{}
	service.pkr = format.NewPropKeyResolver(service.config)

	if err := service.setDefaults(); err != nil {
		return err
	}

	return service.config.setURL(&service.pkr, configURL)
}

// Send a notification message to PagerDuty
// See: https://developer.pagerduty.com/docs/events-api-v2-overview
func (service *Service) Send(message string, params *types.Params) error {
	config := service.config
	endpointURL := fmt.Sprintf(eventEndpointTemplate, config.Host, config.Port)

	payload, err := service.newEventPayload(message, params)
	if err != nil {
		return err
	}

	return service.sendAlert(endpointURL, payload)
}

func (service *Service) newEventPayload(message string, params *types.Params) (EventPayload, error) {
	if params == nil {
		params = &types.Params{}
	}

	// Defensive copy
	payloadFields := *service.config

	if err := service.pkr.UpdateConfigFromParams(&payloadFields, params); err != nil {
		return EventPayload{}, err
	}

	// The maximum permitted length of this property is 1024 characters.
	if len(message) > 1024 {
		message = message[:1024]
	}

	result := EventPayload{
		Payload: Payload{
			Summary:  message,
			Severity: payloadFields.Severity,
			Source:   payloadFields.Source,
		},
		RoutingKey:  payloadFields.IntegrationKey,
		EventAction: payloadFields.Action,
	}
	return result, nil
}

func (service *Service) setDefaults() error {
	if err := service.pkr.SetDefaultProps(service.config); err != nil {
		return err
	}

	setUrlDefaults(service.config)
	return nil
}
