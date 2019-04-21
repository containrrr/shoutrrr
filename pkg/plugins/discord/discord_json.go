package discord

import (
	"encoding/json"
	"errors"
)


type DiscordJson struct {
	Text string `json:"content"`
}

func CreateJsonToSend(message string) ([]byte, error) {
	if message == "" {
		return nil, errors.New("message was empty")
	}
	if len(message) > maxlength {
		return nil, errors.New("the supplied message exceeds the max length for discord")
	}
	return json.Marshal(DiscordJson {
		Text: message,
	})
}
