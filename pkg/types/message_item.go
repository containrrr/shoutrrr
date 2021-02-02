package types

import (
	"strings"
	"time"
)

// MessageLevel is used to denote the urgency of a message item
type MessageLevel uint8

const (
	// Unknown is the default message level
	Unknown MessageLevel = iota
	// Debug is the lowest kind of known message level
	Debug
	// Info is generally used as the "normal" message level
	Info
	// Warning is generally used to denote messages that might be OK, but can cause problems
	Warning
	// Error is generally used for messages about things that did not go as planned
	Error
	messageLevelCount
	// MessageLevelCount is used to create arrays that maps levels to other values
	MessageLevelCount = int(messageLevelCount)
)

var messageLevelStrings = [MessageLevelCount]string{
	"Unknown",
	"Debug",
	"Info",
	"Warning",
	"Error",
}

func (level MessageLevel) String() string {
	if level >= messageLevelCount {
		return messageLevelStrings[0]
	}
	return messageLevelStrings[level]
}

// MessageItem is an entry in a notification being sent by a service
type MessageItem struct {
	Text      string
	Timestamp time.Time
	Level     MessageLevel
	Fields    []Field
}

// WithField appends the key/value pair to the message items fields
func (mi *MessageItem) WithField(key, value string) *MessageItem {
	mi.Fields = append(mi.Fields, Field{
		Key:   key,
		Value: value,
	})
	return mi
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
