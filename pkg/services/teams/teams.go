package teams

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
)

// Service providing teams as a notification service
type Service struct {
	standard.Standard
	config *Config
}

// Send a notification message to Microsoft Teams
func (service *Service) Send(message string, params *map[string]string) error {
	config := service.config

	postURL := buildURL(config)
	return service.doSend(postURL, message)
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

func (service *Service) doSend(postURL string, message string) error {
	body := JSON{
		CardType: "MessageCard",
		Context:  "http://schema.org/extensions",
		Markdown: true,
		Text:     message,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	res, err := http.Post(postURL, "application/json", bytes.NewBuffer(jsonBody))
	if res.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("failed to send notification to teams, response status code %s", res.Status)
		return errors.New(msg)
	}
	return nil
}

func buildURL(config *Config) string {
	var baseURL = "https://outlook.office.com/webhook"
	return fmt.Sprintf(
		"%s/%s/IncomingWebhook/%s/%s",
		baseURL,
		config.Token.A,
		config.Token.B,
		config.Token.C)
}
