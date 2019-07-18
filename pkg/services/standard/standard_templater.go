package standard

import (
	"io/ioutil"
	"text/template"
)

// Templater is the standard implementation of ApplyTemplate using the "text/template" library
type Templater struct {
	templates map[string]*template.Template
}

// GetTemplate attempts to retrieve the template identified with id
func (templater *Templater) GetTemplate(id string) (template *template.Template, found bool) {
	tpl, found := templater.templates[id]
	return tpl, found
}

// SetTemplateString creates a new template from the body and assigning it the id
func (templater *Templater) SetTemplateString(id string, body string) error {
	tpl, err := template.New("").Parse(body)
	if err != nil {
		return err
	}
	if templater.templates == nil {
		templater.templates = make(map[string]*template.Template, 1)
	}

	templater.templates[id] = tpl
	return nil
}

// SetTemplateFile creates a new template from the file and assigning it the id
func (templater *Templater) SetTemplateFile(id string, file string) error {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return templater.SetTemplateString(id, string(bytes))
}
