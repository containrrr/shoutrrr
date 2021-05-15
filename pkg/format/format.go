package format

import (
	"strconv"
	"strings"
)

// ParseBool returns true for "1","true","yes" or false for "0","false","no" or defaultValue for any other value
func ParseBool(value string, defaultValue bool) (parsedValue bool, ok bool) {
	switch strings.ToLower(value) {
	case "true", "1", "yes", "y":
		return true, true
	case "false", "0", "no", "n":
		return false, true
	default:
		return defaultValue, false
	}
}

// PrintBool returns "Yes" if value is true, otherwise returns "No"
func PrintBool(value bool) string {
	if value {
		return "Yes"
	}

	return "No"

}

// IsNumber returns whether the specified string is number-like
func IsNumber(value string) bool {
	_, err := strconv.ParseFloat(value, 64)
	return err == nil
}
