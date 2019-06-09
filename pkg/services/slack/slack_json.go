package slack

import "encoding/json"

// JSON used within the Slack service
type JSON struct {
	Text    string `json:"text"`
	BotName string `json:"username"`
}

// CreateJSONPayload compatible with the slack webhook api
func CreateJSONPayload(config *Config, message string) ([]byte, error) {
	return json.Marshal(
		JSON{
			Text:    message,
			BotName: config.BotName,
		})
}