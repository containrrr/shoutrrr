package teams

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/types"
	"net/http"
	"net/url"
)

// Service providing teams as a notification service
type Service struct{}

// Send a notification message to Microsoft Teams
func (plugin *Service) Send(url *url.URL, message string, opts types.ServiceOpts) error {
	config, err := plugin.CreateConfigFromURL(url)
	if err != nil {
		return err
	}

	postURL := buildURL(config)
	return plugin.doSend(postURL, message)
}

// GetConfig returns an empty ServiceConfig for this Service
func (plugin *Service) GetConfig() types.ServiceConfig {
	return &Config{}
}

func (plugin *Service) doSend(postURL string, message string) error {
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