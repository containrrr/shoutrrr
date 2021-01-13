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
}

// CreateConfigFromWebhookURL creates a new config from a teams webhook URL
func CreateConfigFromWebhookURL(webhookURL string) (*Config, error) {
	parts, err := parseWebhookURL(webhookURL)
	if err != nil {
		return nil, err
	}

	return &Config{WebhookParts: parts}, nil
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
		Path:       "/" + parts[2] + parts[3] + "/",
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

	if err := VerifyWebhookParts(webhookParts); err != nil {
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

func buildWebhookURL(parts [4]string) string {
	return fmt.Sprintf(
		"https://outlook.office.com/webhook/%s@%s/IncomingWebhook/%s/%s",
		parts[0],
		parts[1],
		parts[2],
		parts[3])
}

func parseWebhookURL(webhookURL string) (parts [4]string, err error) {
	if len(webhookURL) < 195 {
		return parts, fmt.Errorf("invalid webhook URL format")
	}
	return [4]string{
		webhookURL[35:71],
		webhookURL[72:108],
		webhookURL[125:157],
		webhookURL[158:194],
	}, nil
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

// CreateConfigFromURL for use within the teams plugin
func (service *Service) CreateConfigFromURL(url *url.URL) (*Config, error) {
	config := Config{}
	err := config.SetURL(url)
	return &config, err
}

const (
	// Scheme is the identifying part of this service's configuration URL
	Scheme = "teams"
)
