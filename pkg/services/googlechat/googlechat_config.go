//go:generate go run ../../../cmd/shoutrrr-gen --lang go ../../../spec/googlechat.yml
package googlechat

import (
	"errors"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Config for use within the Google Chat plugin.
type LegacyConfig struct {
	standard.EnumlessConfig
	Host  string `default:"chat.googleapis.com"`
	Path  string
	Token string
	Key   string
}

// GetURL returns a URL representation of it's current field values
func (config *LegacyConfig) GetURL() *url.URL {
	resolver := format.NewPropKeyResolver(config)
	return config.getURL(&resolver)
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *LegacyConfig) SetURL(url *url.URL) error {
	resolver := format.NewPropKeyResolver(config)
	return config.setURL(&resolver, url)
}

// SetURL updates a ServiceConfig from a URL representation of it's field values.
func (config *LegacyConfig) setURL(_ types.ConfigQueryResolver, serviceURL *url.URL) error {
	config.Host = serviceURL.Host
	config.Path = serviceURL.Path

	query := serviceURL.Query()
	config.Key = query.Get("key")
	config.Token = query.Get("token")

	if config.Key == "" {
		return errors.New("missing field 'key'")
	}

	if config.Key == "" {
		return errors.New("missing field 'token'")
	}

	return nil
}

func (config *LegacyConfig) getURL(_ types.ConfigQueryResolver) *url.URL {
	query := url.Values{}
	query.Set("key", config.Key)
	query.Set("token", config.Token)

	return &url.URL{
		Host:     config.Host,
		Path:     config.Path,
		RawQuery: query.Encode(),
		Scheme:   Scheme,
	}
}

const (
	// Scheme is the identifying part of this service's configuration URL.
	Scheme = "googlechat"
)
