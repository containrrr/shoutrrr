//go:generate go run ../../../cmd/shoutrrr-gen
package pushbullet

import (
	"errors"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/conf"
)

const (
	//Scheme is the scheme part of the service configuration URL
	Scheme = "pushbullet"
)

// ErrorTokenIncorrectSize is the error returned when the token size is incorrect
var ErrorTokenIncorrectSize = errors.New("token has incorrect size")

func (config *Config) UpdateLegacyURL(legacyURL *url.URL) *url.URL {
	if len(legacyURL.Fragment) > 0 {
		updatedURL := *legacyURL
		updatedURL.Fragment = ""
		paths := append(conf.SplitPath(legacyURL.Path), "#"+legacyURL.Fragment)
		updatedURL.Path = conf.JoinPath(paths...)
		return &updatedURL
	}
	return legacyURL
}
