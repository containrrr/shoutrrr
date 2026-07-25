package pinglet

import (
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
)

type priority int

type priorityVals struct {
	Silent priority
	Normal priority
	Urgent priority
	Enum   types.EnumFormatter
}

// Priority is the enum of delivery priorities supported by Pinglet
var Priority = &priorityVals{
	Silent: 0,
	Normal: 1,
	Urgent: 2,
	Enum: format.CreateEnumFormatter(
		[]string{
			"silent",
			"normal",
			"urgent",
		}, map[string]int{
			// Prefix aliases (s => silent, n => normal, u => urgent)
			"s": 0,
			"n": 1,
			"u": 2,
		}),
}

func (p priority) String() string {
	return Priority.Enum.Print(int(p))
}
