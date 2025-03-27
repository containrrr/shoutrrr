package ntfy

import (
	"github.com/nicholas-fedor/shoutrrr/pkg/format"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
)

// Priority levels as constants.
const (
	PriorityMin     priority = 1
	PriorityLow     priority = 2
	PriorityDefault priority = 3
	PriorityHigh    priority = 4
	PriorityMax     priority = 5
)

type priority int

type priorityVals struct {
	Min     priority
	Low     priority
	Default priority
	High    priority
	Max     priority
	Enum    types.EnumFormatter
}

// Priority defines the notification priority levels.
var Priority = &priorityVals{
	Min:     PriorityMin,
	Low:     PriorityLow,
	Default: PriorityDefault,
	High:    PriorityHigh,
	Max:     PriorityMax,
	Enum: format.CreateEnumFormatter(
		[]string{
			"",
			"Min",
			"Low",
			"Default",
			"High",
			"Max",
		}, map[string]int{
			"1":      int(PriorityMin),
			"2":      int(PriorityLow),
			"3":      int(PriorityDefault),
			"4":      int(PriorityHigh),
			"5":      int(PriorityMax),
			"urgent": int(PriorityMax),
		}),
}

func (p priority) String() string {
	return Priority.Enum.Print(int(p))
}
