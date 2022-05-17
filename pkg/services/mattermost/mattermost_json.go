package mattermost

import (
	"encoding/json"
	"regexp"

	"github.com/containrrr/shoutrrr/pkg/types"
)

// JSON payload for mattermost notifications
type JSON struct {
	Text      string `json:"text"`
	UserName  string `json:"username,omitempty"`
	Channel   string `json:"channel,omitempty"`
	IconEmoji string `json:"icon_emoji,omitempty"`
	IconURL   string `json:"icon_url,omitempty"`
}

var iconURLPattern = regexp.MustCompile(`https?://`)

// SetIcon sets the appropriate icon field in the payload based on whether the input is a URL or not
func (j *JSON) SetIcon(icon string) {
	j.IconURL = ""
	j.IconEmoji = ""

	if icon != "" {
		if iconURLPattern.MatchString(icon) {
			j.IconURL = icon
		} else {
			j.IconEmoji = icon
		}
	}
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
	payload.SetIcon(config.Icon)

	return json.Marshal(payload)
}
