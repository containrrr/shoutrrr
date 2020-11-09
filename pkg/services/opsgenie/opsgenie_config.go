package opsgenie

import (
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
)

// Config for use within the opsgenie plugin
type Config struct {
	ApiKey string `desc:"The OpsGenie API key"`
	Host   string `desc:"The OpsGenie API host. Use 'api.opsgenie.com' for US and 'api.eu.opsgenie.com' for EU instances"`
	standard.QuerylessConfig
	standard.EnumlessConfig
}

// GetURL returns a URL representation of it's current field values
func (config *Config) GetURL() *url.URL {
	return &url.URL{
		Host:       config.Host,
		Scheme:     Scheme,
		ForceQuery: false,
		Path:       config.ApiKey,
	}
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *Config) SetURL(url *url.URL) error {
	config.Host = url.Hostname() + ":" + url.Port()
	config.ApiKey = url.Path[1:]
	return nil
}

const (
	// Scheme is the identifying part of this service's configuration URL
	Scheme = "opsgenie"
)
