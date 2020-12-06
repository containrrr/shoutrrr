package hangouts

import (
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
)

// Config for use within the Hangouts Chat plugin.
type Config struct {
	standard.EnumlessConfig
	URL *url.URL
}

// SetURL updates a ServiceConfig from a URL representation of it's field values.
func (config *Config) SetURL(url *url.URL) error {
	config.URL = url
	config.URL.Scheme = "https"

	return nil
}

const (
	// Scheme is the identifying part of this service's configuration URL.
	Scheme = "hangouts"
)
