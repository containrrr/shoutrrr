package ref

import (
	"github.com/containrrr/shoutrrr/pkg/format"
)

// ColorizeToken colorizes the value according to the tokenType
func ColorizeToken(value string, tokenType NodeTokenType) string {
	switch tokenType {
	case NumberToken:
		return format.ColorizeNumber(value)
	case EnumToken:
		return format.ColorizeEnum(value)
	case TrueToken:
		return format.ColorizeTrue(value)
	case FalseToken:
		return format.ColorizeFalse(value)
	case PropToken:
		return format.ColorizeProp(value)
	case ErrorToken:
		return format.ColorizeError(value)
	case ContainerToken:
		return format.ColorizeContainer(value)
	case StringToken:
		return format.ColorizeString(value)
	case UnknownToken:
	default:
	}
	return value
}
