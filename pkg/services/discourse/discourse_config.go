package discourse

import (
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
	"net/url"
	"strings"
)

// Config is the configuration needed to send Discourse notifications
type Config struct {
	APIKey   string `url:"pass"`
	Username string `url:"user"`
	Title    string `key:"title"      default:""`
	Host     string `url:"host" required:""`

	Type       postType `url:"path1" default:"regular"`
	Topic      int      `key:"topic" optional:"" default:"0"`
	Recipients string   `key:"recipients"`
	EmbedURL   string   `key:"embed-url"`
	Category   int      `key:"category" optional:"" default:"0"`
}

// GetURL returns a URL representation of its current field values
func (config *Config) GetURL() *url.URL {
	resolver := format.NewPropKeyResolver(config)
	return config.getURL(&resolver)
}

// SetURL updates a ServiceConfig from a URL representation of its field values
func (config *Config) SetURL(url *url.URL) error {
	resolver := format.NewPropKeyResolver(config)
	return config.setURL(&resolver, url)
}

func (config *Config) getURL(resolver types.ConfigQueryResolver) (u *url.URL) {
	u = &url.URL{
		User:     url.UserPassword(config.Username, config.APIKey),
		Host:     config.Host,
		Scheme:   Scheme,
		RawQuery: format.BuildQuery(resolver),
	}

	u.Path = "/" + config.Type.String()

	return u
}

// SetURL updates a ServiceConfig from a URL representation of its field values
func (config *Config) setURL(resolver types.ConfigQueryResolver, url *url.URL) error {

	config.Host = url.Host
	config.APIKey, _ = url.User.Password()
	config.Username = url.User.Username()

	paths := strings.Split(url.Path, "/")

	if len(paths) > 1 {
		config.Type = postType(PostTypes.Enum.Parse(paths[1]))
	}

	for key, vals := range url.Query() {
		if err := resolver.Set(key, vals[0]); err != nil {
			return err
		}
	}

	return nil
}

// Enums returns the fields that should use a corresponding EnumFormatter to Print/Parse their values
func (config Config) Enums() map[string]types.EnumFormatter {
	return map[string]types.EnumFormatter{
		"Type": PostTypes.Enum,
	}
}

// Scheme is the identifying part of this service's configuration URL
const Scheme = "discourse"
