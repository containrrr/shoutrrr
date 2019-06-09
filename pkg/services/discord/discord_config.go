package discord

import (
	"errors"
	"github.com/containrrr/shoutrrr/pkg/types"
	"net/url"
)

// Config is the configuration needed to send discord notifications
type Config struct {
	Channel string
	Token string
}

// QueryFields returns the fields that are part of the Query of the service URL
func (config Config) QueryFields() []string {
	return []string{}
}

// Enums returns the fields that should use a corresponding EnumFormatter to Print/Parse their values
func (config Config) Enums() map[string]types.EnumFormatter {
	return map[string]types.EnumFormatter{}
}

// Get returns the value of a Query field
func (config Config) Get(string) (string, error) {
	return "", nil
}

// Set updates the value of a Query field
func (config Config) Set(string, string) error {
	return nil
}

// GetURL returns a URL representation of it's current field values
func (config Config) GetURL() *url.URL {
	return &url.URL{
		User: url.User(config.Token),
		Host: config.Channel,
		Scheme: Scheme,
		ForceQuery: false,
	}
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config Config) SetURL(url *url.URL) error {

	config.Channel = url.Host
	config.Token = url.User.Username()

	if len(url.Path) > 0 {
		return errors.New("illegal argument in config URL")
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