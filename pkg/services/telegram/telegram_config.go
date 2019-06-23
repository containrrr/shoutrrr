package telegram

import (
	"errors"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
	"net/url"
	"strings"
)

// Config for use within the telegram plugin
type Config struct {
	Token    string
	Channels []string
}

// QueryFields returns the fields that are part of the Query of the service URL
func (config *Config) QueryFields() []string {
	return []string{
		"channels",
	}
}

// Enums returns the fields that should use a corresponding EnumFormatter to Print/Parse their values
func (config *Config) Enums() map[string]types.EnumFormatter {
	return map[string]types.EnumFormatter{}
}

// Get returns the value of a Query field
func (config *Config) Get(key string) (string, error) {
	switch key {
	case "channels":
		return strings.Join(config.Channels, ","), nil
	}
	return "", fmt.Errorf("invalid query key \"%s\"", key)
}

// Set updates the value of a Query field
func (config *Config) Set(key string, value string) error {
	switch key {
	case "channels":
		config.Channels = strings.Split(value, ",")
	default:
		return fmt.Errorf("invalid query key \"%s\"", key)
	}
	return nil
}

// GetURL returns a URL representation of it's current field values
func (config *Config) GetURL() *url.URL {

	return &url.URL{
		User:       url.UserPassword("Token", config.Token),
		Host:       Scheme,
		Scheme:     Scheme,
		ForceQuery: true,
		RawQuery:   format.BuildQuery(config),
	}

}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *Config) SetURL(url *url.URL) error {

	password, _ := url.User.Password()

	token := url.User.Username() + ":" + password
	if !IsTokenValid(token) {
		return fmt.Errorf("invalid telegram token %s", token)
	}


	for key, vals := range url.Query() {
		if err := config.Set(key, vals[0]); err != nil {
			return err
		}
	}

	if len(config.Channels) < 1 {
		return errors.New("no channels defined in config URL")
	}

	config.Token = token

	return nil
}

// Scheme is the identifying part of this service's configuration URL
const (
	Scheme = "telegram"
)