package teams

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/nicholas-fedor/shoutrrr/pkg/format"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/standard"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
)

const (
	Scheme   = "teams"
	dummyURL = "teams://dummy@dummy.com"
)

// Config represents the configuration for the Teams service.
type Config struct {
	standard.EnumlessConfig
	Group      string `optional:"" url:"user"`
	Tenant     string `optional:"" url:"host"`
	AltID      string `optional:"" url:"path1"`
	GroupOwner string `optional:"" url:"path2"`
	ExtraID    string `optional:"" url:"path3"`

	Title string `key:"title" optional:""`
	Color string `key:"color" optional:""`
	Host  string `key:"host"  optional:""` // Required, no default
}

// WebhookParts returns the webhook components as an array.
func (config *Config) WebhookParts() [5]string {
	return [5]string{config.Group, config.Tenant, config.AltID, config.GroupOwner, config.ExtraID}
}

// SetFromWebhookURL updates the Config from a Teams webhook URL.
func (config *Config) SetFromWebhookURL(webhookURL string) error {
	orgPattern, err := regexp.Compile(`https://([^.]+)` + WebhookDomain + `/`)
	if err == nil {
		orgGroups := orgPattern.FindStringSubmatch(webhookURL)
		if len(orgGroups) == 2 {
			config.Host = orgGroups[1] + WebhookDomain
		} else {
			return fmt.Errorf("invalid webhook URL format - must contain organization domain")
		}
	}
	parts, err := parseAndVerifyWebhookURL(webhookURL)
	if err != nil {
		return err
	}
	config.setFromWebhookParts(parts)
	return nil
}

// ConfigFromWebhookURL creates a new Config from a parsed Teams webhook URL.
func ConfigFromWebhookURL(webhookURL url.URL) (*Config, error) {
	webhookURL.RawQuery = "" // Clear query params
	config := &Config{
		Host: webhookURL.Host,
	}
	if err := config.SetFromWebhookURL(webhookURL.String()); err != nil {
		return nil, err
	}
	return config, nil
}

// GetURL constructs a URL from the Config fields.
func (config *Config) GetURL() *url.URL {
	resolver := format.NewPropKeyResolver(config)
	return config.getURL(&resolver)
}

func (config *Config) getURL(resolver types.ConfigQueryResolver) *url.URL {
	if config.Host == "" {
		return nil // Host is required
	}
	return &url.URL{
		User:     url.User(config.Group),
		Host:     config.Tenant,
		Path:     "/" + config.AltID + "/" + config.GroupOwner + "/" + config.ExtraID,
		Scheme:   Scheme,
		RawQuery: format.BuildQuery(resolver),
	}
}

// SetURL updates the Config from a URL.
func (config *Config) SetURL(url *url.URL) error {
	resolver := format.NewPropKeyResolver(config)
	return config.setURL(&resolver, url)
}

// setURL updates the Config from a URL using the provided resolver.
// It parses the URL parts, sets query parameters, and ensures the host is specified.
// Returns an error if the URL is invalid or the host is missing.
func (config *Config) setURL(resolver types.ConfigQueryResolver, url *url.URL) error {
	parts, err := parseURLParts(url)
	if err != nil {
		return err
	}
	config.setFromWebhookParts(parts)

	if err := config.setQueryParams(resolver, url.Query()); err != nil {
		return err
	}

	if config.Host == "" {
		return fmt.Errorf("missing required host parameter (organization.webhook.office.com)")
	}
	return nil
}

// parseURLParts extracts and validates webhook components from a URL.
// It handles path splitting and verification, returning the parts as an array.
// Returns an error if the URL format is invalid or components fail verification.
func parseURLParts(url *url.URL) ([5]string, error) {
	var parts [5]string
	if url.String() == dummyURL {
		return parts, nil
	}

	pathParts := strings.Split(url.Path, "/")
	if pathParts[0] == "" {
		pathParts = pathParts[1:]
	}
	if len(pathParts) < 3 {
		return parts, errors.New("invalid URL format: missing extraId component")
	}

	parts = [5]string{
		url.User.Username(),
		url.Hostname(),
		pathParts[0],
		pathParts[1],
		pathParts[2],
	}
	if err := verifyWebhookParts(parts); err != nil {
		return parts, fmt.Errorf("invalid URL format: %w", err)
	}
	return parts, nil
}

// setQueryParams applies query parameters to the Config using the resolver.
// It resets Color, Host, and Title, then updates them based on query values.
// Returns an error if the resolver fails to set any parameter.
func (config *Config) setQueryParams(resolver types.ConfigQueryResolver, query url.Values) error {
	config.Color = ""
	config.Host = ""
	config.Title = ""

	for key, vals := range query {
		if len(vals) > 0 && vals[0] != "" {
			switch key {
			case "color":
				config.Color = vals[0]
			case "host":
				config.Host = vals[0]
			case "title":
				config.Title = vals[0]
			}
			if err := resolver.Set(key, vals[0]); err != nil {
				return err
			}
		}
	}
	return nil
}

func (config *Config) setFromWebhookParts(parts [5]string) {
	config.Group = parts[0]
	config.Tenant = parts[1]
	config.AltID = parts[2]
	config.GroupOwner = parts[3]
	config.ExtraID = parts[4]
}
