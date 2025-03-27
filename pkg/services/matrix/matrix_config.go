package matrix

import (
	"net/url"

	"github.com/nicholas-fedor/shoutrrr/pkg/format"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/standard"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
)

// Config is the configuration for the matrix service.
type Config struct {
	standard.EnumlessConfig

	User       string   `desc:"Username or empty when using access token" optional:""      url:"user"`
	Password   string   `desc:"Password or access token"                  url:"password"`
	DisableTLS bool     `default:"No"                                     key:"disableTLS"`
	Host       string   `url:"host"`
	Rooms      []string `desc:"Room aliases, or with ! prefix, room IDs"  key:"rooms,room" optional:""`
	Title      string   `default:""                                       key:"title"`
}

// GetURL returns a URL representation of it's current field values.
func (c *Config) GetURL() *url.URL {
	resolver := format.NewPropKeyResolver(c)

	return c.getURL(&resolver)
}

// SetURL updates a ServiceConfig from a URL representation of it's field values.
func (c *Config) SetURL(url *url.URL) error {
	resolver := format.NewPropKeyResolver(c)

	return c.setURL(&resolver, url)
}

func (c *Config) getURL(resolver types.ConfigQueryResolver) *url.URL {
	return &url.URL{
		User:       url.UserPassword(c.User, c.Password),
		Host:       c.Host,
		Scheme:     Scheme,
		ForceQuery: true,
		RawQuery:   format.BuildQuery(resolver),
	}
}

func (c *Config) setURL(resolver types.ConfigQueryResolver, configURL *url.URL) error {
	c.User = configURL.User.Username()
	password, _ := configURL.User.Password()
	c.Password = password
	c.Host = configURL.Host

	for key, vals := range configURL.Query() {
		if err := resolver.Set(key, vals[0]); err != nil {
			return err
		}
	}

	for r, room := range c.Rooms {
		// If room does not begin with a '#' let's prepend it
		if room[0] != '#' && room[0] != '!' {
			c.Rooms[r] = "#" + room
		}
	}

	return nil
}
