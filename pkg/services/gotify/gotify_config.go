//go:generate go run ../../../cmd/shoutrrr-gen
package gotify

import (
	"net/url"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/pkr"
	"github.com/containrrr/shoutrrr/pkg/types"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
)

// Config for use within the gotify plugin
type LegacyConfig struct {
	standard.EnumlessConfig
	Token      string `url:"path2" desc:"Application token" required:""`
	Host       string `url:"host,port" desc:"Server hostname (and optionally port)" required:""`
	Path       string `optional:"" url:"path1" desc:"Server subpath"`
	Priority   int    `key:"priority" default:"0"`
	Title      string `key:"title" default:"Shoutrrr notification"`
	DisableTLS bool   `key:"disabletls" default:"No"`
}

// GetURL returns a URL representation of it's current field values
func (config *LegacyConfig) GetURL() *url.URL {
	resolver := pkr.NewPropKeyResolver(config)
	return config.getURL(&resolver)
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *LegacyConfig) SetURL(url *url.URL) error {
	resolver := pkr.NewPropKeyResolver(config)
	return config.setURL(&resolver, url)
}

func (config *LegacyConfig) getURL(resolver types.ConfigQueryResolver) *url.URL {
	return &url.URL{
		Host:       config.Host,
		Scheme:     Scheme,
		ForceQuery: false,
		Path:       config.Path + config.Token,
		RawQuery:   pkr.BuildQuery(resolver),
	}
}

func (config *LegacyConfig) setURL(resolver types.ConfigQueryResolver, url *url.URL) error {

	tokenIndex := strings.LastIndex(url.Path, "/")
	config.Path = url.Path[:tokenIndex]

	config.Host = url.Host
	config.Token = url.Path[tokenIndex:]

	for key, vals := range url.Query() {
		if err := resolver.Set(key, vals[0]); err != nil {
			return err
		}
	}
	return nil
}

const (
	// Scheme is the identifying part of this service's configuration URL
	Scheme = "gotify"
)
