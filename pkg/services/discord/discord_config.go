package discord

import (
	"errors"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
)

// Config is the configuration needed to send discord notifications
type Config struct {
	standard.EnumlessConfig
	Channel string
	Token   string
	JSON    bool
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
