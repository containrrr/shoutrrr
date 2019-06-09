package plugin

import "strings"

const EnumInvalid = -1

type EnumFormatter struct {
	Names []string
}

func (ef EnumFormatter) Print(e int) string {
	if e >= len(ef.Names) || e < 0 {
		return "Invalid"
	} else {
		return ef.Names[e]
	}
}

func (ef EnumFormatter) Parse(s string) int {
	target := strings.ToLower(s)
	for index, name := range ef.Names {
		if target == strings.ToLower(name) {
			return index
		}
	}
	return EnumInvalid
}