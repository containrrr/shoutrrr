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

var EncMethods = &encMethodVals{
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
	return EncMethods.Enum.Print(int(at))
}

func useImplicitTLS(encryption encMethod, port uint16) bool {
	switch encryption {
	case EncMethods.ImplicitTLS:
		return true
	case EncMethods.Auto:
		return port == ImplicitTLSPort
	default:
		return false
	}
}

const ImplicitTLSPort = 465
