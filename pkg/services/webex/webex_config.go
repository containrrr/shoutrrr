package webex

import (
	"errors"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Config is the configuration needed to send webex notifications
type Config struct {
	standard.EnumlessConfig
	BotToken string   `url:"user"`
	Rooms    []string `key:"rooms"`
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

func (config *Config) getURL(resolver types.ConfigQueryResolver) (u *url.URL) {
	u = &url.URL{
		User:       url.User(config.BotToken),
		Host:       Scheme,
		Scheme:     Scheme,
		RawQuery:   format.BuildQuery(resolver),
		ForceQuery: true,
	}

	return u
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *Config) setURL(resolver types.ConfigQueryResolver, url *url.URL) error {
	config.BotToken = url.User.Username()

	if len(url.Path) > 0 {
		switch url.Path {
		// todo: implement markdown and card functionality separately
		default:
			return errors.New("illegal argument in config URL")
		}
	}

	for key, vals := range url.Query() {
		if err := resolver.Set(key, vals[0]); err != nil {
			return err
		}
	}

	if len(config.Rooms) < 1 {
		return errors.New("no rooms defined in config URL")
	}

	if len(config.BotToken) < 1 {
		return errors.New("bot token missing from config URL")
	}

	for key, vals := range url.Query() {
		if err := resolver.Set(key, vals[0]); err != nil {
			return err
		}
	}

	return nil
}

// Scheme is the identifying part of this service's configuration URL
const Scheme = "webex"
