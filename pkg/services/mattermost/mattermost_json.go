package mattermost

import (
	"encoding/json"

	"github.com/containrrr/shoutrrr/pkg/types"
)

// JSON payload for mattermost notifications
type JSON struct {
	Text     string `json:"text"`
	UserName string `json:"username,omitempty"`
	Channel  string `json:"channel,omitempty"`
}

// CreateJSONPayload for usage with the mattermost service
func CreateJSONPayload(config *Config, message string, params *types.Params) ([]byte, error) {
	payload := JSON{
		Text:     message,
		UserName: config.UserName,
		Channel:  config.Channel,
	}

	if params != nil {
		if value, found := (*params)["username"]; found {
			payload.UserName = value
		}
		if value, found := (*params)["channel"]; found {
			payload.Channel = value
		}
	}
	return json.Marshal(payload)
}
