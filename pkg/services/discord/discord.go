package discord

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Service providing Discord as a notification service
type Service struct {
	standard.Standard
	config *Config
}

const (
	hookURL   = "https://discordapp.com/api/webhooks"
	maxlength = 2000
)

// Send a notification message to discord
func (service *Service) Send(message string, params *types.Params) error {

	payload, err := CreateJSONToSend(message, service.config.JSON)
	if err != nil {
		return err
	}

	postURL := CreateAPIURLFromConfig(service.config)

	return doSend(payload, postURL)
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger *log.Logger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{}
	if err := service.config.SetURL(configURL); err != nil {
		return err
	}

	return nil
}

// CreateAPIURLFromConfig takes a discord config object and creates a post url
func CreateAPIURLFromConfig(config *Config) string {
	return fmt.Sprintf(
		"%s/%s/%s",
		hookURL,
		config.Channel,
		config.Token)
}

func doSend(payload []byte, postURL string) error {
	res, err := http.Post(postURL, "application/json", bytes.NewBuffer(payload))

	if res == nil && err == nil {
		err = fmt.Errorf("unknown error")
	}

	if err == nil && res.StatusCode != http.StatusNoContent {
		err = fmt.Errorf("response status code %s", res.Status)
	}

	if err != nil {
		return fmt.Errorf("failed to send discord notification: %v", err)
	}

	return nil
}
