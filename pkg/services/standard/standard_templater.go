package standard

import (
	"strings"
	"text/template"
)

// Templater is the standard implementation of ApplyTemplate using the "text/template" library
type Templater struct {}

// ApplyTemplate applies the params to the input template and returns the result
func (templater *Templater) ApplyTemplate (tpl string, params *map[string]string) (string, error) {
	engine, err := template.New("").Parse(tpl);
	if err != nil {
		return "", err
	}

	writer := strings.Builder{}

	if err := engine.Execute(&writer, params); err != nil {
		return "", err
	}

	return writer.String(), err

}