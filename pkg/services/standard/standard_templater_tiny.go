//go:build tinygo
// +build tinygo

package standard

import (
	"errors"
	"io"

	"github.com/containrrr/shoutrrr/pkg/types"
)

var NotSupportedError = errors.New("not supported in TinyGO")

// Templater is the standard implementation of ApplyTemplate using the "text/template" library
type Templater struct{}

type FakeTemplate struct{}

func (*FakeTemplate) Execute(wr io.Writer, data interface{}) error {
	return NotSupportedError
}

// GetTemplate attempts to retrieve the template identified with id
func (templater *Templater) GetTemplate(id string) (template types.StdTemplate, found bool) {
	return &FakeTemplate{}, false
}

// SetTemplateString creates a new template from the body and assigning it the id
func (templater *Templater) SetTemplateString(id string, body string) error {
	return NotSupportedError
}

// SetTemplateFile creates a new template from the file and assigning it the id
func (templater *Templater) SetTemplateFile(id string, file string) error {
	return NotSupportedError
}
