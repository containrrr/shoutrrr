package types

// Params is the string map used to provide additional variables to the service templates
type Params map[string]string

const (
	// TitleKey is the common key for the title prop
	TitleKey = "title"
	// MessageKey is the common key for the message prop
	MessageKey = "message"
)

// SetTitle sets the "title" param to the specified value
func (p Params) SetTitle(title string) {
	p[TitleKey] = title
}

// Title returns the "title" param
func (p Params) Title() (title string, found bool) {
	title, found = p[TitleKey]
	return
}

// SetMessage sets the "message" param to the specified value
func (p Params) SetMessage(message string) {
	p[MessageKey] = message
}
