package teams

// JSON is the actual payload being sent to the teams api
type Payload struct {
	CardType   string    `json:"@type"`
	Context    string    `json:"@context"`
	Markdown   bool      `json:"markdown,bool"`
	Text       string    `json:"text,omitempty"`
	Title      string    `json:"title,omitempty"`
	Summary    string    `json:"summary,omitempty"`
	Sections   []Section `json:"sections,omitempty"`
	ThemeColor string    `json:"themeColor,omitempty"`
}

type Section struct {
	Text         string `json:"text,omitempty"`
	ActivityText string `json:"activityText,omitempty"`
	StartGroup   bool   `json:"startGroup"`
	Facts        []Fact `json:"facts,omitempty"`
}

type Fact struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
