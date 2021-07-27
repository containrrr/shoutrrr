package slack

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/types"
	"net/url"
	"regexp"
	"strings"
)

var _ types.ConfigProp = &Token{}

const (
	hookTokenIdentifier = "hook"
	userTokenIdentifier = "xoxp"
	botTokenIdentifier  = "xoxb"
)

// Token is a Slack API token or a Slack webhook token
type Token struct {
	raw string
}

// SetFromProp updates it's state according to the passed string
// (implementation of the types.ConfigProp interface)
func (token *Token) SetFromProp(propValue string) error {
	if len(propValue) < 3 {
		return ErrorInvalidToken
	}

	match := tokenPattern.FindStringSubmatch(propValue)
	if match == nil || len(match) != tokenMatchCount {
		return ErrorInvalidToken
	}

	typeIdentifier := match[tokenMatchType]
	if typeIdentifier == "" {
		typeIdentifier = hookTokenIdentifier
	}

	token.raw = fmt.Sprintf("%s:%s-%s-%s",
		typeIdentifier, match[tokenMatchPart1], match[tokenMatchPart2], match[tokenMatchPart3])

	if match[tokenMatchSep1] != match[tokenMatchSep2] {
		return ErrorMismatchedTokenSeparators
	}

	return nil
}

// GetPropValue returns a deserializable string representation of the token
// (implementation of the types.ConfigProp interface)
func (token *Token) GetPropValue() (string, error) {
	if token == nil {
		return "", nil
	}

	return token.raw, nil
}

// TypeIdentifier returns the type identifier of the token
func (token Token) TypeIdentifier() string {
	return token.raw[:4]
}

// ParseToken parses and normalizes a token string
func ParseToken(str string) (*Token, error) {
	token := &Token{}
	if err := token.SetFromProp(str); err != nil {
		return nil, err
	}
	return token, nil
}

const (
	tokenMatchFull = iota
	tokenMatchType
	tokenMatchPart1
	tokenMatchSep1
	tokenMatchPart2
	tokenMatchSep2
	tokenMatchPart3
	tokenMatchCount
)

var tokenPattern = regexp.MustCompile(`(?:(?P<type>xox.|hook)[-:]|:?)(?P<p1>[A-Z0-9]{9,})(?P<s1>[-/,])(?P<p2>[A-Z0-9]{9,})(?P<s2>[-/,])(?P<p3>[A-Za-z0-9]{24,})`)

// String returns the token in normalized format with dashes (-) as separator
func (token *Token) String() string {
	return token.raw
}

// UserInfo returns a url.Userinfo struct populated from the token
func (token *Token) UserInfo() *url.Userinfo {
	return url.UserPassword(token.raw[:4], token.raw[5:])
}

// IsAPIToken returns whether the identifier is set to anything else but the webhook identifier (`hook`)
func (token *Token) IsAPIToken() bool {
	return token.TypeIdentifier() != hookTokenIdentifier
}

const webhookBase = "https://hooks.slack.com/services/"

// WebhookURL returns the corresponding Webhook URL for the Token
func (token Token) WebhookURL() string {
	sb := strings.Builder{}
	sb.WriteString(webhookBase)
	sb.Grow(len(token.raw) - 5)
	for i := 5; i < len(token.raw); i++ {
		c := token.raw[i]
		if c == '-' {
			c = '/'
		}
		sb.WriteByte(c)
	}
	return sb.String()
}

// Authorization returns the corresponding `Authorization` HTTP header value for the Token
func (token *Token) Authorization() string {
	sb := strings.Builder{}
	sb.WriteString("Bearer ")
	sb.Grow(len(token.raw))
	sb.WriteString(token.raw[:4])
	sb.WriteRune('-')
	sb.WriteString(token.raw[5:])
	return sb.String()
}
