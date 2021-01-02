package telegram

// JSON to be used as a notification payload for the telegram notification service
type SendMessagePayload struct {
	Text                string `json:"text"`
	ID                  string `json:"chat_id"`
	ParseMode           string `json:"parse_mode,omitempty"`
	DisablePreview      bool   `json:"disable_web_page_preview"`
	DisableNotification bool   `json:"disable_notification"`
}

func createSendMessagePayload(message string, channel string, config *Config) SendMessagePayload {
	payload := SendMessagePayload{
		Text:                message,
		ID:                  channel,
		DisableNotification: !config.Notification,
		DisablePreview:      !config.Preview,
	}

	if config.ParseMode != parseModes.None {
		payload.ParseMode = config.ParseMode.String()
	}

	return payload
}
