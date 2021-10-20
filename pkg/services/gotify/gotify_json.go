package gotify

// payload is the actual payload being sent to the Gotify API
type payload struct {
	Message  string `json:"message"`
	Title    string `json:"title"`
	Priority int    `json:"priority"`
}

type payloadResponse struct{}
