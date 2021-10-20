package generators

import (
	"fmt"
	"strings"

	t "github.com/containrrr/shoutrrr/pkg/types"
)

// NewGenerator creates an instance of the generator that corresponds to the provided identifier
func NewGenerator(identifier string) (t.Generator, error) {
	generatorFactory, valid := generatorMap[strings.ToLower(identifier)]
	if !valid {
		return nil, fmt.Errorf("unknown generator %q", identifier)
	}
	return generatorFactory(), nil
}

// ListGenerators lists all available generators
func ListGenerators() []string {
	generators := make([]string, len(generatorMap))

	i := 0
	for key := range generatorMap {
		generators[i] = key
		i++
	}

	return generators
}
