// Package pinglet implements Pinglet as a shoutrrr service
package pinglet

import (
	"fmt"
	"net/url"

	"github.com/containrrr/shoutrrr/internal/meta"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util/jsonclient"
)

// Pinglet renders at most 3 badges per notification, and rejects over-length
// keys/values outright, so they are truncated client-side to keep the
// notification deliverable.
const (
	maxBadgeCount    = 3
	maxBadgeKeyLen   = 24
	maxBadgeValueLen = 32
	maxDataKeyLen    = 64
	maxDataValueLen  = 256
)

// Service sends notifications to Pinglet
type Service struct {
	standard.Standard
	config *Config
	pkr    format.PropKeyResolver
}

// Send a notification message to Pinglet
func (service *Service) Send(message string, params *types.Params) error {
	config := service.config

	if err := service.pkr.UpdateConfigFromParams(config, params); err != nil {
		return err
	}

	if err := service.sendAPI(config, message); err != nil {
		return fmt.Errorf("failed to send pinglet notification: %w", err)
	}

	return nil
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{}
	service.pkr = format.NewPropKeyResolver(service.config)

	_ = service.pkr.SetDefaultProps(service.config)

	return service.config.setURL(&service.pkr, configURL)
}

func (service *Service) sendAPI(config *Config, message string) error {
	response := apiResponse{}
	payload := pushPayload{
		Message:  message,
		Title:    config.Title,
		Priority: config.Priority.String(),
		Badges:   service.truncatedBadges(config),
		Data:     service.truncatedData(config),
	}

	jsonClient := jsonclient.NewClient()
	headers := jsonClient.Headers()
	headers.Set("User-Agent", "shoutrrr/"+meta.Version)
	headers.Set("Authorization", "Bearer "+config.Token)

	if err := jsonClient.Post(config.GetAPIURL(), &payload, &response); err != nil {
		if jsonClient.ErrorResponse(err, &response) {
			// apiResponse implements Error
			return &response
		}
		return err
	}

	return nil
}

// truncatedBadges returns at most maxBadgeCount badges with keys and values
// truncated to the lengths accepted by the Pinglet API
func (service *Service) truncatedBadges(config *Config) map[string]string {
	if len(config.Badges) == 0 {
		return nil
	}

	badges := make(map[string]string, maxBadgeCount)
	for key, value := range config.Badges {
		if len(badges) >= maxBadgeCount {
			service.Logger.Logf("Pinglet renders at most %d badges; additional entries were ignored.", maxBadgeCount)
			break
		}
		if len(key) > maxBadgeKeyLen || len(value) > maxBadgeValueLen {
			service.Logger.Logf("Pinglet badge %q was truncated.", key)
		}
		badges[truncate(key, maxBadgeKeyLen)] = truncate(value, maxBadgeValueLen)
	}
	return badges
}

// truncatedData returns metadata with keys and values truncated to the lengths
// accepted by the Pinglet API
func (service *Service) truncatedData(config *Config) map[string]string {
	if len(config.Data) == 0 {
		return nil
	}

	data := make(map[string]string, len(config.Data))
	for key, value := range config.Data {
		if len(key) > maxDataKeyLen || len(value) > maxDataValueLen {
			service.Logger.Logf("Pinglet metadata %q was truncated.", key)
		}
		data[truncate(key, maxDataKeyLen)] = truncate(value, maxDataValueLen)
	}
	return data
}

func truncate(value string, maxLen int) string {
	if len(value) > maxLen {
		return value[:maxLen]
	}
	return value
}
