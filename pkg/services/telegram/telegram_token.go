package telegram

import "regexp"

// IsTokenValid for use with telegram
func IsTokenValid(token string) bool {
	matched, err := regexp.MatchString("^[0-9]+:[a-zA-Z0-9_-]+$", token)
	return matched && err == nil
}
