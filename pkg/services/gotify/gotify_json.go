package gotify

// JSON is the actual payload being sent to the Gotify API
type JSON struct {
	Message  string `json:"message"`
	Title    string `json:"title"`
	Priority int    `json:"priority"`
}
