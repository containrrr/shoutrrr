package telegram

import (
	"errors"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/plugin"
	"net/url"
	"strings"
)

// Config for use within the telegram plugin
type Config struct {
	Token    string
	Channels []string
}

func (config Config) QueryFields() []string {
	return []string{
		"channels",
	}
}

func (config Config) Enums() map[string]plugin.EnumFormatter {
	return map[string]plugin.EnumFormatter{}
}

func (config Config) Get(key string) (string, error) {
	switch key {
	case "channels":
		return strings.Join(config.Channels, ","), nil
	}
	return "", fmt.Errorf("invalid query key \"%s\"", key)
}

func (config Config) Set(key string, value string) error {
	switch key {
	case "channels":
		config.Channels = strings.Split(value, ",")
	default:
		return fmt.Errorf("invalid query key \"%s\"", key)
	}
	return nil
}

func (config Config) GetURL() url.URL {

	return url.URL{
		User: url.UserPassword("Token", config.Token),
		Host: Scheme,
		Scheme: Scheme,
		ForceQuery: true,
		RawQuery: plugin.FormatQuery(config),
	}

}

func (config Config) SetURL(url url.URL) error {

	password, _ := url.User.Password()

	config.Token = password
	if !IsTokenValid(password) {
		return errors.New("invalid telegram token")
	}

	for key, vals := range url.Query() {
		if err := config.Set(key, vals[0]); err != nil {
			return err
		}
	}

	if len(config.Channels) < 1 {
		return errors.New("no channels defined in config URL")
	}

	return nil
}

const (
	Scheme = "telegram"
)