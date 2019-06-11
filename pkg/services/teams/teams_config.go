package teams

import (
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"net/url"
)

// Config for use within the teams plugin
type Config struct {
	standard.QuerylessConfig
	Token Token
}

// GetURL returns a URL representation of it's current field values
func (config *Config) GetURL() *url.URL {
	return &url.URL{
		User: url.UserPassword("Token", config.Token.String()),
		Host: "Teams",
		Scheme: Scheme,
		ForceQuery: false,
	}
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *Config) SetURL(url *url.URL) error {

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
func (service *Service) CreateConfigFromURL(url *url.URL) (*Config, error) {
	config := Config{}
	err := config.SetURL(url)
	return &config, err
}

const (
	// Scheme is the identifying part of this service's configuration URL
	Scheme = "teams"
)