package generators

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/generators/basic"
	t "github.com/containrrr/shoutrrr/pkg/types"
	"strings"
)

var generatorMap = map[string]func() t.Generator{
	"basic": func() t.Generator { return &basic.Generator{} },
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
