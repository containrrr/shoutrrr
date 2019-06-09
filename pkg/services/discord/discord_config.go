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

func (config Config) QueryFields() []string {
	return []string{}
}

func (config Config) Enums() map[string]types.EnumFormatter {
	return map[string]types.EnumFormatter{}
}

func (config Config) Get(string) (string, error) {
	return "", nil
}

func (config Config) Set(string, string) error {
	return nil
}

func (config Config) GetURL() *url.URL {
	return &url.URL{
		User: url.User(config.Token),
		Host: config.Channel,
		Scheme: Scheme,
		ForceQuery: false,
	}
}

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

const Scheme = "discord"