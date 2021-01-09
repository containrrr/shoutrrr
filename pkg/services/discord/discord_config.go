package discord

import (
	"errors"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types/rich"
	"net/url"
)

// Config is the configuration needed to send discord notifications
type Config struct {
	standard.EnumlessConfig
	Channel    string
	Token      string
	Title      string `key:"title" default:""`
	Color      int    `key:"color" default:"50D9FF" desc:"The color of the left border for plain messages" base:"16"`
	ColorError int    `key:"colorError" default:"D60510" desc:"The color of the left border for error messages" base:"16"`
	ColorWarn  int    `key:"colorWarn" default:"FFC441" desc:"The color of the left border for warning messages" base:"16"`
	ColorInfo  int    `key:"colorInfo" default:"2488FF" desc:"The color of the left border for info messages" base:"16"`
	ColorDebug int    `key:"colorDebug" default:"7B00AB" desc:"The color of the left border for debug messages" base:"16"`
	SplitLines bool   `key:"splitLines" default:"yes" desc:"Whether to send each line as a separate embedded item"`
	JSON       bool   `desc:"Whether to send the whole message as the JSON payload instead of using it as the 'content' field"`
}

func (config *Config) LevelColors() (colors [rich.MessageLevelCount]int) {
	colors[rich.Unknown] = config.Color
	colors[rich.Error] = config.ColorError
	colors[rich.Warning] = config.ColorWarn
	colors[rich.Info] = config.ColorInfo
	colors[rich.Debug] = config.ColorDebug

	return colors
}

// GetURL returns a URL representation of it's current field values
func (config *Config) GetURL() *url.URL {
	return &url.URL{
		User:       url.User(config.Token),
		Host:       config.Channel,
		Scheme:     Scheme,
		ForceQuery: false,
	}
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *Config) SetURL(url *url.URL) error {

	config.Channel = url.Host
	config.Token = url.User.Username()

	if len(url.Path) > 0 {
		switch url.Path {
		case "/raw":
			config.JSON = true
			break
		default:
			return errors.New("illegal argument in config URL")
		}
	}

	if config.Channel == "" {
		return errors.New("channel missing from config URL")
	}

	if len(config.Token) < 1 {
		return errors.New("token missing from config URL")
	}

	return nil
}

// Scheme is the identifying part of this service's configuration URL
const Scheme = "discord"
