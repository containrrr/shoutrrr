package teams

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/nicholas-fedor/shoutrrr/pkg/format"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/standard"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
)

const (
	MaxSummaryLength    = 20
	TruncatedSummaryLen = 21
)

// Service implements the Teams notification service.
type Service struct {
	standard.Standard
	Config *Config
	pkr    format.PropKeyResolver
}

// Send delivers a notification message to Microsoft Teams.
func (service *Service) Send(message string, params *types.Params) error {
	config := service.Config
	if err := service.pkr.UpdateConfigFromParams(config, params); err != nil {
		service.Logf("Failed to update params: %v", err)
	}
	return service.doSend(config, message)
}

// Initialize sets up the Service with a URL and logger.
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)
	service.Config = &Config{}
	service.pkr = format.NewPropKeyResolver(service.Config)
	return service.Config.SetURL(configURL)
}

// GetID returns the service identifier.
func (service *Service) GetID() string {
	return Scheme
}

// GetConfigURLFromCustom converts a custom URL to a service URL.
func (service *Service) GetConfigURLFromCustom(customURL *url.URL) (*url.URL, error) {
	webhookURLStr := strings.TrimPrefix(customURL.String(), "teams+")
	tempURL, err := url.Parse(webhookURLStr)
	if err != nil {
		return nil, err
	}
	webhookURL := &url.URL{
		Scheme: tempURL.Scheme,
		Host:   tempURL.Host,
		Path:   tempURL.Path,
	}
	config, err := ConfigFromWebhookURL(*webhookURL)
	if err != nil {
		return nil, err
	}
	config.Color = ""
	config.Title = ""

	query := customURL.Query()
	for key, vals := range query {
		if vals[0] != "" {
			switch key {
			case "color":
				config.Color = vals[0]
			case "host":
				config.Host = vals[0]
			case "title":
				config.Title = vals[0]
			}
		}
	}
	return config.GetURL(), nil
}

func (service *Service) doSend(config *Config, message string) error {
	lines := strings.Split(message, "\n")
	sections := make([]section, 0, len(lines))
	for _, line := range lines {
		sections = append(sections, section{Text: line})
	}
	summary := config.Title
	if summary == "" && len(sections) > 0 {
		summary = sections[0].Text
		if len(summary) > MaxSummaryLength {
			summary = summary[:TruncatedSummaryLen]
		}
	}
	payload, err := json.Marshal(payload{
		CardType:   "MessageCard",
		Context:    "http://schema.org/extensions",
		Markdown:   true,
		Title:      config.Title,
		ThemeColor: config.Color,
		Summary:    summary,
		Sections:   sections,
	})
	if err != nil {
		return err
	}
	if config.Host == "" {
		return fmt.Errorf("host is required but not specified in the configuration")
	}
	postURL := BuildWebhookURL(
		config.Host,
		config.Group,
		config.Tenant,
		config.AltID,
		config.GroupOwner,
		config.ExtraID,
	)
	res, err := http.Post(postURL, "application/json", bytes.NewBuffer(payload))
	if err == nil && res.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"failed to send notification to teams, response status code %s",
			res.Status,
		)
	}
	if err != nil {
		return fmt.Errorf("an error occurred while sending notification to teams: %s", err.Error())
	}
	return nil
}
