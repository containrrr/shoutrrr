package telegram

import (
	"github.com/nicholas-fedor/shoutrrr/pkg/format"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
)

type parseMode int

const (
	ParseModeNone       parseMode = iota // 0
	ParseModeMarkdown                    // 1
	ParseModeHTML                        // 2
	ParseModeMarkdownV2                  // 3
)

type parseModeVals struct {
	None       parseMode
	Markdown   parseMode
	HTML       parseMode
	MarkdownV2 parseMode
	Enum       types.EnumFormatter
}

// ParseModes is an enum helper for parseMode.
var ParseModes = &parseModeVals{
	None:       ParseModeNone,
	Markdown:   ParseModeMarkdown,
	HTML:       ParseModeHTML,
	MarkdownV2: ParseModeMarkdownV2,
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
