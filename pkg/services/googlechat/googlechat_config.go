package googlechat

import (
	"errors"
	"net/url"

	"github.com/nicholas-fedor/shoutrrr/pkg/format"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/standard"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
)

type Config struct {
	standard.EnumlessConfig
	Host  string `default:"chat.googleapis.com"`
	Path  string
	Token string
	Key   string
}

func (config *Config) GetURL() *url.URL {
	resolver := format.NewPropKeyResolver(config)

	return config.getURL(&resolver)
}

func (config *Config) SetURL(url *url.URL) error {
	resolver := format.NewPropKeyResolver(config)

	return config.setURL(&resolver, url)
}

func (config *Config) setURL(_ types.ConfigQueryResolver, serviceURL *url.URL) error {
	config.Host = serviceURL.Host
	config.Path = serviceURL.Path

	query := serviceURL.Query()
	config.Key = query.Get("key")
	config.Token = query.Get("token")

	// Only enforce if explicitly provided but empty
	if query.Has("key") && config.Key == "" {
		return errors.New("missing field 'key'")
	}

	if query.Has("token") && config.Token == "" {
		return errors.New("missing field 'token'")
	}

	return nil
}

func (config *Config) getURL(_ types.ConfigQueryResolver) *url.URL {
	query := url.Values{}
	query.Set("key", config.Key)
	query.Set("token", config.Token)

	return &url.URL{
		Host:     config.Host,
		Path:     config.Path,
		RawQuery: query.Encode(),
		Scheme:   Scheme,
	}
}

const (
	Scheme = "googlechat"
)
