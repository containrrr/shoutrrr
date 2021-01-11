package matrix

import (
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	t "github.com/containrrr/shoutrrr/pkg/types"
	"net/url"
)

type Config struct {
	standard.EnumlessConfig

	User       string `optional:""`
	Password   string
	DisableTLS bool `key:"disableTLS" default:"No"`
	Host       string
	Rooms      []string `key:"rooms" optional:""`
}

// GetURL returns a URL representation of it's current field values
func (c *Config) GetURL() *url.URL {
	resolver := format.NewPropKeyResolver(c)
	return c.getURL(&resolver)
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (c *Config) SetURL(url *url.URL) error {
	resolver := format.NewPropKeyResolver(c)
	return c.setURL(&resolver, url)
}

func (c *Config) getURL(resolver t.ConfigQueryResolver) *url.URL {
	return &url.URL{
		User:       url.UserPassword(c.User, c.Password),
		Host:       c.Host,
		Scheme:     Scheme,
		ForceQuery: true,
		RawQuery:   format.BuildQuery(resolver),
	}

}

func (c *Config) setURL(resolver t.ConfigQueryResolver, configURL *url.URL) error {

	c.User = configURL.User.Username()
	password, _ := configURL.User.Password()
	c.Password = password
	c.Host = configURL.Host

	for key, vals := range configURL.Query() {
		if err := resolver.Set(key, vals[0]); err != nil {
			return err
		}
	}

	return nil
}
