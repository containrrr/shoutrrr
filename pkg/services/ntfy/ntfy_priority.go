package ntfy

import (
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
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

// Priority ...
var Priority = &priorityVals{
	Min:     1,
	Low:     2,
	Default: 3,
	High:    4,
	Max:     5,
	Enum: format.CreateEnumFormatter(
		[]string{
			"",
			"Min",
			"Low",
			"Default",
			"High",
			"Max",
		}, map[string]int{
			"1":      1,
			"2":      2,
			"3":      3,
			"4":      4,
			"5":      5,
			"urgent": 5,
		}),
}

func (p priority) String() string {
	return Priority.Enum.Print(int(p))
}
