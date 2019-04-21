package slack

import "encoding/json"

type SlackJson struct {
	Text string `json:"text"`
	Botname string `json:"username"`
}

func CreateJsonPayload(config *SlackConfig, message string) ([]byte, error) {
	return json.Marshal(
		SlackJson {
			Text: message,
			Botname: config.Botname,
		})
}