package teams

import (
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
)

// Config for use within the teams plugin
type Config struct {
	standard.EnumlessConfig
	Token Token
}

// GetURL returns a URL representation of it's current field values
func (config *Config) GetURL() *url.URL {
	return &url.URL{
		User:       url.UserPassword(config.Token.A, config.Token.B),
		Host:       config.Token.C,
		Scheme:     Scheme,
		ForceQuery: false,
	}
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *Config) SetURL(url *url.URL) error {

	tokenA := url.User.Username()
	tokenB, _ := url.User.Password()
	tokenC := url.Hostname()

	config.Token = Token{
		A: tokenA,
		B: tokenB,
		C: tokenC,
	}
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
