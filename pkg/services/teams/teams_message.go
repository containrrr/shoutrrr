package teams

// payload is the main structure for a Teams message card.
type payload struct {
	CardType   string    `json:"@type"`
	Context    string    `json:"@context"`
	ThemeColor string    `json:"themeColor,omitempty"`
	Summary    string    `json:"summary"`
	Title      string    `json:"title,omitempty"`
	Markdown   bool      `json:"markdown"`
	Sections   []section `json:"sections"`
}

// section represents a section of a Teams message card.
type section struct {
	ActivityTitle    string    `json:"activityTitle,omitempty"`
	ActivitySubtitle string    `json:"activitySubtitle,omitempty"`
	ActivityImage    string    `json:"activityImage,omitempty"`
	Facts            []fact    `json:"facts,omitempty"`
	Text             string    `json:"text,omitempty"`
	Images           []image   `json:"images,omitempty"`
	Actions          []action  `json:"potentialAction,omitempty"`
	HeroImage        *heroCard `json:"heroImage,omitempty"`
}

// fact represents a key-value pair in a Teams message card.
type fact struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// image represents an image in a Teams message card.
type image struct {
	Image string `json:"image"`
	Title string `json:"title,omitempty"`
}

// action represents an action button in a Teams message card.
type action struct {
	Type    string      `json:"@type"`
	Name    string      `json:"name"`
	Targets []target    `json:"targets,omitempty"`
	Actions []subAction `json:"actions,omitempty"`
}

// target represents a target for an action in a Teams message card.
type target struct {
	OS  string `json:"os"`
	URI string `json:"uri"`
}

// subAction represents a sub-action in a Teams message card.
type subAction struct {
	Type string `json:"@type"`
	Name string `json:"name"`
	URI  string `json:"uri"`
}

// heroCard represents a hero image in a Teams message card.
type heroCard struct {
	Image string `json:"image"`
}
