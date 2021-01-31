package teams

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
	"net/url"
	"regexp"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
)

// Config for use within the teams plugin
type Config struct {
	standard.EnumlessConfig
	WebhookParts [4]string
	Title        string `key:"title" optional:""`
	Color        string `key:"color" optional:""`
	Host         string `key:"host" optional:"" default:"outlook.office.com"`
}

// SetFromWebhookURL updates the config WebhookParts from a teams webhook URL
func (config *Config) SetFromWebhookURL(webhookURL string) error {
	parts, err := parseAndVerifyWebhookURL(webhookURL)
	if err != nil {
		return err
	}

	config.WebhookParts = parts
	return nil
}

// ConfigFromWebhookURL creates a new Config from a parsed Teams Webhook URL
func ConfigFromWebhookURL(webhookURL url.URL) (*Config, error) {
	config := &Config{
		Host: webhookURL.Host,
	}

	if err := config.SetFromWebhookURL(webhookURL.String()); err != nil {
		return nil, err
	}

	return config, nil
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
	parts := config.WebhookParts

	return &url.URL{
		User:       url.User(parts[0]),
		Host:       parts[1],
		Path:       "/" + parts[2] + "/" + parts[3],
		Scheme:     Scheme,
		ForceQuery: false,
		RawQuery:   format.BuildQuery(resolver),
	}
}

func (config *Config) setURL(resolver types.ConfigQueryResolver, url *url.URL) error {
	var webhookParts [4]string

	if pass, legacyFormat := url.User.Password(); legacyFormat {
		parts := strings.Split(url.User.Username(), "@")
		if len(parts) != 2 {
			return fmt.Errorf("invalid URL format")
		}
		webhookParts = [4]string{parts[0], parts[1], pass, url.Hostname()}
	} else {
		parts := strings.Split(url.Path, "/")
		if parts[0] == "" {
			parts = parts[1:]
		}
		webhookParts = [4]string{url.User.Username(), url.Hostname(), parts[0], parts[1]}
	}

	if err := verifyWebhookParts(webhookParts); err != nil {
		return fmt.Errorf("invalid URL format: %v", err)
	}

	config.WebhookParts = webhookParts

	for key, vals := range url.Query() {
		if err := resolver.Set(key, vals[0]); err != nil {
			return err
		}
	}

	return nil
}

func buildWebhookURL(host string, parts [4]string) string {
	return fmt.Sprintf(
		"https://%s/webhook/%s@%s/IncomingWebhook/%s/%s",
		host,
		parts[0],
		parts[1],
		parts[2],
		parts[3])
}

func parseAndVerifyWebhookURL(webhookURL string) (parts [4]string, err error) {
	pattern, err := regexp.Compile(`([0-9a-f-]{36})@([0-9a-f-]{36})/[^/]+/([0-9a-f]{32})/([0-9a-f-]{36})`)
	if err != nil {
		return parts, err
	}

	groups := pattern.FindStringSubmatch(webhookURL)
	if len(groups) != 5 {
		return parts, fmt.Errorf("invalid webhook URL format")
	}

	copy(parts[:], groups[1:])
	return parts, nil
}

const (
	// Scheme is the identifying part of this service's configuration URL
	Scheme = "teams"
	// DefaultHost is the default host for the webhook request
	DefaultHost = "outlook.office.com"
)
