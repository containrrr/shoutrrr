//go:generate go run ../../../cmd/shoutrrr-gen --lang go
package smtp

import (
	"net/url"
)

// Scheme is the identifying part of this service's configuration URL
const Scheme = "smtp"

func (config *Config) Hostname() string {
	return (&url.URL{Host: config.Host}).Hostname()
}
