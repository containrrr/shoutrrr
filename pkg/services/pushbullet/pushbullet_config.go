package pushbullet

import (
	"errors"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
	"net/url"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
)

// Config ...
type Config struct {
	standard.EnumlessConfig
	Targets []string `url:"path"`
	Token   string   `url:"host"`
	Title   string   `key:"title" default:"Shoutrrr notification"`
}

// GetURL returns a URL representation of it's current field values
func (config *Config) GetURL() *url.URL {
	resolver := format.NewPropKeyResolver(config)
	return config.getURL(&resolver)
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *Config) SetURL(url *url.URL) error {
	resolver := format.NewPropKeyResolver(config)
	return config.setURL(&resolver, url)
}

func (config *Config) getURL(resolver types.ConfigQueryResolver) *url.URL {
	return &url.URL{
		Host:       config.Token,
		Path:       "/" + strings.Join(config.Targets, "/"),
		Scheme:     Scheme,
		ForceQuery: false,
		RawQuery:   format.BuildQuery(resolver),
	}
}

func (config *Config) setURL(resolver types.ConfigQueryResolver, url *url.URL) error {
	path := url.Path

	if len(path) > 0 && path[0] == '/' {
		// Remove initial slash to skip empty first target
		path = path[1:]
	}

	if url.Fragment != "" {
		path += fmt.Sprintf("/#%s", url.Fragment)
	}

	targets := strings.Split(path, "/")

	token := url.Hostname()
	if err := validateToken(token); err != nil {
		return err
	}

	config.Token = token
	config.Targets = targets

	for key, vals := range url.Query() {
		if err := resolver.Set(key, vals[0]); err != nil {
			return err
		}
	}

	return nil
}

func validateToken(token string) error {
	if len(token) != 34 {
		return ErrorTokenIncorrectSize
	}
	return nil
}

const (
	//Scheme is the scheme part of the service configuration URL
	Scheme = "pushbullet"
)

// ErrorTokenIncorrectSize is the error returned when the token size is incorrect
var ErrorTokenIncorrectSize = errors.New("token has incorrect size")
