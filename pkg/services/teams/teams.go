package teams

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/format"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Service providing teams as a notification service
type Service struct {
	standard.Standard
	config *Config
	pkr    format.PropKeyResolver
}

// Send a notification message to Microsoft Teams
func (service *Service) Send(message string, params *types.Params) error {
	config := service.config

	if err := service.pkr.UpdateConfigFromParams(config, params); err != nil {
		service.Logf("Failed to update params: %v", err)
	}

	return service.doSend(config, message)
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger *log.Logger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{}

	service.pkr = format.NewPropKeyResolver(service.config)

	if err := service.config.setURL(&service.pkr, configURL); err != nil {
		return err
	}

	return nil
}

func (service *Service) doSend(config *Config, message string) error {
	var sections []Section

	for _, line := range strings.Split(message, "\n") {
		sections = append(sections, Section{
			Text: line,
		})
	}

	// Teams need a summary for the webhook, use title or first (truncated) row
	summary := config.Title
	if summary == "" && len(sections) > 0 {
		summary = sections[0].Text
		if len(summary) > 20 {
			summary = summary[:21]
		}
	}

	payload, err := json.Marshal(Payload{
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

	postURL := buildWebhookURL(config.WebhookParts)

	res, err := http.Post(postURL, "application/json", bytes.NewBuffer(payload))
	if err == nil && res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send notification to teams, response status code %s", res.Status)
	}
	if err != nil {
		return fmt.Errorf(
			"an error occurred while sending notification to teams: %s",
			err.Error(),
		)
	}
	return nil
}
