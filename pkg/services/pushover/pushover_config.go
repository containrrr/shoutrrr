//go:generate go run ../../../cmd/shoutrrr-gen --lang go
package pushover

import (
	"errors"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/pkr"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Config for the Pushover notification service service
type LegacyConfig struct {
	Token    string   `url:"pass" desc:"API Token/Key"`
	User     string   `url:"host" desc:"User Key"`
	Devices  []string `key:"devices" optional:""`
	Priority int8     `key:"priority" default:"0"`
	Title    string   `key:"title" optional:""`
}

// Enums returns the fields that should use a corresponding EnumFormatter to Print/Parse their values
func (config *LegacyConfig) Enums() map[string]types.EnumFormatter {
	return map[string]types.EnumFormatter{}
}

// GetURL returns a URL representation of it's current field values
func (config *LegacyConfig) GetURL() *url.URL {
	resolver := pkr.NewPropKeyResolver(config)
	return &url.URL{
		User:       url.UserPassword("Token", config.Token),
		Host:       config.User,
		Scheme:     Scheme,
		ForceQuery: true,
		RawQuery:   pkr.BuildQuery(&resolver),
	}
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *LegacyConfig) SetURL(url *url.URL) error {
	resolver := pkr.NewPropKeyResolver(config)
	return config.setURL(&resolver, url)
}

func (config *LegacyConfig) setURL(resolver types.ConfigQueryResolver, url *url.URL) error {
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

func (service *Service) GetLegacyConfig() types.ServiceConfig {
	return &LegacyConfig{}
}
