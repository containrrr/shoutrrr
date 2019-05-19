package discord

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/plugins"
	"net/http"
)

// Plugin providing Discord as a notification service
type Plugin struct {}

const (
	hookURL   = "https://discordapp.com/api/webhooks"
	maxlength = 2000
)

// Send a notification message to discord
func (plugin *Plugin) Send(url string, message string) error {
	config, err := plugin.CreateConfigFromURL(url)
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

// CreateAPIURLFromConfig takes a discord config object and creates a post url
func CreateAPIURLFromConfig(config Config) string {
	return fmt.Sprintf(
		"%s/%s/%s",
		hookURL,
		config.Channel,
		config.Token)
}

// CreateConfigFromURL creates a Config struct given a valid discord notification url
func (plugin *Plugin) CreateConfigFromURL(url string) (Config, error) {
	args, err := plugins.ExtractArguments(url)
	if err != nil {
		return Config{}, err
	}
	if len(args) != 2 {
		return Config{}, errors.New("the discord plugin expects exactly two url path arguments")
	}

	return Config{
		Channel: args[0],
		Token: args[1],
	}, nil
}

func doSend(payload []byte, postURL string) error {
	res, err := http.Post(postURL, "application/json", bytes.NewBuffer(payload))
	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to send notification to discord, response status code %s", res.Status)
	}
	return err
}