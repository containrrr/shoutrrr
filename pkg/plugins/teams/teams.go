package teams

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type TeamsPlugin struct{}

// Send a notification message to Microsoft Teams
func (plugin *TeamsPlugin) Send(url string, message string) error {
	config, err := plugin.CreateConfigFromURL(url)
	if err != nil {
		return err
	}

	postUrl := buildURL(config)
	return plugin.doSend(postUrl, message)
}

func (plugin *TeamsPlugin) doSend(postUrl string, message string) error {
	body := TeamsJson{
		CardType: "MessageCard",
		Context:  "http://schema.org/extensions",
		Markdown: true,
		Text:     message,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	res, err := http.Post(postUrl, "application/json", bytes.NewBuffer(jsonBody))
	if res.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("failed to send notification to teams, response status code %s", res.Status)
		return errors.New(msg)
	}
	return nil
}

func buildURL(config *TeamsConfig) string {
	var baseUrl = "https://outlook.office.com/webhook"
	return fmt.Sprintf(
		"%s/%s/IncomingWebhook/%s/%s",
		baseUrl,
		config.Token.A,
		config.Token.B,
		config.Token.C)
}