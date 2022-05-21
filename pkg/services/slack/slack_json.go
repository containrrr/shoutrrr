package slack

import (
	"regexp"
	"strings"
)

// MessagePayload used within the Slack service
type MessagePayload struct {
	Text        string       `json:"text"`
	BotName     string       `json:"username,omitempty"`
	Blocks      []block      `json:"blocks,omitempty"`
	Attachments []attachment `json:"attachments,omitempty"`
	ThreadTS    string       `json:"thread_ts,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	IconURL     string       `json:"icon_url,omitempty"`
}

var iconURLPattern = regexp.MustCompile(`https?://`)

// SetIcon sets the appropriate icon field in the payload based on whether the input is a URL or not
func (p *MessagePayload) SetIcon(icon string) {
	p.IconURL = ""
	p.IconEmoji = ""

	if icon != "" {
		if iconURLPattern.MatchString(icon) {
			p.IconURL = icon
		} else {
			p.IconEmoji = icon
		}
	}
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

// APIResponse is the default generic response message sent from the API
type APIResponse struct {
	Ok       bool   `json:"ok"`
	Error    string `json:"error"`
	Warning  string `json:"warning"`
	MetaData struct {
		Warnings []string `json:"warnings"`
	} `json:"response_metadata"`
}

// CreateJSONPayload compatible with the slack post message API
func CreateJSONPayload(config *Config, message string) interface{} {

	var atts []attachment
	for i, line := range strings.Split(message, "\n") {
		// When 100 attachments have been reached, append the remaining line to the last
		// attachment to prevent reaching the slack API limit
		if i >= 100 {
			atts[len(atts)-1].Text += "\n" + line
			continue
		}
		atts = append(atts, attachment{
			Text:  line,
			Color: config.Color,
		})
	}

	// Remove last attachment if empty
	if atts[len(atts)-1].Text == "" {
		atts = atts[:len(atts)-1]
	}

	payload := MessagePayload{
		ThreadTS:    config.ThreadTS,
		Text:        config.Title,
		BotName:     config.BotName,
		Attachments: atts,
	}

	payload.SetIcon(config.Icon)

	if config.Channel != "webhook" {
		payload.Channel = config.Channel
	}

	return payload
}
