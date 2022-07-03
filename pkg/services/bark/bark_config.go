//go:generate go run ../../../cmd/shoutrrr-gen --lang go ../../../spec/bark.yml

package bark

import (
	"net/url"
	"strings"
)

// GetAPIURL returns the API URL corresponding to the passed endpoint based on the configuration
func (config *Config) GetAPIURL(endpoint string) string {

	path := strings.Builder{}
	if !strings.HasPrefix(config.Path, "/") {
		path.WriteByte('/')
	}
	_, _ = path.WriteString(config.Path)
	if !strings.HasSuffix(path.String(), "/") {
		path.WriteByte('/')
	}
	path.WriteString(endpoint)

	apiURL := url.URL{
		Scheme: config.Scheme,
		Host:   config.Host,
		Path:   path.String(),
	}
	return apiURL.String()
}

// Scheme is the identifying part of this service's configuration URL
const (
	Scheme = "bark"
)
