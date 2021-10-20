package googlechat

import (
	"fmt"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util/jsonclient"
)

// Service providing Google Chat as a notification service.
type Service struct {
	standard.Standard
	config *Config
}

// EmptyConfig returns an empty types.ServiceConfig for the service
func (service *Service) EmptyConfig() types.ServiceConfig {
	return &Config{}
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service.
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{}

	err := service.config.SetURL(configURL)

	return err
}

// Send a notification message to Google Chat.
func (service *Service) Send(message string, _ *types.Params) error {
	config := service.config

	request := payload{
		Text: message,
	}

	postURL := getAPIURL(config)

	if err := jsonclient.Post(postURL.String(), request, nil); err != nil {
		return fmt.Errorf("failed to send notification to Google Chat: %s", err)
	}

	return nil
}

func getAPIURL(config *Config) *url.URL {
	query := url.Values{}
	query.Set("key", config.Key)
	query.Set("token", config.Token)

	return &url.URL{
		Path:     config.Path,
		Host:     config.Host,
		Scheme:   "https",
		RawQuery: query.Encode(),
	}
}
