package discord

import (
	"encoding/json"
	"errors"
)


// JSON is the actual notification payload
type JSON struct {
	Text string `json:"content"`
}

// CreateJSONToSend creates a JSON payload to be sent to the discord webhook API
func CreateJSONToSend(message string) ([]byte, error) {
	if message == "" {
		return nil, errors.New("message was empty")
	}
	if len(message) > maxlength {
		return nil, errors.New("the supplied message exceeds the max length for discord")
	}
	return json.Marshal(JSON{
		Text: message,
	})
}
