package types

import (
	"strings"
	"time"
)

// MessageLevel is used to denote the importance of a MessageItem
type MessageLevel int

const (
	// Unknown MessageLevel (default)
	Unknown MessageLevel = iota
	// Debug MessageLevel
	Debug
	// Info MessageLevel
	Info
	// Warning MessageLevel
	Warning
	// Error MessageLevel
	Error
	messageLevelCount
	// MessageLevelCount is the number of MessageLevel values
	MessageLevelCount = int(messageLevelCount)
)

func (level MessageLevel) String() string {
	switch level {
	case Debug:
		return "Debug"
	case Info:
		return "Info"
	case Warning:
		return "Warning"
	case Error:
		return "Error"
	case Unknown:
	default:
	}
	return "Unknown"
}

// MessageItem is a notification message with some additional meta data
type MessageItem struct {
	Text      string
	Timestamp *time.Time
	Level     MessageLevel
}

// ItemsToPlain joins together the MessageItems' Text using newlines
// Used implement the rich sender API by redirecting to the plain sender implementation
func ItemsToPlain(items []MessageItem) string {
	builder := strings.Builder{}
	for _, item := range items {
		builder.WriteString(item.Text)
		builder.WriteRune('\n')
	}
	return builder.String()
}
