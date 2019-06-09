package discord

import (
	"errors"
	"github.com/containrrr/shoutrrr/pkg/plugin"
	"net/url"
)

// Config is the configuration needed to send discord notifications
type Config struct {
	Channel string
	Token string
}

func (config Config) QueryFields() []string {
	return []string{}
}

func (config Config) Enums() map[string]plugin.EnumFormatter {
	return map[string]plugin.EnumFormatter{}
}

func (config Config) Get(string) (string, error) {
	return "", nil
}

func (config Config) Set(string, string) error {
	return nil
}

func (config Config) GetURL() url.URL {

	return url.URL{
		User: url.UserPassword("Token", config.Token),
		Host: config.Channel,
		Scheme: Scheme,
		ForceQuery: false,
	}

}

func (config Config) SetURL(url url.URL) error {

	password, _ := url.User.Password()

	config.Channel = url.Host
	config.Token = password

	if len(config.Channel) < 1 {
		return errors.New("channel missing from config URL")
	}

	if len(config.Token) < 1 {
		return errors.New("token missing from config URL")
	}

	return nil
}

const Scheme = "discord"