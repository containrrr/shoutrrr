package zulip

import (
	"fmt"
	"log"
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
func (service *Service) Initialize(configURL *url.URL, logger *log.Logger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{}

	if err := service.config.SetURL(nil, configURL); err != nil {
		return err
	}

	return nil
}

func (service *Service) doSend(config *Config, message string) error {
	apiURL := service.getURL(config)
	payload := CreatePayload(config, message)
	res, err := http.Post(apiURL, "application/x-www-form-urlencoded", strings.NewReader(payload.Encode()))

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send notification to service, response status code %s", res.Status)
	}

	return err
}

func (service *Service) getURL(config *Config) string {
	url := url.URL{
		User:   url.UserPassword(config.BotMail, config.BotKey),
		Host:   config.Host,
		Path:   config.Path,
		Scheme: "https",
	}

	return url.String()
}
