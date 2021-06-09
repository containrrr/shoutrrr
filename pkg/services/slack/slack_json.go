package slack

import (
	"encoding/json"
	"strings"
)

// JSON used within the Slack service
type JSON struct {
	Text        string       `json:"text"`
	BotName     string       `json:"username,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	Emoji       string       `json:"icon_emoji,omitempty"`
	Blocks      []block      `json:"blocks,omitempty"`
	Attachments []attachment `json:"attachments,omitempty"`
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

// CreateJSONPayload compatible with the slack webhook api
func CreateJSONPayload(config *Config, message string) ([]byte, error) {

	var atts []attachment
	for _, line := range strings.Split(message, "\n") {
		atts = append(atts, attachment{
			Text:  line,
			Color: config.Color,
		})
	}

	return json.Marshal(
		JSON{
			Text:        config.Title,
			BotName:     config.BotName,
			Channel:     config.Channel,
			Emoji:       config.Emoji,
			Attachments: atts,
		})
}
