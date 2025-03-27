package mattermost

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"

	"github.com/nicholas-fedor/shoutrrr/pkg/format"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/standard"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
)

// Service sends notifications to a pre-configured channel or user.
type Service struct {
	standard.Standard
	Config     *Config
	pkr        format.PropKeyResolver
	httpClient *http.Client
}

// GetHTTPClient returns the service's HTTP client for testing purposes.
func (service *Service) GetHTTPClient() *http.Client {
	return service.httpClient
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service.
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)
	service.Config = &Config{}
	service.pkr = format.NewPropKeyResolver(service.Config)

	err := service.Config.setURL(&service.pkr, configURL)
	if err != nil {
		return err
	}

	var transport *http.Transport
	if service.Config.DisableTLS {
		transport = &http.Transport{
			TLSClientConfig: nil, // Plain HTTP
		}
	} else {
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,            // Explicitly safe when TLS is enabled
				MinVersion:         tls.VersionTLS12, // Enforce TLS 1.2 or higher
			},
		}
	}

	service.httpClient = &http.Client{Transport: transport}

	return nil
}

// GetID returns the service identifier.
func (service *Service) GetID() string {
	return Scheme
}

// Send a notification message to Mattermost.
func (service *Service) Send(message string, params *types.Params) error {
	config := service.Config
	apiURL := buildURL(config)

	if err := service.pkr.UpdateConfigFromParams(config, params); err != nil {
		return err
	}

	json, _ := CreateJSONPayload(config, message, params)

	res, err := service.httpClient.Post(apiURL, "application/json", bytes.NewReader(json))
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send notification to service, response status code %s", res.Status)
	}

	return nil
}

// Builds the actual URL the request should go to.
func buildURL(config *Config) string {
	scheme := "https"
	if config.DisableTLS {
		scheme = "http"
	}

	return fmt.Sprintf("%s://%s/hooks/%s", scheme, config.Host, config.Token)
}
