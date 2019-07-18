package telegram

// JSON to be used as a notification payload for the telegram notification service
type JSON struct {
	Text string `json:"text"`
	ID   string `json:"chat_id"`
}
