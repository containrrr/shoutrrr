// Package ntfy implements Ntfy as a shoutrrr service
package ntfy

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/nicholas-fedor/shoutrrr/internal/meta"
	"github.com/nicholas-fedor/shoutrrr/pkg/format"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/standard"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
	"github.com/nicholas-fedor/shoutrrr/pkg/util/jsonclient"
)

// Service sends notifications Ntfy.
type Service struct {
	standard.Standard
	Config *Config
	pkr    format.PropKeyResolver
}

// Send a notification message to Ntfy.
func (service *Service) Send(message string, params *types.Params) error {
	config := service.Config

	if err := service.pkr.UpdateConfigFromParams(config, params); err != nil {
		return err
	}

	if err := service.sendAPI(config, message); err != nil {
		return fmt.Errorf("failed to send ntfy notification: %w", err)
	}

	return nil
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service.
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)
	service.Config = &Config{}
	service.pkr = format.NewPropKeyResolver(service.Config)

	_ = service.pkr.SetDefaultProps(service.Config)

	return service.Config.setURL(&service.pkr, configURL)
}

// GetID returns the service identifier.
func (service *Service) GetID() string {
	return Scheme
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
