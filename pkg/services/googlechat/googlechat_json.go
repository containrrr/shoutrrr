package googlechat

// payload is the actual payload being sent to the Google Chat API.
type payload struct {
	Text string `json:"text"`
}
