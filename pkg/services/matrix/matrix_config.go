//go:generate go run ../../../cmd/shoutrrr-gen
package matrix

import (
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/conf"
	"github.com/containrrr/shoutrrr/pkg/pkr"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	t "github.com/containrrr/shoutrrr/pkg/types"
)

// Config is the configuration for the matrix service
type LegacyConfig struct {
	standard.EnumlessConfig

	User       string   `optional:"" url:"user" desc:"Username or empty when using access token"`
	Password   string   `url:"password" desc:"Password or access token"`
	DisableTLS bool     `key:"disableTLS" default:"No"`
	Host       string   `url:"host"`
	Rooms      []string `key:"rooms,room" optional:"" desc:"Room aliases, or with ! prefix, room IDs"`
	Title      string   `key:"title" default:""`
}

// GetURL returns a URL representation of it's current field values
func (c *LegacyConfig) GetURL() *url.URL {
	resolver := pkr.NewPropKeyResolver(c)
	return c.getURL(&resolver)
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (c *LegacyConfig) SetURL(url *url.URL) error {
	resolver := pkr.NewPropKeyResolver(c)
	return c.setURL(&resolver, url)
}

func (c *LegacyConfig) getURL(resolver t.ConfigQueryResolver) *url.URL {
	return &url.URL{
		User:       url.UserPassword(c.User, c.Password),
		Host:       c.Host,
		Scheme:     Scheme,
		ForceQuery: true,
		RawQuery:   pkr.BuildQuery(resolver),
	}

}

func (c *LegacyConfig) setURL(resolver t.ConfigQueryResolver, configURL *url.URL) error {

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

func (config *Config) UpdateLegacyURL(legacyURL *url.URL) *url.URL {
	updatedURL := *legacyURL
	query := legacyURL.Query()

	for _, key := range []string{"rooms", "room"} {
		rooms, _ := conf.ParseListValue(query.Get(key), ",")
		if len(rooms) < 1 {
			continue
		}
		for r, room := range rooms {
			// If room does not begin with a '#' let's prepend it
			if room[0] != '#' && room[0] != '!' {
				rooms[r] = "#" + room
			}
		}

		query.Set(key, conf.FormatListValue(rooms, ","))
	}

	updatedURL.RawQuery = query.Encode()
	return &updatedURL
}
