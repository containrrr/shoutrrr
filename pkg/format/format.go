package format

import (
    "strconv"
    "strings"
)

// NotifyFormat describes the format used in the notification body
type NotifyFormat int

const (
    // Markdown is the default notification format
    Markdown NotifyFormat = 0
)

// ParseBool returns true for "1","true","yes" or false for "0","false","no" or defaultValue for any other value
func ParseBool(value string, defaultValue bool)  (bool, bool) {
    switch strings.ToLower(value) {
    case "true": fallthrough
    case "1": fallthrough
    case "yes": return true, true
    case "false": fallthrough
    case "0": fallthrough
    case "no": return false, true
    default:
        return  defaultValue, false
    }
}

// PrintBool returns "Yes" if value is true, otherwise returns "No"
func PrintBool(value bool) string {
    if value {
        return "Yes"
    }

    return "No"

}

func IsNumber(value string) bool {
    _, err := strconv.ParseFloat(value, 64)
    return err == nil
}