//go:generate go run ../../../cmd/shoutrrr-gen
package mattermost

import "net/url"

//ErrorMessage for error events within the mattermost service
type ErrorMessage string

const (
	// Scheme is the identifying part of this service's configuration URL
	Scheme = "mattermost"
	// NotEnoughArguments provided in the service URL
	NotEnoughArguments ErrorMessage = "the apiURL does not include enough arguments, either provide 1 or 3 arguments (they may be empty)"
)

// CreateConfigFromURL to use within the mattermost service
func CreateConfigFromURL(serviceURL *url.URL) (*Config, error) {
	config := Config{}
	err := config.SetURL(serviceURL)
	return &config, err
}
