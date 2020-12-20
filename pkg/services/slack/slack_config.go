package slack

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"net/url"
	"strings"
)

// Config for the slack service
type Config struct {
	standard.EnumlessConfig
	BotName string   `default:"Shoutrrr"`
	Token   []string `description:"List of comma separated token parts"`
}

// GetURL returns a URL representation of it's current field values
func (config *Config) GetURL() *url.URL {
	return &url.URL{
		User:       url.User(config.BotName),
		Host:       config.Token[0],
		Path:       fmt.Sprintf("/%s/%s", config.Token[1], config.Token[2]),
		Scheme:     Scheme,
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

	token := strings.Split(serviceURL.Path, "/")
	token[0] = host

	if len(token) < 2 {
		token = []string{"", "", ""}
	}

	config.BotName = botName
	config.Token = token

	if err := ValidateToken(config.Token); err != nil {
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
