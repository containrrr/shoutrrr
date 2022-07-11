//go:generate go run ../../../cmd/shoutrrr-gen
package mattermost

import (
	"errors"
	"net/url"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

//LegacyConfig object holding all information
type LegacyConfig struct {
	standard.EnumlessConfig
	UserName string `url:"user" optional:"" desc:"Override webhook user"`
	Channel  string `url:"path2" optional:"" desc:"Override webhook channel"`
	Host     string `url:"host,port" desc:"Mattermost server host"`
	Token    string `url:"path1" desc:"Webhook token"`
}

// GetURL returns a URL representation of it's current field values
func (config *LegacyConfig) GetURL() *url.URL {
	paths := []string{"", config.Token, config.Channel}
	if config.Channel == "" {
		paths = paths[:2]
	}
	var user *url.Userinfo
	if config.UserName != "" {
		user = url.User(config.UserName)
	}
	return &url.URL{
		User:       user,
		Host:       config.Host,
		Path:       strings.Join(paths, "/"),
		Scheme:     Scheme,
		ForceQuery: false,
	}
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *LegacyConfig) SetURL(serviceURL *url.URL) error {

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

func (service *Service) GetLegacyConfig() types.ServiceConfig {
	return &LegacyConfig{}
}
