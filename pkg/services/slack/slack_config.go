package slack

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"net/url"
	"strings"
)

// Config for the slack service
type Config struct {
	standard.QuerylessConfig
	standard.EnumlessConfig
	BotName string
	Token   Token
}

// GetURL returns a URL representation of it's current field values
func (config *Config) GetURL() *url.URL {
	return &url.URL{
		User: url.UserPassword(config.BotName, config.Token.String()),
		Host: config.Token.A,
		Path: fmt.Sprintf("/%s/%s", config.Token.B, config.Token.C),
		Scheme: Scheme,
		ForceQuery: false,
	}
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *Config) SetURL(serviceURL *url.URL) error {

	botName := serviceURL.User.Username()
	if botName == "" {
		botName = DefaultUser
	}

	host := serviceURL.Hostname()

	path := strings.Split(serviceURL.Path, "/")

	if len(path) <2 {
		path = []string { "", "", "" }
	}

	config.BotName = botName
	config.Token = Token{
		A: host,
		B: path[1],
		C: path[2],
	}

	if err := validateToken(config.Token); err != nil {
		return err
	}

	return nil
}

const (
	// DefaultUser for sending notifications to slack
	DefaultUser = "Shoutrrr"
	// Scheme is the identifying part of this service's configuration URL
	Scheme = "slack"
)

// CreateConfigFromURL to use within the slack service
func CreateConfigFromURL(serviceURL *url.URL) (*Config, error) {
	config := Config{}
	err := config.SetURL(serviceURL)
	return &config, err
}