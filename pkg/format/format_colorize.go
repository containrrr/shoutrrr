package format

import "github.com/fatih/color"

var ColorizeDesc = color.New(color.FgHiBlack).SprintFunc()
var ColorizeTrue = color.New(color.FgHiGreen).SprintFunc()
var ColorizeFalse = color.New(color.FgHiRed).SprintFunc()
var ColorizeNumber = color.New(color.FgHiBlue).SprintFunc()
var ColorizeString = color.New(color.FgHiYellow).SprintFunc()
var ColorizeEnum = color.New(color.FgHiCyan).SprintFunc()

func ColorizeValue(value string, isEnum bool) string {
	if isEnum {
		return ColorizeEnum(value)
	}

	if isTrue, isType := ParseBool(value, false); isType {
		if isTrue {
			return ColorizeTrue(value)
		}
		return ColorizeFalse(value)
	}

	if IsNumber(value) {
		return ColorizeNumber(value)
	}

	return ColorizeString(value)
}