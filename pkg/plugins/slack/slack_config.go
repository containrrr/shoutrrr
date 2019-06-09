package slack

import (
	"errors"
	"github.com/containrrr/shoutrrr/pkg/plugin"
	netUrl "net/url"
)

// Config for the slack plugin
type Config struct {
	BotName string
	Token   Token
	Channel string
}

func (config *Config) QueryFields() []string {
	return []string{}
}

func (config *Config) Enums() map[string]plugin.EnumFormatter {
	return map[string]plugin.EnumFormatter{}
}

func (config *Config) Get(string) (string, error) {
	return "", nil
}

func (config *Config) Set(string, string) error {
	return nil
}

func (config *Config) GetURL() netUrl.URL {
	return netUrl.URL{
		User: netUrl.UserPassword(config.BotName, config.Token.String()),
		Host: config.Channel,
		Scheme: Scheme,
		ForceQuery: false,
	}
}

func (config *Config) SetURL(url netUrl.URL) error {

	password, _ := url.User.Password()

	config.Channel = url.Host
	config.Token = ParseToken(password)

	if len(config.Channel) < 1 {
		return errors.New("channel missing from config URL")
	}

	if err := validateToken(config.Token); err != nil {
		return err
	}

	return nil
}

const (
	// DefaultUser for sending notifications to slack
	DefaultUser = "Shoutrrr"
	Scheme = "slack"
)

// CreateConfigFromURL to use within the slack plugin
func CreateConfigFromURL(url netUrl.URL) (*Config, error) {
	config := Config{}
	err := config.SetURL(url)
	return &config, err
}