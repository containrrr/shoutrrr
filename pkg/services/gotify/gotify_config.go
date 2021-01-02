package gotify

import (
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
)

// Config for use within the gotify plugin
type Config struct {
	standard.EnumlessConfig
	Token   string
	Host    string
	Prority int
}

// GetURL returns a URL representation of it's current field values
func (config *Config) GetURL() *url.URL {
	return &url.URL{
		Host:       config.Host,
		Scheme:     Scheme,
		ForceQuery: false,
		Path:       config.Token,
	}
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *Config) SetURL(url *url.URL) error {
	config.Host = url.Hostname()
	config.Token = url.Path
	return nil
}

const (
	// Scheme is the identifying part of this service's configuration URL
	Scheme = "gotify"
)
