package smtp

import (
	"github.com/nicholas-fedor/shoutrrr/pkg/format"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
)

type authType int

const (
	AuthNone    authType = iota // 0
	AuthPlain                   // 1
	AuthCRAMMD5                 // 2
	AuthUnknown                 // 3
	AuthOAuth2                  // 4
)

type authTypeVals struct {
	None    authType
	Plain   authType
	CRAMMD5 authType
	Unknown authType
	OAuth2  authType
	Enum    types.EnumFormatter
}

// AuthTypes is the enum helper for populating the Auth field.
var AuthTypes = &authTypeVals{
	None:    AuthNone,
	Plain:   AuthPlain,
	CRAMMD5: AuthCRAMMD5,
	Unknown: AuthUnknown,
	OAuth2:  AuthOAuth2,
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
	return AuthTypes.Enum.Print(int(at))
}
