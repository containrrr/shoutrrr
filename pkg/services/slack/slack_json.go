package slack

import (
	"regexp"
	"strings"
)

// messagePayload used within the Slack service
type messagePayload struct {
	Text        string       `json:"text"`
	BotName     string       `json:"username,omitempty"`
	Blocks      []block      `json:"blocks,omitempty"`
	Attachments []attachment `json:"attachments,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	IconURL     string       `json:"icon_url,omitempty"`
}

type block struct {
	Type string    `json:"type"`
	Text blockText `json:"text"`
}

type blockText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type attachment struct {
	Title    string        `json:"title,omitempty"`
	Fallback string        `json:"fallback,omitempty"`
	Text     string        `json:"text"`
	Color    string        `json:"color,omitempty"`
	Fields   []legacyField `json:"fields,omitempty"`
	Footer   string        `json:"footer,omitempty"`
	Time     int           `json:"ts,omitempty"`
}

type legacyField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short,omitempty"`
}

type APIResponse struct {
	Ok       bool   `json:"ok"`
	Error    string `json:"error"`
	Warning  string `json:"warning"`
	MetaData struct {
		Warnings []string `json:"warnings"`
	} `json:"response_metadata"`
}

var iconUrlPattern = regexp.MustCompile(`https?://`)

// CreateJSONPayload compatible with the slack post message API
func CreateJSONPayload(config *Config, message string) interface{} {

	var atts []attachment
	for _, line := range strings.Split(message, "\n") {
		atts = append(atts, attachment{
			Text:  line,
			Color: config.Color,
		})
	}

	payload := messagePayload{
		Text:        config.Title,
		BotName:     config.BotName,
		Attachments: atts,
	}

	if config.Icon != "" {
		if iconUrlPattern.MatchString(config.Icon) {
			payload.IconURL = config.Icon
		} else {
			payload.IconEmoji = config.Icon
		}
	}

	if config.Channel != "webhook" {
		payload.Channel = config.Channel
	}

	return payload
}
