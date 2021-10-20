package generic

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Service providing a generic notification service
type Service struct {
	standard.Standard
	config *Config
	pkr    format.PropKeyResolver
}

// EmptyConfig returns an empty types.ServiceConfig for the service
func (service *Service) EmptyConfig() types.ServiceConfig {
	return &Config{}
}

// Send a notification message to a generic webhook endpoint
func (service *Service) Send(message string, params *types.Params) error {
	config := service.config

	if err := service.pkr.UpdateConfigFromParams(config, params); err != nil {
		service.Logf("Failed to update params: %v", err)
	}

	if err := service.doSend(config, message, params); err != nil {
		return fmt.Errorf("an error occurred while sending notification to generic webhook: %s", err.Error())
	}

	return nil
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)
	config, pkr := DefaultConfig()
	service.config = config
	service.pkr = pkr

	return service.config.setURL(&service.pkr, configURL)
}

// GetConfigURLFromCustom creates a regular service URL from one with a custom host
func (*Service) GetConfigURLFromCustom(customURL *url.URL) (serviceURL *url.URL, err error) {
	webhookURL := *customURL
	if strings.HasPrefix(webhookURL.Scheme, Scheme) {
		webhookURL.Scheme = webhookURL.Scheme[len(Scheme)+1:]
	}
	config, pkr, err := ConfigFromWebhookURL(webhookURL)
	if err != nil {
		return nil, err
	}
	return config.getURL(&pkr), nil
}

func (service *Service) getPayload(template string, message string, params *types.Params) (io.Reader, error) {
	if template == "" {
		return bytes.NewBufferString(message), nil
	}
	tpl, found := service.GetTemplate(template)
	if !found {
		return nil, fmt.Errorf("template %q has not been loaded", template)
	}

	if params == nil {
		params = &types.Params{}
	}
	params.SetMessage(message)
	bb := &bytes.Buffer{}
	err := tpl.Execute(bb, params)
	return bb, err
}
