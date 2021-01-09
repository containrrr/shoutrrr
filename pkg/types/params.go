package types

// Params is the string map used to provide additional variables to the service templates
type Params map[string]string

const titleKey = "title"

func (p Params) SetTitle(title string) {
	p[titleKey] = title
}

func (p Params) Title() (title string, found bool) {
	title, found = p[titleKey]
	return
}