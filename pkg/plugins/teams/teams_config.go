package teams

import (
	"github.com/containrrr/shoutrrr/pkg/plugin"
	"net/url"
)

// Config for use within the teams plugin
type Config struct {
	Token Token
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
		User: url.UserPassword("Token", config.Token.String()),
		Host: "Teams",
		Scheme: Scheme,
		ForceQuery: false,
	}
}

func (config Config) SetURL(url url.URL) error {

	password, _ := url.User.Password()

	if token, err := ParseToken(password); err != nil {
		return err
	} else {
		config.Token = token
	}

	return nil
}

// CreateConfigFromURL for use within the teams plugin
func (plugin *Plugin) CreateConfigFromURL(url url.URL) (*Config, error) {
	config := Config{}
	err := config.SetURL(url)
	return &config, err
}

const (
	Scheme = "teams"
)