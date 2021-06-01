package googlechat

// JSON is the actual payload being sent to the Google Chat API.
type JSON struct {
	Text string `json:"text"`
}
