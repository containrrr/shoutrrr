package telegram

import (
	"bytes"
	json2 "encoding/json"
	"errors"
	"fmt"
	. "github.com/containrrr/shoutrrr/pkg/plugins"
	"net/http"
	"regexp"
)

const (
	url = "https://api.telegram.org/bot"
	maxlength = 4096
)

type TelegramConfig struct {
	ApiToken string
	Channels []string
}

type TelegramJson struct {
	Text string `json:"text"`
	Id string `json:"chat_id"`
}

type TelegramPlugin struct {}

func (plugin *TelegramPlugin) Send(url string, message string) error {
	if len(message) > maxlength {
		return errors.New("message exceeds the max length")
	}
	config, err := plugin.CreateConfigFromUrl(url)
	if err != nil {
		return err
	}

	return sendMessageForChatIds(config, message)
}

func sendMessageForChatIds(config *TelegramConfig, message string) error {
	for _, channel := range config.Channels {
		if err := sendMessageToApi(message, channel, config.ApiToken); err != nil {
			return err
		}
	}
	return nil
}

func sendMessageToApi(message string, channel string, apiToken string) error {
	postUrl := fmt.Sprintf("%s%s/sendMessage", url, apiToken)
	json, _ := json2.Marshal(
		TelegramJson {
			Text: message,
			Id: channel,
		})

	res, err := http.Post(postUrl, "application/json", bytes.NewBuffer(json))
	if res.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("failed to send notification to \"%s\", response status code %s", channel, res.Status))
	}
	return err
}


func (plugin *TelegramPlugin) CreateConfigFromUrl(url string) (*TelegramConfig, error) {
	arguments, err := ExtractArguments(url)
	if err != nil {
		return nil, err
	}
	if len(arguments) < 2 {
		return nil, errors.New("the telegram plugin expects at least two arguments")
	}
	if !IsTokenValid(arguments[0]) {
		return nil, errors.New("invalid telegram token")
	}
	return &TelegramConfig{
		ApiToken: arguments[0],
		Channels: arguments[1:],
	}, nil
}

func IsTokenValid(token string) bool {
	matched, err := regexp.MatchString("^[0-9]+:[a-zA-Z0-9_-]+$", token)
	return matched && err == nil
}