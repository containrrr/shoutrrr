package teams

import (
	"github.com/containrrr/shoutrrr/pkg/types"
	"net/url"
)

// Config for use within the teams plugin
type Config struct {
	Token Token
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
		User: url.UserPassword("Token", config.Token.String()),
		Host: "Teams",
		Scheme: Scheme,
		ForceQuery: false,
	}
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config Config) SetURL(url *url.URL) error {

	password, _ := url.User.Password()

	var err error
	var token Token

	if token, err = ParseToken(password); err != nil {
		return err
	}

	config.Token = token
	return nil
}

// CreateConfigFromURL for use within the teams plugin
func (plugin *Service) CreateConfigFromURL(url *url.URL) (*Config, error) {
	config := Config{}
	err := config.SetURL(url)
	return &config, err
}

const (
	// Scheme is the identifying part of this service's configuration URL
	Scheme = "teams"
)