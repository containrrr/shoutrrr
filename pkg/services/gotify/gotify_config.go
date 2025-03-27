package gotify

import (
	"net/url"
	"strings"

	"github.com/nicholas-fedor/shoutrrr/pkg/format"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"

	"github.com/nicholas-fedor/shoutrrr/pkg/services/standard"
)

// Config for use within the gotify plugin.
type Config struct {
	standard.EnumlessConfig
	Token      string `desc:"Application token"                     required:""      url:"path2"`
	Host       string `desc:"Server hostname (and optionally port)" required:""      url:"host,port"`
	Path       string `desc:"Server subpath"                        optional:""      url:"path1"`
	Priority   int    `default:"0"                                  key:"priority"`
	Title      string `default:"Shoutrrr notification"              key:"title"`
	DisableTLS bool   `default:"No"                                 key:"disabletls"`
}

// GetURL returns a URL representation of it's current field values.
func (config *Config) GetURL() *url.URL {
	resolver := format.NewPropKeyResolver(config)

	return config.getURL(&resolver)
}

// SetURL updates a ServiceConfig from a URL representation of it's field values.
func (config *Config) SetURL(url *url.URL) error {
	resolver := format.NewPropKeyResolver(config)

	return config.setURL(&resolver, url)
}

func (config *Config) getURL(resolver types.ConfigQueryResolver) *url.URL {
	return &url.URL{
		Host:       config.Host,
		Scheme:     Scheme,
		ForceQuery: false,
		Path:       config.Path + config.Token,
		RawQuery:   format.BuildQuery(resolver),
	}
}

func (config *Config) setURL(resolver types.ConfigQueryResolver, url *url.URL) error {
	path := url.Path
	if len(path) > 0 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	tokenIndex := strings.LastIndex(path, "/") + 1

	config.Path = path[:tokenIndex]
	if config.Path == "/" {
		config.Path = config.Path[1:]
	}

	config.Host = url.Host
	config.Token = path[tokenIndex:]

	for key, vals := range url.Query() {
		if err := resolver.Set(key, vals[0]); err != nil {
			return err
		}
	}

	return nil
}

const (
	// Scheme is the identifying part of this service's configuration URL.
	Scheme = "gotify"
)
