package types

import (
	"text/template"
)

// Templater is the interface for the service template API
type Templater interface {
	GetTemplate(id string) (template *template.Template, found bool)
	SetTemplateString(id string, body string) error
	SetTemplateFile(id string, file string) error
}
