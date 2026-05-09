package format

import (
	"strings"

	"github.com/containrrr/shoutrrr/pkg/types"
)

// EnumInvalid is the constant value that an enum gets assigned when it could not be parsed
const EnumInvalid = -1

// EnumFormatter is the helper methods for enum-like types
type EnumFormatter struct {
	names       []string
	firstOffset int
	aliases     map[string]int
}

// Names is the list of the valid Enum string values
func (ef EnumFormatter) Names() []string {
	return ef.names[ef.firstOffset:]
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
	if index, found := ef.aliases[s]; found {
		return index
	}
	return EnumInvalid
}

// CreateEnumFormatter creates a EnumFormatter struct
func CreateEnumFormatter(names []string, optAliases ...map[string]int) types.EnumFormatter {
	aliases := map[string]int{}
	if len(optAliases) > 0 {
		aliases = optAliases[0]
	}
	firstOffset := 0
	for i, name := range names {
		if name != "" {
			firstOffset = i
			break
		}
	}
	return &EnumFormatter{
		names,
		firstOffset,
		aliases,
	}
}
