package rich

import "time"

type MessageLevel int

const (
	Unknown MessageLevel = iota
	Debug
	Info
	Warning
	Error
	messageLevelCount
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

type MessageItem struct {
	Text      string
	Timestamp *time.Time
	Level     MessageLevel
}
