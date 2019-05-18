package discord

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/plugins"
	"net/http"
)

type DiscordPlugin struct {}

const (
	hookUrl = "https://discordapp.com/api/webhooks"
	maxlength = 2000
)

// Send a notification message to discord
func (plugin *DiscordPlugin) Send(url string, message string) error {
	config, err := plugin.CreateConfigFromURL(url)
	if err != nil {
		return err
	}

	payload, err := CreateJsonToSend(message)
	if err != nil {
		return err
	}
	fmt.Println(string(payload))

	apiUrl := CreateApiURLFromConfig(config)
	fmt.Println(apiUrl)

	return doSend(payload, apiUrl)
}

// CreateApiURLFromConfig takes a discord config object and creates a post url
func CreateApiURLFromConfig(config DiscordConfig) string {
	return fmt.Sprintf(
		"%s/%s/%s",
		hookUrl,
		config.Channel,
		config.Token)
}

// CreateConfigFromURL creates a DiscordConfig struct given a valid discord notification url
func (plugin *DiscordPlugin) CreateConfigFromURL(url string) (DiscordConfig, error) {
	args, err := plugins.ExtractArguments(url)
	if err != nil {
		return DiscordConfig{}, err
	}
	if len(args) != 2 {
		return DiscordConfig{}, errors.New("the discord plugin expects exactly two url path arguments")
	}

	return DiscordConfig{
		Channel: args[0],
		Token: args[1],
	}, nil
}

func doSend(payload []byte, postUrl string) error {
	res, err := http.Post(postUrl, "application/json", bytes.NewBuffer(payload))
	if res.StatusCode != http.StatusNoContent {
		return errors.New(fmt.Sprintf("failed to send notification to discord, response status code %s", res.Status))
	}
	return err
}