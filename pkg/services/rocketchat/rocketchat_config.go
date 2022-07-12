//go:generate go run ../../../cmd/shoutrrr-gen
package rocketchat

import (
	"net/url"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/conf"
	"github.com/containrrr/shoutrrr/pkg/types"
)

const (
	// Scheme is the identifying part of this service's configuration URL
	Scheme = "rocketchat"
	// NotEnoughArguments provided in the service URL
	NotEnoughArguments = "the apiURL does not include enough arguments"
)

// CreateConfigFromURL to use within the rocket.chat service
func CreateConfigFromURL(_ types.ConfigQueryResolver, serviceURL *url.URL) (*Config, error) {
	config := Config{}
	err := config.SetURL(serviceURL)
	return &config, err
}

func (*Config) UpdateLegacyURL(legacyURL *url.URL) *url.URL {
	updatedURL := *legacyURL

	path := conf.SplitPath(legacyURL.Path)
	channel := ""

	if legacyURL.Fragment != "" {
		channel = "#" + legacyURL.Fragment
		legacyURL.Fragment = ""
		if len(path) < 3 {
			path = append(path, channel)
		} else {
			path[2] = channel
		}
	} else if len(path) > 2 {
		if !strings.HasPrefix(path[2], "@") {
			path[2] = "#" + path[2]
		}
	}

	updatedURL.Path = conf.JoinPath(path...)
	return &updatedURL
}
