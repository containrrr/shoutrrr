package generic

import (
	"net/url"

	"github.com/nicholas-fedor/shoutrrr/pkg/format"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/standard"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
)

// Config for use within the generic service.
type Config struct {
	standard.EnumlessConfig
	webhookURL    *url.URL
	headers       map[string]string
	extraData     map[string]string
	ContentType   string `default:"application/json"                                desc:"The value of the Content-Type header"            key:"contenttype"`
	DisableTLS    bool   `default:"No"                                              key:"disabletls"`
	Template      string `desc:"The template used for creating the request payload" key:"template"                                         optional:""`
	Title         string `default:""                                                key:"title"`
	TitleKey      string `default:"title"                                           desc:"The key that will be used for the title value"   key:"titlekey"`
	MessageKey    string `default:"message"                                         desc:"The key that will be used for the message value" key:"messagekey"`
	RequestMethod string `default:"POST"                                            key:"method"`
}

// DefaultConfig creates a PropKeyResolver and uses it to populate the default values of a new Config, returning both.
func DefaultConfig() (*Config, format.PropKeyResolver) {
	config := &Config{}
	pkr := format.NewPropKeyResolver(config)
	_ = pkr.SetDefaultProps(config)

	return config, pkr
}

// ConfigFromWebhookURL creates a new Config from a parsed Webhook URL.
func ConfigFromWebhookURL(webhookURL url.URL) (*Config, format.PropKeyResolver, error) {
	config, pkr := DefaultConfig()

	// Process webhook URL query parameters, preserving custom params
	webhookQuery := webhookURL.Query()
	// First strip custom headers and extra data
	headers, extraData := stripCustomQueryValues(webhookQuery)
	// Then process remaining query parameters as potential escaped config params
	escapedQuery := url.Values{}

	for key, values := range webhookQuery {
		if len(values) > 0 {
			escapedQuery.Set(format.EscapeKey(key), values[0])
		}
	}
	// Apply escaped config parameters
	_, err := format.SetConfigPropsFromQuery(&pkr, escapedQuery)
	if err != nil {
		return nil, pkr, err
	}

	// Restore original query parameters for the webhook URL
	webhookURL.RawQuery = webhookQuery.Encode()
	config.webhookURL = &webhookURL
	config.headers = headers
	config.extraData = extraData
	config.DisableTLS = webhookURL.Scheme == "http"

	return config, pkr, nil
}

// WebhookURL returns a url.URL that is synchronized with the config props.
func (config *Config) WebhookURL() *url.URL {
	webhookURL := *config.webhookURL
	webhookURL.Scheme = DefaultWebhookScheme

	if config.DisableTLS {
		webhookURL.Scheme = webhookURL.Scheme[:4]
	}

	return &webhookURL
}

// GetURL returns a URL representation of its current field values.
func (config *Config) GetURL() *url.URL {
	resolver := format.NewPropKeyResolver(config)

	return config.getURL(&resolver)
}

// SetURL updates a ServiceConfig from a URL representation of its field values.
func (config *Config) SetURL(serviceURL *url.URL) error {
	resolver := format.NewPropKeyResolver(config)

	return config.setURL(&resolver, serviceURL)
}

func (config *Config) getURL(resolver types.ConfigQueryResolver) *url.URL {
	serviceURL := *config.webhookURL
	webhookQuery := config.webhookURL.Query()
	serviceQuery := format.BuildQueryWithCustomFields(resolver, webhookQuery)
	appendCustomQueryValues(serviceQuery, config.headers, config.extraData)
	serviceURL.RawQuery = serviceQuery.Encode()
	serviceURL.Scheme = Scheme

	return &serviceURL
}

func (config *Config) setURL(resolver types.ConfigQueryResolver, serviceURL *url.URL) error {
	webhookURL := *serviceURL
	serviceQuery := serviceURL.Query()
	headers, extraData := stripCustomQueryValues(serviceQuery)

	customQuery, err := format.SetConfigPropsFromQuery(resolver, serviceQuery)
	if err != nil {
		return err
	}

	webhookURL.RawQuery = customQuery.Encode()
	config.webhookURL = &webhookURL
	config.headers = headers
	config.extraData = extraData

	return nil
}

const (
	// Scheme is the identifying part of this service's configuration URL.
	Scheme = "generic"
	// DefaultWebhookScheme is the scheme used for webhook URLs unless overridden.
	DefaultWebhookScheme = "https"
)
