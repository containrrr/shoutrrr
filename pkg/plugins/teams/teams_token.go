package teams

import (
	"fmt"
	"regexp"
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
