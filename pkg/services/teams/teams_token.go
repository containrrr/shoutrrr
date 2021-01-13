package teams

import (
	"fmt"
	"regexp"
)

var uuid4Pattern = "[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}"
var hex32Pattern = "[A-Za-z0-9]{32}"

func Part124Valid(token string) bool {
	return matchesRegexp(uuid4Pattern, token)
}

func Part3Valid(token string) bool {
	return matchesRegexp(hex32Pattern, token)
}

func VerifyWebhookParts(p [4]string) error {
	if !Part124Valid(p[0]) {
		return fmt.Errorf("first token part is invalid: '%v'", p[0])
	}
	if !Part124Valid(p[1]) {
		return fmt.Errorf("second token part is invalid: '%v'", p[1])
	}
	if !Part3Valid(p[2]) {
		return fmt.Errorf("third token part is invalid: '%v'", p[2])
	}
	if !Part124Valid(p[3]) {
		return fmt.Errorf("forth token part is invalid: '%v'", p[3])
	}
	return nil
}

func matchesRegexp(pattern string, token string) bool {
	matched, err := regexp.MatchString(pattern, token)
	return matched && err == nil
}
