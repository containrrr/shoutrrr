// Package ntfy implements Ntfy as a shoutrrr service
package ntfy

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/containrrr/shoutrrr/internal/meta"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/util/jsonclient"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Service sends notifications Ntfy
type Service struct {
	standard.Standard
	config *Config
	pkr    format.PropKeyResolver
}

// Send a notification message to Ntfy
func (service *Service) Send(message string, params *types.Params) error {
	config := service.config

	if err := service.pkr.UpdateConfigFromParams(config, params); err != nil {
		return err
	}

	// Prevent enabling templated messages without custom data
	if config.Template && message == "" {
		return fmt.Errorf("body must be set if template is enabled")
	}

	if config.Template && config.Message == "" && config.Title == "" {
		return fmt.Errorf("at least title or message must be set if template is enabled")
	}

	if config.Message != "" && !config.Template {
		return fmt.Errorf("message should not be filled if template is disabled, use body instead")
	}

	if err := service.sendAPI(config, message); err != nil {
		return fmt.Errorf("failed to send ntfy notification: %w", err)
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
	request := message
	jsonClient := jsonclient.NewClient()

	headers := jsonClient.Headers()
	headers.Del("Content-Type")
	headers.Set("User-Agent", "shoutrrr/"+meta.Version)
	addHeaderIfNotEmpty(&headers, "Title", config.Title)
	addHeaderIfNotEmpty(&headers, "Priority", config.Priority.String())
	addHeaderIfNotEmpty(&headers, "Tags", strings.Join(config.Tags, ","))
	addHeaderIfNotEmpty(&headers, "Delay", config.Delay)
	addHeaderIfNotEmpty(&headers, "Actions", strings.Join(config.Actions, ";"))
	addHeaderIfNotEmpty(&headers, "Click", config.Click)
	addHeaderIfNotEmpty(&headers, "Attach", config.Attach)
	addHeaderIfNotEmpty(&headers, "X-Icon", config.Icon)
	addHeaderIfNotEmpty(&headers, "Filename", config.Filename)
	addHeaderIfNotEmpty(&headers, "Email", config.Email)

	if !config.Cache {
		headers.Add("Cache", "no")
	}
	if !config.Firebase {
		headers.Add("Firebase", "no")
	}
	if config.Template {
		headers.Add("Template", "yes")
		headers.Add("Message", config.Message)
	}

	if err := jsonClient.Post(config.GetAPIURL(), request, &response); err != nil {
		if jsonClient.ErrorResponse(err, &response) {
			// apiResponse implements Error
			return &response
		}
		return err
	}

	return nil
}

func addHeaderIfNotEmpty(headers *http.Header, key string, value string) {
	if value != "" {
		headers.Add(key, value)
	}
}
