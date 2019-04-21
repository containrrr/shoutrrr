package telegram

import "regexp"

func IsTokenValid(token string) bool {
	matched, err := regexp.MatchString("^[0-9]+:[a-zA-Z0-9_-]+$", token)
	return matched && err == nil
}