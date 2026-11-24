package gotify

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util/jsonclient"
)

// Service providing Gotify as a notification service
type Service struct {
	standard.Standard
	config     *Config
	pkr        format.PropKeyResolver
	httpClient *http.Client
	client     jsonclient.Client
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{
		Title: "Shoutrrr notification",
	}
	service.pkr = format.NewPropKeyResolver(service.config)
	err := service.config.SetURL(configURL)

	service.httpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				// If DisableTLS is specified, we might still need to disable TLS verification
				// since the default configuration of Gotify redirects HTTP to HTTPS
				// Note that this cannot be overridden using params, only using the config URL
				InsecureSkipVerify: service.config.DisableTLS,
			},
		},
		// Set a reasonable timeout to prevent one bad transfer from block all subsequent ones
		Timeout: 10 * time.Second,
	}
	service.client = jsonclient.NewWithHTTPClient(service.httpClient)

	return err
}

const tokenChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.-_"

// isTokenValid checks application-token syntax for legacy and Gotify 3 tokens.
// Modern tokens may include signed representations, so their length is not fixed.
// Gotify validates key encoding, ownership, signatures, and expiration:
// https://github.com/gotify/server/blob/v3.0.0/auth/token.go
func isTokenValid(token string) bool {
	legacy := len(token) == 15 && token[0] == 'A'
	modern := strings.HasPrefix(token, "gtfya.") && len(token) > len("gtfya.")
	if !legacy && !modern {
		return false
	}
	for _, c := range token {
		if !strings.ContainsRune(tokenChars, c) {
			return false
		}
	}
	return true
}

func buildURL(config *Config) (string, error) {
	token := config.Token
	if !isTokenValid(token) {
		return "", fmt.Errorf("invalid gotify token %q", token)
	}
	scheme := "https"
	if config.DisableTLS {
		scheme = scheme[:4]
	}
	return fmt.Sprintf("%s://%s%s/message?token=%s", scheme, config.Host, config.Path, token), nil
}

// Send a notification message to Gotify
func (service *Service) Send(message string, params *types.Params) error {
	if params == nil {
		params = &types.Params{}
	}
	config := service.config
	if err := service.pkr.UpdateConfigFromParams(config, params); err != nil {
		service.Logf("Failed to update params: %v", err)
	}

	postURL, err := buildURL(config)
	if err != nil {
		return err
	}

	request := &messageRequest{
		Message:  message,
		Title:    config.Title,
		Priority: config.Priority,
	}
	response := &messageResponse{}
	err = service.client.Post(postURL, request, response)
	if err != nil {
		errorRes := &errorResponse{}
		if service.client.ErrorResponse(err, errorRes) {
			return errorRes
		}
		return fmt.Errorf("failed to send notification to Gotify: %s", err)
	}

	return nil
}

// GetHTTPClient is only supposed to be used for mocking the httpclient when testing
func (service *Service) GetHTTPClient() *http.Client {
	return service.httpClient
}
