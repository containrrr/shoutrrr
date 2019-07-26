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

	payload, err := CreateJSONToSend(message)
	if err != nil {
		return err
	}
	fmt.Println(string(payload))

	postURL := CreateAPIURLFromConfig(service.config)
	fmt.Println(postURL)

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
	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to send notification to discord, response status code %s", res.Status)
	}
	return err
}
