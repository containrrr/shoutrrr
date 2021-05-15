package format

import "github.com/fatih/color"

// ColorizeDesc colorizes the input string as "Description"
var ColorizeDesc = color.New(color.FgHiBlack).SprintFunc()

// ColorizeTrue colorizes the input string as "True"
var ColorizeTrue = color.New(color.FgHiGreen).SprintFunc()

// ColorizeFalse colorizes the input string as "False"
var ColorizeFalse = color.New(color.FgHiRed).SprintFunc()

// ColorizeNumber colorizes the input string as "Number"
var ColorizeNumber = color.New(color.FgHiBlue).SprintFunc()

// ColorizeString colorizes the input string as "String"
var ColorizeString = color.New(color.FgHiYellow).SprintFunc()

// ColorizeEnum colorizes the input string as "Enum"
var ColorizeEnum = color.New(color.FgHiCyan).SprintFunc()

// ColorizeProp colorizes the input string as "Prop"
var ColorizeProp = color.New(color.FgHiMagenta).SprintFunc()

// ColorizeError colorizes the input string as "Error"
var ColorizeError = ColorizeFalse

// ColorizeContainer colorizes the input string as "Container"
var ColorizeContainer = ColorizeDesc

// ColorizeLink colorizes the input string as "Link"
var ColorizeLink = color.New(color.FgHiBlue).SprintFunc()

// ColorizeValue colorizes the input string according to what type appears to be
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

// ColorizeToken colorizes the value according to the tokenType
func ColorizeToken(value string, tokenType NodeTokenType) string {
	switch tokenType {
	case NumberToken:
		return ColorizeNumber(value)
	case EnumToken:
		return ColorizeEnum(value)
	case TrueToken:
		return ColorizeTrue(value)
	case FalseToken:
		return ColorizeFalse(value)
	case PropToken:
		return ColorizeProp(value)
	case ErrorToken:
		return ColorizeError(value)
	case ContainerToken:
		return ColorizeContainer(value)
	case StringToken:
		return ColorizeString(value)
	case UnknownToken:
	default:
	}
	return value
}
