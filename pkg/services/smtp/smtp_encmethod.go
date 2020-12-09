package smtp

import (
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
)

type encMethod int

type encMethodVals struct {
	None        encMethod
	ExplicitTLS encMethod
	ImplicitTLS encMethod
	Auto        encMethod

	Enum types.EnumFormatter
}

var encMethods = &encMethodVals{
	None:        0,
	ExplicitTLS: 1,
	ImplicitTLS: 2,
	Auto:        3,

	Enum: format.CreateEnumFormatter(
		[]string{
			"None",
			"ExplicitTLS",
			"ImplicitTLS",
			"Auto",
		}),
}

func (at encMethod) String() string {
	return encMethods.Enum.Print(int(at))
}

func parseEncryption(s string) encMethod {
	return encMethod(encMethods.Enum.Parse(s))
}

func useImplicitTLS(encryption encMethod, port uint16) bool {
	switch encryption {
	case encMethods.ImplicitTLS:
		return true
	case encMethods.Auto:
		return port == ImplicitTLSPort
	default:
		return false
	}
}

const ImplicitTLSPort = 465
