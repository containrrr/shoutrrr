package teams

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Plugin providing teams as a notification service
type Plugin struct{}

// Send a notification message to Microsoft Teams
func (plugin *Plugin) Send(url string, message string) error {
	config, err := plugin.CreateConfigFromURL(url)
	if err != nil {
		return err
	}

	postURL := buildURL(config)
	return plugin.doSend(postURL, message)
}

func (plugin *Plugin) doSend(postURL string, message string) error {
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