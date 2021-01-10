package util

import "strings"

const hex int = 16

// StripNumberPrefix returns a number string with any base prefix stripped and it's corresponding base
// If no prefix was found, returns 0 to let strconv try to identify the base
func StripNumberPrefix(input string) (number string, base int) {

	if strings.HasPrefix(input, "#") {
		return input[1:], hex
	}

	return input, 0
}
