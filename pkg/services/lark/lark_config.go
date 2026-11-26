package lark

import (
	"net/url"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
)

const Scheme = "lark"

type Config struct {
	Host   string `desc:"Custom bot URL Host" default:"open.larksuite.com" url:"Host"`
	Secret string `desc:"Custom bot secret" default:"" key:"secret"`
	Path   string `desc:"Custom bot token" url:"Path"`
	Title  string `desc:"Message Title" default:"" key:"title"`
}

func (config *Config) Enums() map[string]types.EnumFormatter {
	return map[string]types.EnumFormatter{}
}

func (config *Config) GetURL() *url.URL {
	resolver := format.NewPropKeyResolver(config)
	return config.getURL(&resolver)
}

func (config *Config) getURL(resolver types.ConfigQueryResolver) *url.URL {

	return &url.URL{
		Host:       config.Host,
		Path:       "/" + config.Path,
		Scheme:     Scheme,
		ForceQuery: true,
		RawQuery:   format.BuildQuery(resolver),
	}

}

func (config *Config) SetURL(url *url.URL) error {
	resolver := format.NewPropKeyResolver(config)
	return config.setURL(&resolver, url)
}

func (config *Config) setURL(resolver types.ConfigQueryResolver, url *url.URL) error {

	config.Host = url.Host
	config.Path = strings.Trim(url.Path, "/")
	// config.Password = url.Hostname()

	for key, vals := range url.Query() {
		if err := resolver.Set(key, vals[0]); err != nil {
			return err
		}
	}

	return nil
}
