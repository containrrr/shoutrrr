package mattermost

import (
	"errors"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"net/url"
	"strings"
)

//Config object holding all information
type Config struct {
	standard.EnumlessConfig
	UserName string
	Channel  string
	Host     string
	Token    string
}

// GetURL returns a URL representation of it's current field values
func (config *Config) GetURL() *url.URL {
	return &url.URL{
		Host:       config.Host,
		Path:       fmt.Sprintf("/hooks/%s", config.Token),
		Scheme:     Scheme,
		ForceQuery: false,
	}
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *Config) SetURL(serviceURL *url.URL) error {

	config.Host = serviceURL.Hostname()
	if serviceURL.Path == "" || serviceURL.Path == "/" {
		return errors.New(string(NotEnoughArguments))
	}
	config.UserName = serviceURL.User.Username()
	path := strings.Split(serviceURL.Path[1:], "/")

	if len(path) < 1 {
		return errors.New(string(NotEnoughArguments))
	}

	config.Token = path[0]
	if len(path) > 1 {
		if path[1] != "" {
			config.Channel = path[1]
		}
	}

	return nil
}

//ErrorMessage for error events within the mattermost service
type ErrorMessage string

const (
	// Scheme is the identifying part of this service's configuration URL
	Scheme = "mattermost"
	// NotEnoughArguments provided in the service URL
	NotEnoughArguments ErrorMessage = "the apiURL does not include enough arguments, either provide 1 or 3 arguments (they may be empty)"
)
