//go:generate go run ../../../cmd/shoutrrr-gen
package generic

import (
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/conf"
)

// DefaultConfig creates a PropKeyResolver and uses it to populate the default values of a new Config, returning both
func DefaultConfig() *Config {
	config := &Config{}
	conf.SetDefaults(config)
	return config
}

// ConfigFromWebhookURL creates a new Config from a parsed Webhook URL
func ConfigFromWebhookURL(webhookURL url.URL) (*Config, error) {
	config := DefaultConfig()

	query := url.Values{}
	for key, value := range webhookURL.Query() {
		escaped := conf.EscapeCustomQueryKey(config, key)
		query.Set(escaped, value[0])
	}
	webhookURL.RawQuery = query.Encode()

	if err := config.SetURL(&webhookURL); err != nil {
		return nil, err
	}

	// TODO: Decide what to do with custom URL queries. Right now they are passed
	//       to the inner url.URL and not processed by PKR.

	config.DisableTLS = webhookURL.Scheme == "http"
	return config, nil
}

// WebhookURL returns a url.URL that is synchronized with the config props
func (config *Config) WebhookURL() *url.URL {
	webhookURL := *config.GetURL()
	// webhookURL
	query := url.Values{}
	for key, value := range config.Query {
		escaped := conf.UnescapeCustomQueryKey(key)
		query.Set(escaped, value[0])
	}
	webhookURL.RawQuery = query.Encode()
	webhookURL.Scheme = DefaultWebhookScheme
	if config.DisableTLS {
		webhookURL.Scheme = webhookURL.Scheme[:4]
	}
	return &webhookURL
}

const (
	// Scheme is the identifying part of this service's configuration URL
	Scheme = "generic"
	// DefaultWebhookScheme is the scheme used for webhook URLs unless overridden
	DefaultWebhookScheme = "https"
)
