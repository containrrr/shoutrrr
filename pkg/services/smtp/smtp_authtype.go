package smtp

import (
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
)

type authType int

type authTypeVals struct {
	None    authType
	Plain   authType
	CRAMMD5 authType
	Unknown authType
	OAuth2  authType
	Enum    types.EnumFormatter
}

var authTypes = &authTypeVals{
	None:    0,
	Plain:   1,
	CRAMMD5: 2,
	Unknown: 3,
	OAuth2:  4,
	Enum: format.CreateEnumFormatter(
		[]string{
			"None",
			"Plain",
			"CRAMMD5",
			"Unknown",
			"OAuth2",
		}),
}

func (at authType) String() string {
	return authTypes.Enum.Print(int(at))
}

func parseAuth(s string) authType {
	return authType(authTypes.Enum.Parse(s))
}
