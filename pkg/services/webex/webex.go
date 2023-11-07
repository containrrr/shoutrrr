package webex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Service providing Webex as a notification service
type Service struct {
	standard.Standard
	config *Config
	pkr    format.PropKeyResolver
}

const (
	MessagesEndpoint = "https://webexapis.com/v1/messages"
)

// MessagePayload is the message endpoint payload
type MessagePayload struct {
	RoomID   string `json:"roomId"`
	Markdown string `json:"markdown,omitempty"`
}

// Send a notification message to webex
func (service *Service) Send(message string, params *types.Params) error {
	err := doSend(message, service.config)
	if err != nil {
		return fmt.Errorf("failed to send webex notification: %v", err)
	}

	return nil
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{}
	service.pkr = format.NewPropKeyResolver(service.config)

	if err := service.pkr.SetDefaultProps(service.config); err != nil {
		return err
	}

	if err := service.config.SetURL(configURL); err != nil {
		return err
	}

	return nil
}

func doSend(message string, config *Config) error {
	var firstErr error

	for _, room := range config.Rooms {
		req, firstErr := BuildRequestFromPayloadAndConfig(message, room, config)
		if firstErr != nil {
			return firstErr
		}

		res, firstErr := http.DefaultClient.Do(req)

		if res == nil && firstErr == nil {
			firstErr = fmt.Errorf("unknown error")
		}

		if firstErr == nil && res.StatusCode != http.StatusOK {
			firstErr = fmt.Errorf("response status code %s", res.Status)
		}

		if firstErr != nil {
			return firstErr
		}
	}

	return firstErr
}

func BuildRequestFromPayloadAndConfig(message string, room string, config *Config) (*http.Request, error) {
	var err error
	payload := MessagePayload{
		RoomID:   room,
		Markdown: message,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", MessagesEndpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+config.BotToken)
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}
