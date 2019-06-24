package pushbullet

import (
	"errors"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
)

// Config ...
type Config struct {
	standard.QuerylessConfig
	standard.EnumlessConfig
	Targets[] string
	Token string
}

// GetURL returns a URL representation of it's current field values
func (config *Config) GetURL() *url.URL {
	return &url.URL{
		Host: config.Token,
		Scheme: Scheme,
		ForceQuery: false,
	}
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *Config) SetURL(url *url.URL) error {


	return errors.New("not implemented")
}