package teams

type TeamsJson struct {
	CardType string `json:"@type"`
	Context  string `json:"@context"`
	Markdown bool   `json:"markdown,bool"`
	Text     string `json:"text,omitempty"`
}
