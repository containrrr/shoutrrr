package format

import (
	"github.com/containrrr/shoutrrr/pkg/types"
	"strings"
)

const EnumInvalid = -1

type EnumFormatter struct {
	names []string
}

func (ef EnumFormatter) Names() []string {
	return ef.names
}

func (ef EnumFormatter) Print(e int) string {
	if e >= len(ef.names) || e < 0 {
		return "Invalid"
	}
	return ef.names[e]
}

func (ef EnumFormatter) Parse(s string) int {
	target := strings.ToLower(s)
	for index, name := range ef.names {
		if target == strings.ToLower(name) {
			return index
		}
	}
	return EnumInvalid
}

func CreateEnumFormatter(names []string) types.EnumFormatter {
	return &EnumFormatter{
		names,
	}
}