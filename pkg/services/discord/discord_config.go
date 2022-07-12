//go:generate go run ../../../cmd/shoutrrr-gen --lang go
package discord

import (
	"fmt"

	"github.com/containrrr/shoutrrr/pkg/types"
)

// LevelColors returns an array of colors with a MessageLevel index
func (config *Config) LevelColors() (colors [types.MessageLevelCount]uint) {
	colors[types.Unknown] = uint(config.Color)
	colors[types.Error] = uint(config.ColorError)
	colors[types.Warning] = uint(config.ColorWarn)
	colors[types.Info] = uint(config.ColorInfo)
	colors[types.Debug] = uint(config.ColorDebug)

	return colors
}

type rawModeType string

func (config *Config) getRawMode() string {
	if config.JSON {
		return "raw"
	} else {
		return ""
	}
}

func (config *Config) setRawMode(v string) (rawModeType, error) {
	if v == "raw" {
		config.JSON = true
		return rawModeType(v), nil
	} else if v == "" {
		return rawModeType(""), nil
	}

	return "", fmt.Errorf("invalid value raw mode value %q", v)
}

// Scheme is the identifying part of this service's configuration URL
const Scheme = "discord"
