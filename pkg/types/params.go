package types

// Params is the string map used to provide additional variables to the service templates
type Params map[string]string

const titleKey = "title"

// SetTitle sets the "title" param to the specified value
func (p Params) SetTitle(title string) {
	p[titleKey] = title
}

// Title returns the "title" param
func (p Params) Title() (title string, found bool) {
	title, found = p[titleKey]
	return
}
