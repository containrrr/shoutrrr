package slack

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
	"net/url"
	"strings"
)

// Config for the slack service
type Config struct {
	standard.EnumlessConfig
	BotName string   `default:"" optional:""`
	Token   []string `description:"List of comma separated token parts"`
	Color   string   `key:"color" optional:""`
	Title   string   `key:"title" optional:""`
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
		User:       url.User(config.BotName),
		Host:       config.Token[0],
		Path:       fmt.Sprintf("/%s/%s", config.Token[1], config.Token[2]),
		Scheme:     Scheme,
		ForceQuery: false,
		RawQuery:   format.BuildQuery(resolver),
	}
}

func (config *Config) setURL(resolver types.ConfigQueryResolver, serviceURL *url.URL) error {

	botName := serviceURL.User.Username()

	host := serviceURL.Hostname()

	token := strings.Split(serviceURL.Path, "/")
	token[0] = host

	if len(token) < 2 {
		token = []string{"", "", ""}
	}

	config.BotName = botName
	config.Token = token

	if err := ValidateToken(config.Token); err != nil {
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
