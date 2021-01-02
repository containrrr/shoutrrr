package zulip

import (
	"errors"
	"github.com/containrrr/shoutrrr/pkg/types"
	"net/url"
)

// Config for the zulip service
type Config struct {
	BotMail string
	BotKey  string
	Host    string
	Path    string
	Stream  string `key:"stream"`
	Topic   string `key:"topic"`
}

// GetURL returns a URL representation of it's current field values
func (config *Config) GetURL(_ types.ConfigQueryResolver) *url.URL {
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
		Path:     config.Path,
		RawQuery: query.Encode(),
		Scheme:   Scheme,
	}
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *Config) SetURL(_ types.ConfigQueryResolver, serviceURL *url.URL) error {
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

	config.Path = "api/v1/messages"
	config.Stream = serviceURL.Query().Get("stream")
	config.Topic = serviceURL.Query().Get("topic")

	return nil
}

// Clone the config to a new Config struct
func (config *Config) Clone() *Config {
	return &Config{
		BotMail: config.BotMail,
		BotKey:  config.BotKey,
		Host:    config.Host,
		Path:    config.Path,
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
	err := config.SetURL(nil, serviceURL)

	return &config, err
}
