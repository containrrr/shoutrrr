package gotify

import (
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
)

// Config for use within the gotify plugin
type Config struct {
	standard.EnumlessConfig
	Token    string
	Host     string
	Priority int    `key:"priority" default:"0"`
	Title    string `key:"title" default:"Shoutrrr notification"`
}

// GetURL returns a URL representation of it's current field values
func (config *Config) GetURL() *url.URL {
	resolver := format.NewPropKeyResolver(config)
	return config.getURL(&resolver)
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *Config) SetURL(url *url.URL) error {
	resolver := format.NewPropKeyResolver(config)
	return config.setURL(&resolver, url)
}

// GetURL returns a URL representation of it's current field values
func (config *Config) getURL(resolver types.ConfigQueryResolver) *url.URL {
	return &url.URL{
		Host:       config.Host,
		Scheme:     Scheme,
		ForceQuery: false,
		Path:       config.Token,
		RawQuery:   format.BuildQuery(resolver),
	}
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *Config) setURL(resolver types.ConfigQueryResolver, url *url.URL) error {
	config.Host = url.Host
	config.Token = url.Path

	for key, vals := range url.Query() {
		if err := resolver.Set(key, vals[0]); err != nil {
			return err
		}
	}
	return nil
}

const (
	// Scheme is the identifying part of this service's configuration URL
	Scheme = "gotify"
)
