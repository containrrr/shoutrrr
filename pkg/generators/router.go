package generators

import (
	"fmt"
	"strings"

	"github.com/nicholas-fedor/shoutrrr/pkg/generators/basic"
	"github.com/nicholas-fedor/shoutrrr/pkg/generators/xouath2"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/telegram"
	t "github.com/nicholas-fedor/shoutrrr/pkg/types"
)

var generatorMap = map[string]func() t.Generator{
	"basic":    func() t.Generator { return &basic.Generator{} },
	"oauth2":   func() t.Generator { return &xouath2.Generator{} },
	"telegram": func() t.Generator { return &telegram.Generator{} },
}

// NewGenerator creates an instance of the generator that corresponds to the provided identifier.
func NewGenerator(identifier string) (t.Generator, error) {
	generatorFactory, valid := generatorMap[strings.ToLower(identifier)]
	if !valid {
		return nil, fmt.Errorf("unknown generator %q", identifier)
	}

	return generatorFactory(), nil
}

// ListGenerators lists all available generators.
func ListGenerators() []string {
	generators := make([]string, len(generatorMap))

	i := 0

	for key := range generatorMap {
		generators[i] = key
		i++
	}

	return generators
}
