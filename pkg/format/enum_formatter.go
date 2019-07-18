package format

import (
	"github.com/containrrr/shoutrrr/pkg/types"
	"strings"
)

// EnumInvalid is the constant value that an enum gets assigned when it could not be parsed
const EnumInvalid = -1

// EnumFormatter is the helper methods for enum-like types
type EnumFormatter struct {
	names []string
}

// Names is the list of the valid Enum string values
func (ef EnumFormatter) Names() []string {
	return ef.names
}

// Print takes a enum mapped int and returns it's string representation or "Invalid"
func (ef EnumFormatter) Print(e int) string {
	if e >= len(ef.names) || e < 0 {
		return "Invalid"
	}
	return ef.names[e]
}

// Parse takes an enum mapped string and returns it's int representation or EnumInvalid (-1)
func (ef EnumFormatter) Parse(s string) int {
	target := strings.ToLower(s)
	for index, name := range ef.names {
		if target == strings.ToLower(name) {
			return index
		}
	}
	return EnumInvalid
}

// CreateEnumFormatter creates a EnumFormatter struct
func CreateEnumFormatter(names []string) types.EnumFormatter {
	return &EnumFormatter{
		names,
	}
}
