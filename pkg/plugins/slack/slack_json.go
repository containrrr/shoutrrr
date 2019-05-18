package slack

import "encoding/json"

type SlackJSON struct {
	Text string `json:"text"`
	Botname string `json:"username"`
}

// CreateJSONPayload compatible with the slack webhook api
func CreateJSONPayload(config *SlackConfig, message string) ([]byte, error) {
	return json.Marshal(
		SlackJSON{
			Text: message,
			Botname: config.Botname,
		})
}