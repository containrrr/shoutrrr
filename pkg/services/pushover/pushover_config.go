package pushover

import (
	"errors"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
	"net/url"
)

// Config for the Pushover notification service service
type Config struct {
	Token    string   `url:"pass" desc:"API Token/Key"`
	User     string   `url:"host" desc:"User Key"`
	Devices  []string `key:"devices" optional:""`
	Priority int8     `key:"priority" default:"0"`
	Title    string   `key:"title" optional:""`
}

// Enums returns the fields that should use a corresponding EnumFormatter to Print/Parse their values
func (config *Config) Enums() map[string]types.EnumFormatter {
	return map[string]types.EnumFormatter{}
}

// GetURL returns a URL representation of it's current field values
func (config *Config) GetURL() *url.URL {
	resolver := format.NewPropKeyResolver(config)
	return &url.URL{
		User:       url.UserPassword("Token", config.Token),
		Host:       config.User,
		Scheme:     Scheme,
		ForceQuery: true,
		RawQuery:   format.BuildQuery(&resolver),
	}
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *Config) SetURL(url *url.URL) error {
	resolver := format.NewPropKeyResolver(config)
	return config.setURL(&resolver, url)
}

func (config *Config) setURL(resolver types.ConfigQueryResolver, url *url.URL) error {
	password, _ := url.User.Password()

	config.User = url.Host
	config.Token = password

	for key, vals := range url.Query() {
		if err := resolver.Set(key, vals[0]); err != nil {
			return err
		}
	}

	if len(config.User) < 1 {
		return errors.New(string(UserMissing))
	}

	if len(config.Token) < 1 {
		return errors.New(string(TokenMissing))
	}

	return nil
}

// Scheme is the identifying part of this service's configuration URL
const Scheme = "pushover"
