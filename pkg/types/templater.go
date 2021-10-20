package types

import (
	"io"
)

type StdTemplate interface {
	Execute(wr io.Writer, data interface{}) error
}

// Templater is the interface for the service template API
type Templater interface {
	GetTemplate(id string) (template StdTemplate, found bool)
	SetTemplateString(id string, body string) error
	SetTemplateFile(id string, file string) error
}
