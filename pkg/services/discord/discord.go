package discord

import (
	"bytes"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/types"
	"net/http"
	"net/url"
)

// Service providing Discord as a notification service
type Service struct {}

const (
	hookURL   = "https://discordapp.com/api/webhooks"
	maxlength = 2000
)

// Send a notification message to discord
func (plugin *Service) Send(rawURL *url.URL, message string, opts types.ServiceOpts) error {
	config, err := plugin.CreateConfigFromURL(rawURL)
	if err != nil {
		return err
	}

	payload, err := CreateJSONToSend(message)
	if err != nil {
		return err
	}
	fmt.Println(string(payload))

	postURL := CreateAPIURLFromConfig(config)
	fmt.Println(postURL)

	return doSend(payload, postURL)
}

func (plugin *Service) URLToStringMap(url url.URL) (map[string]string, error) {
	return map[string]string {}, nil
}

func (plugin *Service) GetConfig() types.ServiceConfig {
	return &Config{}
}

// CreateAPIURLFromConfig takes a discord config object and creates a post url
func CreateAPIURLFromConfig(config *Config) string {
	return fmt.Sprintf(
		"%s/%s/%s",
		hookURL,
		config.Channel,
		config.Token)
}

// CreateConfigFromURL creates a Config struct given a valid discord notification url
func (plugin *Service) CreateConfigFromURL(rawURL *url.URL) (*Config, error) {
	config := Config{}
	err := config.SetURL(rawURL)
	return &config, err
}

func doSend(payload []byte, postURL string) error {
	res, err := http.Post(postURL, "application/json", bytes.NewBuffer(payload))
	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to send notification to discord, response status code %s", res.Status)
	}
	return err
}