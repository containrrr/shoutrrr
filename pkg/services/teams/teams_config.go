//go:generate go run ../../../cmd/shoutrrr-gen --lang go
package teams

import (
	"fmt"
	"net/url"
	"regexp"
)

func (config *Config) webhookParts() [4]string {
	return [4]string{config.Group, config.Tenant, config.AltID, config.GroupOwner}
}

// SetFromWebhookURL updates the config WebhookParts from a teams webhook URL
func (config *Config) SetFromWebhookURL(webhookURL string) error {
	parts, err := parseAndVerifyWebhookURL(webhookURL)
	if err != nil {
		return err
	}

	config.setFromWebhookParts(parts)
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

func (config *Config) setFromWebhookParts(parts [4]string) {
	config.Group = parts[0]
	config.Tenant = parts[1]
	config.AltID = parts[2]
	config.GroupOwner = parts[3]
}

func buildWebhookURL(host, group, tenant, altID, groupOwner string) string {
	// config.Group, config.Tenant, config.AltID, config.GroupOwner
	path := Path
	if host == LegacyHost {
		path = LegacyPath
	}
	return fmt.Sprintf(
		"https://%s/%s/%s@%s/%s/%s/%s",
		host,
		path,
		group,
		tenant,
		ProviderName,
		altID,
		groupOwner)
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
	// LegacyHost is the default host for legacy webhook requests
	LegacyHost = "outlook.office.com"
	// LegacyPath is the initial path of the webhook URL for legacy webhook requests
	LegacyPath = "webhook"
	// Path is the initial path of the webhook URL for domain-scoped webhook requests
	Path = "webhookb2"
	// ProviderName is the name of the Teams integration provider
	ProviderName = "IncomingWebhook"
)
