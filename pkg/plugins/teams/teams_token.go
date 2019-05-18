package teams

import (
	"fmt"
	"regexp"
)

var uuid4_pattern = "[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}"

type TeamsToken struct {
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
	pattern := fmt.Sprintf("%s@%s", uuid4_pattern, uuid4_pattern)
	return matchesRegexp(pattern, token)
}

func isTokenBValid(token string) bool {
	return matchesRegexp("[A-Za-z0-9]{32}", token)
}

func isTokenCValid(token string) bool {
	return matchesRegexp(uuid4_pattern, token)
}

func matchesRegexp(pattern string, token string) bool {
	matched, err := regexp.MatchString(pattern, token)
	return !matched || err != nil
}
