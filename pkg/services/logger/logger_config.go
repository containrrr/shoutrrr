package logger

import (
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
)

// Config is the configuration object for the Logger Service
type Config struct {
	standard.EnumlessConfig
}

// GetURL returns a URL representation of it's current field values
func (config *Config) GetURL() *url.URL {
	return &url.URL{
		Scheme: Scheme,
	}
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *Config) SetURL(_ *url.URL) error {
	return nil
}

const (
	// Scheme is the identifying part of this service's configuration URL
	Scheme = "logger"
)
