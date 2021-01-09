package discord

import (
	"encoding/json"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/types/rich"
	"github.com/containrrr/shoutrrr/pkg/util"
	"time"
)

// JSON is the actual notification payload
type WebhookPayload struct {
	Embeds []EmbedItem `json:"content"`
}

// JSON is the actual notification payload
type EmbedItem struct {
	Title     string       `json:"title,omitempty"`
	Content   string       `json:"description,omitempty"`
	URL       string       `json:"url,omitempty"`
	Timestamp string       `json:"timestamp,omitempty"`
	Color     int          `json:"color,omitempty"`
	Footer    *EmbedFooter `json:"footer,omitempty"`
}

type EmbedFooter struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url,omitempty"`
}

// CreateJSONToSend creates a JSON payload to be sent to the discord webhook API
func CreatePayloadFromItems(items []rich.MessageItem, title string, colors [rich.MessageLevelCount]int, omitted int) ([]byte, error) {

	itemCount := util.Min(9, len(items))
	embeds := make([]EmbedItem, 1, itemCount+1)

	for _, item := range items {

		color := 0
		if item.Level >= rich.Unknown && int(item.Level) < len(colors) {
			color = colors[item.Level]
		}

		ei := EmbedItem{
			Content: item.Text,
			Color:   color,
		}

		if item.Level != rich.Unknown {
			ei.Footer = &EmbedFooter{
				Text: item.Level.String(),
			}
		}

		if item.Timestamp != nil {
			ei.Timestamp = item.Timestamp.UTC().Format(time.RFC3339)
		}

		embeds = append(embeds, ei)
	}

	embeds[0].Title = title
	if omitted > 0 {
		embeds[0].Footer = &EmbedFooter{
			Text: fmt.Sprintf("... (%v character(s) where omitted)", omitted),
		}
	}

	return json.Marshal(WebhookPayload{
		Embeds: embeds,
	})
}
