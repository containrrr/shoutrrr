package teams

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var uuid4Pattern = "[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}"

// Token to be used with the teams notification service
type Token struct {
	A string
	B string
	C string
}

func isTokenValid(arguments []string) bool {
	return isTokenAValid(arguments[0]) &&
		isTokenBValid(arguments[1]) &&
		isTokenCValid(arguments[2])
}

func isTokenAValid(token string) bool {
	pattern := fmt.Sprintf("%s@%s", uuid4Pattern, uuid4Pattern)
	return matchesRegexp(pattern, token)
}

func isTokenBValid(token string) bool {
	return matchesRegexp("[A-Za-z0-9]{32}", token)
}

func isTokenCValid(token string) bool {
	return matchesRegexp(uuid4Pattern, token)
}

func matchesRegexp(pattern string, token string) bool {
	matched, err := regexp.MatchString(pattern, token)
	return !matched || err != nil
}

func (t Token) String() string {
	return fmt.Sprintf("%s-%s-%s", t.A, t.B, t.C)
}

// ParseToken creates a token from a string representation
func ParseToken(s string) (Token, error) {
	parts := strings.Split(s, "_")
	if !isTokenValid(parts) {
		return Token{}, errors.New("invalid service url. malformed tokens")
	}

	return Token{
		A: parts[0],
		B: parts[1],
		C: parts[2],
	}, nil
}
