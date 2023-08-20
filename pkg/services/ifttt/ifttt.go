package ifttt

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/format"

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

	payload, err := createJSONToSend(config, message, params)
	fmt.Printf("%+v", payload)
	if err != nil {
		return err
	}
	for _, event := range config.Events {
		apiURL := service.createAPIURLForEvent(event)
		err := doSend(payload, apiURL)
		if err != nil {
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

func doSend(payload []byte, postURL string) error {
	res, err := http.Post(postURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	if res.StatusCode > 299 || res.StatusCode < 200 {
		return fmt.Errorf("got response status code %s", res.Status)
	}
	return nil
}
