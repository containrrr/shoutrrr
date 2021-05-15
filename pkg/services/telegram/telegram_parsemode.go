package telegram

import (
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
)

type parseMode int

type parseModeVals struct {
	None       parseMode
	Markdown   parseMode
	HTML       parseMode
	MarkdownV2 parseMode
	Enum       types.EnumFormatter
}

// ParseModes is an enum helper for parseMode
var ParseModes = &parseModeVals{
	None:       0,
	Markdown:   1,
	HTML:       2,
	MarkdownV2: 3,
	Enum: format.CreateEnumFormatter(
		[]string{
			"None",
			"Markdown",
			"HTML",
			"MarkdownV2",
		}),
}

func (pm parseMode) String() string {
	return ParseModes.Enum.Print(int(pm))
}
