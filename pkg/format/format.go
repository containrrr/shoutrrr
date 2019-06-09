package format

import "strings"

// NotifyFormat describes the format used in the notification body
type NotifyFormat int

const (
    // Markdown is the default notification format
    Markdown NotifyFormat = 0
)


func ParseBool(value string, defaultValue bool)  bool {
    switch strings.ToLower(value) {
    case "true": fallthrough
    case "1": fallthrough
    case "yes": return true
    case "false": fallthrough
    case "0": fallthrough
    case "no": return false
    default:
        return  defaultValue
    }
}

func PrintBool(value bool) string {
    if value {
        return "Yes"
    }

    return "No"

}