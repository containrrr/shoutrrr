package hangouts

// JSON is the actual payload being sent to the Hangouts Chat API.
type JSON struct {
	Text string `json:"text"`
}
