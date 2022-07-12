//go:generate go run ../../../cmd/shoutrrr-gen --lang go
package slack

import (
	"net/url"
	"strings"
)

const (
	// Scheme is the identifying part of this service's configuration URL
	Scheme = "slack"
)

// CreateConfigFromURL to use within the slack service
func CreateConfigFromURL(serviceURL *url.URL) (*Config, error) {
	config := Config{}
	err := config.SetURL(serviceURL)
	return &config, err
}

func (config *Config) setToken(value string) (*Token, error) {
	return ParseToken(value)
}

func (config *Config) getToken() string {
	value, _ := config.Token.GetPropValue()
	return value
}

func (config *Config) emptyToken(value *Token) bool {
	return value.IsEmpty()
}

func (config *Config) UpdateLegacyURL(serviceURL *url.URL) *url.URL {

	if len(serviceURL.Path) > 1 {
		// Reading legacy config URL format
		updatedURL := *serviceURL
		token := strings.Replace(serviceURL.Hostname()+serviceURL.Path, "/", "-", -1)
		updatedURL.User = url.UserPassword(hookTokenIdentifier, token)
		updatedURL.Path = ""
		updatedURL.Host = "webhook"
		query := updatedURL.Query()
		query.Add("botname", serviceURL.User.Username())
		updatedURL.RawQuery = query.Encode()
		return &updatedURL
	}

	return serviceURL
}
