//go:generate go run ../../../cmd/shoutrrr-gen
package zulip

import (
	"net/url"
)

const (
	// Scheme is the identifying part of this service's configuration URL
	Scheme = "zulip"
)

// CreateConfigFromURL to use within the zulip service
func CreateConfigFromURL(serviceURL *url.URL) (*Config, error) {
	config := Config{}
	err := config.SetURL(serviceURL)

	return &config, err
}
