//go:generate go run ../../../cmd/shoutrrr-gen
package zulip

import (
	"errors"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/pkr"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Config for the zulip service
type LegacyConfig struct {
	standard.EnumlessConfig
	BotMail string `url:"user" desc:"Bot e-mail address"`
	BotKey  string `url:"pass" desc:"API Key"`
	Host    string `url:"host,port" desc:"API server hostname"`
	Stream  string `key:"stream" optional:"" description:"Target stream name"`
	Topic   string `key:"topic,title" default:""`
}

// GetURL returns a URL representation of it's current field values
func (config *LegacyConfig) GetURL() *url.URL {
	resolver := pkr.NewPropKeyResolver(config)
	return config.getURL(&resolver)
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *LegacyConfig) SetURL(url *url.URL) error {
	resolver := pkr.NewPropKeyResolver(config)
	return config.setURL(&resolver, url)
}

func (config *LegacyConfig) getURL(_ types.ConfigQueryResolver) *url.URL {
	query := &url.Values{}

	if config.Stream != "" {
		query.Set("stream", config.Stream)
	}

	if config.Topic != "" {
		query.Set("topic", config.Topic)
	}

	return &url.URL{
		User:     url.UserPassword(config.BotMail, config.BotKey),
		Host:     config.Host,
		RawQuery: query.Encode(),
		Scheme:   Scheme,
	}
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *LegacyConfig) setURL(_ types.ConfigQueryResolver, serviceURL *url.URL) error {
	var ok bool

	config.BotMail = serviceURL.User.Username()

	if config.BotMail == "" {
		return errors.New(string(MissingBotMail))
	}

	config.BotKey, ok = serviceURL.User.Password()

	if !ok {
		return errors.New(string(MissingAPIKey))
	}

	config.Host = serviceURL.Hostname()

	if config.Host == "" {
		return errors.New(string(MissingHost))
	}

	config.Stream = serviceURL.Query().Get("stream")
	config.Topic = serviceURL.Query().Get("topic")

	return nil
}

// Clone the config to a new Config struct
func (config *LegacyConfig) Clone() *LegacyConfig {
	return &LegacyConfig{
		BotMail: config.BotMail,
		BotKey:  config.BotKey,
		Host:    config.Host,
		Stream:  config.Stream,
		Topic:   config.Topic,
	}
}

const (
	// Scheme is the identifying part of this service's configuration URL
	Scheme = "zulip"
)

// CreateConfigFromURL to use within the zulip service
func CreateConfigFromURL(serviceURL *url.URL) (*Config, error) {
	config := Config{}
	err := config.SetURL(serviceURL)

	return &config, err
}

func (service *Service) GetLegacyConfig() types.ServiceConfig {
	return &LegacyConfig{}
}
