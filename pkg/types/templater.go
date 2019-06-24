package types

import (
	"text/template"
)

type Templater interface {
// GetTemplate attempts to retrieve the template identified with id
GetTemplate (id string) (template *template.Template, found bool)

SetTemplateString (id string, body string) error

SetTemplateFile (id string, file string) error

}