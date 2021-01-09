package types

import (
	"time"
)

type MessageLevel uint8

const (
	Unknown MessageLevel = iota
	Debug
	Info
	Warning
	Error
	messageLevelCount
	MessageLevelCount = int(messageLevelCount)
)

var messageLevelStrings = [MessageLevelCount]string {
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

type MessageItem struct {
	Text      string
	Timestamp time.Time
	Level     MessageLevel
	Fields	  []Field
}

func (mi *MessageItem) WithField(key, value string) *MessageItem {
	mi.Fields = append(mi.Fields, Field{
		Key:   key,
		Value: value,
	})
	return mi
}
