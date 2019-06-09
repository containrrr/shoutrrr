package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/plugin"
	"net/http"
	"net/url"
)

const (
	apiBase = "https://api.telegram.org/bot"
	maxlength = 4096
)


// Plugin sends notifications to a given telegram chat
type Plugin struct {}

// Send notification to Telegram
func (plugin *Plugin) Send(url url.URL, message string, opts plugin.PluginOpts) error {
	if len(message) > maxlength {
		return errors.New("message exceeds the max length")
	}
	config, err := plugin.CreateConfigFromURL(url)
	if err != nil {
		return err
	}

	return sendMessageForChatIDs(config, message)
}

func (plugin *Plugin) GetConfig() plugin.PluginConfig {
	return Config{}
}

func sendMessageForChatIDs(config *Config, message string) error {
	for _, channel := range config.Channels {
		if err := sendMessageToAPI(message, channel, config.Token); err != nil {
			return err
		}
	}
	return nil
}

func sendMessageToAPI(message string, channel string, apiToken string) error {
	postURL := fmt.Sprintf("%s%s/sendMessage", apiBase, apiToken)
	jsonData, _ := json.Marshal(
		JSON{
			Text: message,
			ID:   channel,
		})

	res, err := http.Post(postURL, "application/jsonData", bytes.NewBuffer(jsonData))
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send notification to \"%s\", response status code %s", channel, res.Status)
	}
	return err
}


// CreateConfigFromURL to use within the telegram plugin
func (plugin *Plugin) CreateConfigFromURL(url url.URL) (*Config, error) {
	config := Config{}
	err := config.SetURL(url)
	return &config, err
}
