package gotify

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util/jsonclient"
)

// Service providing Gotify as a notification service
type Service struct {
	standard.Standard
	config *Config
	pkr    format.PropKeyResolver
	Client jsonclient.Client
}

// EmptyConfig returns an empty types.ServiceConfig for the service
func (service *Service) EmptyConfig() types.ServiceConfig {
	return &Config{}
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{
		Title: "Shoutrrr notification",
	}
	service.pkr = format.NewPropKeyResolver(service.config)
	err := service.config.SetURL(configURL)

	service.Client = jsonclient.NewClient()
	// service.Client = &http.Client{
	// 	Transport: &http.Transport{
	// 		TLSClientConfig: &tls.Config{
	// 			// If DisableTLS is specified, we might still need to disable TLS verification
	// 			// since the default configuration of Gotify redirects HTTP to HTTPS
	// 			// Note that this cannot be overridden using params, only using the config URL
	// 			InsecureSkipVerify: service.config.DisableTLS,
	// 		},
	// 	},
	// 	// Set a reasonable timeout to prevent one bad transfer from block all subsequent ones
	// 	Timeout: 10 * time.Second,
	// }

	return err
}

const tokenChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.-_"

// The validation rules have been taken directly from the Gotify source code.
// These will have to be adapted in case of a change:
// https://github.com/gotify/server/blob/ad157a138b4985086c484a7aabfc2deada5a33dd/auth/token.go#L8
func isTokenValid(token string) bool {
	if len(token) != 15 {
		return false
	} else if token[0] != 'A' {
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
	if len(token) > 0 && token[0] == '/' {
		token = token[1:]
	}
	if !isTokenValid(token) {
		return "", fmt.Errorf("invalid gotify token \"%s\"", token)
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
	request := payload{
		Message:  message,
		Title:    config.Title,
		Priority: config.Priority,
	}

	if err = service.Client.Post(postURL, request, nil); err != nil {
		return fmt.Errorf("failed to send notification to Gotify: %s", err)
	}

	return nil
}
