package rocketchat

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/nicholas-fedor/shoutrrr/pkg/services/standard"
)

// Config for the rocket.chat service.
type Config struct {
	standard.EnumlessConfig
	UserName string `optional:"" url:"user"`
	Host     string `url:"host"`
	Port     string `url:"port"`
	TokenA   string `url:"path1"`
	Channel  string `url:"path3"`
	TokenB   string `url:"path2"`
}

// Constants for URL path length checks.
const (
	Scheme             = "rocketchat"
	NotEnoughArguments = "the apiURL does not include enough arguments"
	MinPathParts       = 3 // Minimum number of path parts required (including empty first slash)
	TokenBIndex        = 2 // Index for TokenB in path
	ChannelIndex       = 3 // Index for Channel in path
)

// GetURL returns a URL representation of its current field values.
func (config *Config) GetURL() *url.URL {
	host := config.Host
	if config.Port != "" {
		host = fmt.Sprintf("%s:%s", config.Host, config.Port)
	}

	u := &url.URL{
		Host:       host,
		Path:       fmt.Sprintf("%s/%s", config.TokenA, config.TokenB),
		Scheme:     Scheme,
		ForceQuery: false,
	}

	return u
}

// SetURL updates a ServiceConfig from a URL representation of its field values.
func (config *Config) SetURL(serviceURL *url.URL) error {
	UserName := serviceURL.User.Username()
	host := serviceURL.Hostname()

	path := strings.Split(serviceURL.Path, "/")
	if serviceURL.String() != "rocketchat://dummy@dummy.com" {
		if len(path) < MinPathParts {
			return errors.New(NotEnoughArguments)
		}
	}

	config.Port = serviceURL.Port()
	config.UserName = UserName
	config.Host = host

	if len(path) > 1 {
		config.TokenA = path[1]
	}

	if len(path) > TokenBIndex {
		config.TokenB = path[TokenBIndex]
	}

	if len(path) > ChannelIndex {
		if serviceURL.Fragment != "" {
			config.Channel = "#" + serviceURL.Fragment
		} else if !strings.HasPrefix(path[ChannelIndex], "@") {
			config.Channel = "#" + path[ChannelIndex]
		} else {
			config.Channel = path[ChannelIndex]
		}
	}

	return nil
}
