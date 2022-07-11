//go:generate go run ../../../cmd/shoutrrr-gen
package rocketchat

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/conf"
	"github.com/containrrr/shoutrrr/pkg/types"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
)

// Config for the rocket.chat service
type LegacyConfig struct {
	standard.EnumlessConfig
	UserName string `url:"user" optional:""`
	Host     string `url:"host"`
	Port     string `url:"port"`
	TokenA   string `url:"path1"`
	Channel  string `url:"path3"`
	TokenB   string `url:"path2"`
}

// GetURL returns a URL representation of it's current field values
func (config *LegacyConfig) GetURL() *url.URL {

	u := &url.URL{
		Host:       fmt.Sprintf("%s:%v", config.Host, config.Port),
		Path:       fmt.Sprintf("%s/%s", config.TokenA, config.TokenB),
		Scheme:     Scheme,
		ForceQuery: false,
	}
	return u
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *LegacyConfig) SetURL(serviceURL *url.URL) error {

	UserName := serviceURL.User.Username()
	host := serviceURL.Hostname()

	path := strings.Split(serviceURL.Path, "/")

	if len(path) < 3 {
		return errors.New(NotEnoughArguments)
	}

	config.Port = serviceURL.Port()
	config.UserName = UserName
	config.Host = host
	config.TokenA = path[1]
	config.TokenB = path[2]
	if len(path) > 3 {
		if serviceURL.Fragment != "" {
			config.Channel = "#" + serviceURL.Fragment
		} else if !strings.HasPrefix(path[3], "@") {
			config.Channel = "#" + path[3]
		} else {
			config.Channel = path[3]
		}
	}
	return nil
}

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
