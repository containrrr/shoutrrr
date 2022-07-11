//go:generate go run ../../../cmd/shoutrrr-gen
package generic

import (
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/conf"
	"github.com/containrrr/shoutrrr/pkg/pkr"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	t "github.com/containrrr/shoutrrr/pkg/types"
)

// Config for use within the generic service
type LegacyConfig struct {
	standard.EnumlessConfig
	webhookURL  *url.URL
	ContentType string `key:"contenttype" default:"application/json" desc:"The value of the Content-Type header"`
	DisableTLS  bool   `key:"disabletls" default:"No"`
	Template    string `key:"template" optional:""`
	Title       string `key:"title" default:""`
}

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

	// config.webhookURL = &webhookURL
	// TODO: Decide what to do with custom URL queries. Right now they are passed
	//       to the inner url.URL and not processed by PKR.
	// customQuery, err := format.SetConfigPropsFromQuery(&pkr, webhookURL.Query())
	// goland:noinspection GoNilness: SetConfigPropsFromQuery always return non-nil
	// config.webhookURL.RawQuery = customQuery.Encode()
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

// GetURL returns a URL representation of it's current field values
func (config *LegacyConfig) GetURL() *url.URL {
	resolver := pkr.NewPropKeyResolver(config)
	return config.getURL(&resolver)
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *LegacyConfig) SetURL(serviceURL *url.URL) error {
	resolver := pkr.NewPropKeyResolver(config)
	return config.setURL(&resolver, serviceURL)
}

func (config *LegacyConfig) getURL(resolver t.ConfigQueryResolver) *url.URL {

	serviceURL := *config.webhookURL
	webhookQuery := config.webhookURL.Query()
	serviceQuery := pkr.BuildQueryWithCustomFields(resolver, webhookQuery)
	serviceURL.RawQuery = serviceQuery.Encode()
	serviceURL.Scheme = Scheme

	return &serviceURL
}

func (config *LegacyConfig) setURL(resolver t.ConfigQueryResolver, serviceURL *url.URL) error {

	webhookURL := *serviceURL
	serviceQuery := serviceURL.Query()

	customQuery, err := pkr.SetConfigPropsFromQuery(resolver, serviceQuery)
	if err != nil {
		return err
	}
	webhookURL.RawQuery = customQuery.Encode()
	config.webhookURL = &webhookURL

	return nil
}

const (
	// Scheme is the identifying part of this service's configuration URL
	Scheme = "generic"
	// DefaultWebhookScheme is the scheme used for webhook URLs unless overridden
	DefaultWebhookScheme = "https"
)
