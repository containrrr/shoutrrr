package slack

import (
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Config for the slack service
type Config struct {
	standard.EnumlessConfig
	BotName string `optional:"uses bot default" key:"botname,username" desc:"Bot name"`
	Icon    string `key:"icon,icon_emoji,icon_url" default:"" optional:"" desc:"Use emoji or URL as icon (based on presence of http(s):// prefix)"`
	Token   Token  `desc:"API Bot token" url:"user,pass"`
	Color   string `key:"color" optional:"default border color" desc:"Message left-hand border color"`
	Title   string `key:"title" optional:"omitted" desc:"Prepended text above the message"`
	Channel string `url:"host" desc:"Channel to send messages to in Cxxxxxxxxxx format"`
	ThreadTS string   `key:"thread_ts" optional:"" desc:"ts value of the parent message (to send message as reply in thread)"`
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
		User:       config.Token.UserInfo(),
		Host:       config.Channel,
		Scheme:     Scheme,
		ForceQuery: false,
		RawQuery:   format.BuildQuery(resolver),
	}
}

func (config *Config) setURL(resolver types.ConfigQueryResolver, serviceURL *url.URL) error {

	var token string
	var err error

	if len(serviceURL.Path) > 1 {
		// Reading legacy config URL format
		token = serviceURL.Hostname() + serviceURL.Path

		config.Channel = "webhook"
		config.BotName = serviceURL.User.Username()
	} else {
		token = serviceURL.User.String()
		config.Channel = serviceURL.Hostname()
	}

	if err = config.Token.SetFromProp(token); err != nil {
		return err
	}

	for key, vals := range serviceURL.Query() {
		if err := resolver.Set(key, vals[0]); err != nil {
			return err
		}
	}

	return nil
}

const (
	// Scheme is the identifying part of this service's configuration URL
	Scheme = "slack"
)

// CreateConfigFromURL to use within the slack service
func CreateConfigFromURL(serviceURL *url.URL) (*Config, error) {
	config := Config{}
	err := config.SetURL(serviceURL)
	return &config, err
}
