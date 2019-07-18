package teams

// JSON is the actual payload being sent to the teams api
type JSON struct {
	CardType string `json:"@type"`
	Context  string `json:"@context"`
	Markdown bool   `json:"markdown,bool"`
	Text     string `json:"text,omitempty"`
}
