package zulip

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Service sends notifications to a pre-configured channel or user
type Service struct {
	standard.Standard
	config *Config
}

const (
	contentMaxSize = 10000 // bytes
	topicMaxLength = 60    // characters
)

// Send a notification message to Zulip
func (service *Service) Send(message string, params *types.Params) error {
	// Clone the config because we might modify stream and/or
	// topic with values from the parameters and they should only
	// change this Send().
	config := service.config.Clone()

	if params != nil {
		if stream, found := (*params)["stream"]; found {
			config.Stream = stream
		}

		if topic, found := (*params)["topic"]; found {
			config.Topic = topic
		}
	}

	topicLength := len([]rune(config.Topic))

	if topicLength > topicMaxLength {
		return fmt.Errorf(string(TopicTooLong), topicMaxLength, topicLength)
	}

	messageSize := len(message)

	if messageSize > contentMaxSize {
		return fmt.Errorf("message exceeds max size (%d bytes): was %d bytes", contentMaxSize, messageSize)
	}

	return service.doSend(config, message)
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{}

	if err := service.config.setURL(nil, configURL); err != nil {
		return err
	}

	return nil
}

func (service *Service) doSend(config *Config, message string) error {
	apiURL := service.getAPIURL(config)
	payload := CreatePayload(config, message)
	res, err := http.Post(apiURL, "application/x-www-form-urlencoded", strings.NewReader(payload.Encode()))

	if err == nil && res.StatusCode != http.StatusOK {
		err = fmt.Errorf("response status code %s", res.Status)
	}

	if err != nil {
		return fmt.Errorf("failed to send zulip message: %s", err)
	}

	return nil
}

func (service *Service) getAPIURL(config *Config) string {
	return (&url.URL{
		User:   url.UserPassword(config.BotMail, config.BotKey),
		Host:   config.Host,
		Path:   "api/v1/messages",
		Scheme: "https",
	}).String()
}
