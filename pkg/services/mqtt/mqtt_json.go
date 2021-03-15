package mqtt

// JSON to be used as a notification payload for the telegram notification service
type SendMessagePayload struct {
	Text                string `json:"text"`
	Topic               string `json:"topic"`
	ParseMode           string `json:"parse_mode,omitempty"`
}

func createSendMessagePayload(message string, topic string, config *Config) SendMessagePayload {
	payload := SendMessagePayload{
		Text:                message,
		Topic:               topic,
	}

	if config.ParseMode != parseModes.None {
		payload.ParseMode = config.ParseMode.String()
	}

	return payload
}
