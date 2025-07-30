package pagerduty

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Config for use within the pagerduty service
type Config struct {
	IntegrationKey string `url:"path" desc:"The PagerDuty API integration key"`
	Host           string `url:"host" desc:"The PagerDuty API host." default:"events.pagerduty.com"`
	Port           uint16 `url:"port" desc:"The PagerDuty API port." default:"443"`
	Severity       string `key:"severity" desc:"The perceived severity of the status the event (critical, error, warning, or info); required" default:"error"`
	Source         string `key:"source" desc:"The unique location of the affected system, preferably a hostname or FQDN; required" default:"default"`
	Action         string `key:"action" desc:"The type of event (trigger, acknowledge, or resolve)" default:"trigger"`
}

// Enums returns an empty map because the PagerDuty service doesn't use Enums
func (config Config) Enums() map[string]types.EnumFormatter {
	return map[string]types.EnumFormatter{}
}

// GetURL is the public version of getURL that creates a new PropKeyResolver when accessed from outside the package
func (config *Config) GetURL() *url.URL {
	resolver := format.NewPropKeyResolver(config)
	return config.getURL(&resolver)
}

// Private version of GetURL that can use an instance of PropKeyResolver instead of rebuilding its model from reflection
func (config *Config) getURL(resolver types.ConfigQueryResolver) *url.URL {
	host := ""
	if config.Port > 0 {
		host = fmt.Sprintf("%s:%d", config.Host, config.Port)
	} else {
		host = config.Host
	}

	result := &url.URL{
		Host:     host,
		Path:     fmt.Sprintf("/%s", config.IntegrationKey),
		Scheme:   Scheme,
		RawQuery: format.BuildQuery(resolver),
	}

	return result
}

// SetURL updates a ServiceConfig from a URL representation of its field values
func (config *Config) SetURL(url *url.URL) error {
	resolver := format.NewPropKeyResolver(config)
	return config.setURL(&resolver, url)
}

// Private version of SetURL that can use an instance of PropKeyResolver instead of rebuilding its model from reflection
func (config *Config) setURL(resolver types.ConfigQueryResolver, url *url.URL) error {
	config.IntegrationKey = url.Path[1:]

	if url.Hostname() != "" {
		config.Host = url.Hostname()
	}

	if url.Port() != "" {
		port, err := strconv.ParseUint(url.Port(), 10, 16)
		if err != nil {
			return err
		}
		config.Port = uint16(port)
	}

	for key, vals := range url.Query() {
		if err := resolver.Set(key, vals[0]); err != nil {
			return err
		}
	}

	return nil
}

const (
	// Scheme is the identifying part of this service's configuration URL
	Scheme = "pagerduty"
)
