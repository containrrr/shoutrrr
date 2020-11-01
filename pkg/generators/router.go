package generators

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/generators/basic"
	"github.com/containrrr/shoutrrr/pkg/generators/xouath2"
	t "github.com/containrrr/shoutrrr/pkg/types"
	"strings"
)

var generatorMap = map[string]func() t.Generator{
	"basic": func() t.Generator { return &basic.Generator{} },
	"oauth2": func() t.Generator { return &xouath2.Generator{} },
}

func NewGenerator(identifier string) (t.Generator, error) {
	generatorFactory, valid := generatorMap[strings.ToLower(identifier)]
	if !valid {
		return nil, fmt.Errorf("unknown generator %q", identifier)
	}
	return generatorFactory(), nil
}

func ListGenerators() []string {
	generators := make([]string, len(generatorMap))

	i := 0
	for key := range generatorMap {
		generators[i] = key
		i++
	}

	return generators
}
