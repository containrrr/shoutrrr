package standard

import (
	"github.com/containrrr/shoutrrr/pkg/types"
)

// EnumlessConfig implements the ServiceConfig interface for services that does not use Enum fields
type EnumlessConfig struct{}

// Enums returns an empty map
func (ec *EnumlessConfig) Enums() map[string]types.EnumFormatter {
	return map[string]types.EnumFormatter{}
}
