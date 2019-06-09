package slack

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/types"
	"net/url"
	"strings"
)

// Config for the slack service
type Config struct {
	BotName string
	Token   Token
}

// QueryFields returns the fields that are part of the Query of the service URL
func (config *Config) QueryFields() []string {
	return []string{}
}

// Enums returns the fields that should use a corresponding EnumFormatter to Print/Parse their values
func (config *Config) Enums() map[string]types.EnumFormatter {
	return map[string]types.EnumFormatter{}
}

// Get returns the value of a Query field
func (config *Config) Get(string) (string, error) {
	return "", nil
}

// Set updates the value of a Query field
func (config *Config) Set(string, string) error {
	return nil
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