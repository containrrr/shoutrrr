package rocketchat

import (
	"errors"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/types"
	"net/url"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
)

// Config for the rocket.chat service
type Config struct {
	standard.EnumlessConfig
	UserName string
	Host     string
	Port     string
	TokenA   string
	Channel  string
	TokenB   string
}

// GetURL returns a URL representation of it's current field values
func (config *Config) GetURL() *url.URL {
	return &url.URL{
		Host:       config.Host,
		Path:       fmt.Sprintf("hooks/%s/%s", config.TokenA, config.TokenB),
		Scheme:     Scheme,
		ForceQuery: false,
	}
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *Config) SetURL(serviceURL *url.URL) error {

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
