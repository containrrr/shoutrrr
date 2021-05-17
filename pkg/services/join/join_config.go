package join

import (
	"errors"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
	"net/url"
)

// Config for the Pushover notification service service
type Config struct {
	APIKey  string   `url:"pass"`
	Devices []string `key:"devices" desc:"Comma separated list of device IDs"`
	Title   string   `key:"title" optional:"" desc:"If set creates a notification"`
	Icon    string   `key:"icon" optional:"" desc:"Icon URL"`
}

// Enums returns the fields that should use a corresponding EnumFormatter to Print/Parse their values
func (config *Config) Enums() map[string]types.EnumFormatter {
	return map[string]types.EnumFormatter{}
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
		User:       url.UserPassword("Token", config.APIKey),
		Host:       "join",
		Scheme:     Scheme,
		ForceQuery: true,
		RawQuery:   format.BuildQuery(resolver),
	}
}

func (config *Config) setURL(resolver types.ConfigQueryResolver, url *url.URL) error {
	password, _ := url.User.Password()

	config.APIKey = password

	for key, vals := range url.Query() {
		if err := resolver.Set(key, vals[0]); err != nil {
			return err
		}
	}

	if len(config.Devices) < 1 {
		return errors.New(string(DevicesMissing))
	}

	if len(config.APIKey) < 1 {
		return errors.New(string(APIKeyMissing))
	}

	return nil
}

// Scheme is the identifying part of this service's configuration URL
const Scheme = "join"
