package discord

import (
	"bytes"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
	"net/http"
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
func (plugin *Service) Send(message string, params *map[string]string) error {

	payload, err := CreateJSONToSend(message)
	if err != nil {
		return err
	}
	fmt.Println(string(payload))

	postURL := CreateAPIURLFromConfig(plugin.config)
	fmt.Println(postURL)

	return doSend(payload, postURL)
}

// NewConfig returns an empty ServiceConfig for this Service
func (plugin *Service) NewConfig() types.ServiceConfig {
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



func doSend(payload []byte, postURL string) error {
	res, err := http.Post(postURL, "application/json", bytes.NewBuffer(payload))
	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to send notification to discord, response status code %s", res.Status)
	}
	return err
}