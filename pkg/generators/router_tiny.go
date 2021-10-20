//go:build tiny
// +build tiny

package generators

import (
	"github.com/containrrr/shoutrrr/pkg/generators/basic"
	"github.com/containrrr/shoutrrr/pkg/generators/xouath2"
	t "github.com/containrrr/shoutrrr/pkg/types"
)

var generatorMap = map[string]func() t.Generator{
	"basic":  func() t.Generator { return &basic.Generator{} },
	"oauth2": func() t.Generator { return &xouath2.Generator{} },
}
