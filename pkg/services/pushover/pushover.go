package pushover

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/util/webclient"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

const (
	hookURL     = "https://api.pushover.net/1/messages.json"
	contentType = "application/x-www-form-urlencoded"
)

// Service providing the notification service Pushover
type Service struct {
	standard.Standard
	config *Config
	pkr    format.PropKeyResolver
}

// EmptyConfig returns an empty types.ServiceConfig for the service
func (service *Service) EmptyConfig() types.ServiceConfig {
	return &Config{}
}

// Send a notification message to Pushover
func (service *Service) Send(message string, params *types.Params) error {
	config := service.config
	if err := service.pkr.UpdateConfigFromParams(config, params); err != nil {
		return err
	}

	device := strings.Join(config.Devices, ",")
	if err := service.sendToDevice(device, message, config); err != nil {
		return fmt.Errorf("failed to send notifications to pushover devices: %v", err)
	}

	return nil
}

func (service *Service) sendToDevice(device string, message string, config *Config) error {

	data := url.Values{
		"device":  []string{device},
		"user":    []string{config.User},
		"token":   []string{config.Token},
		"message": []string{message},
	}

	if len(config.Title) > 0 {
		data.Set("title", config.Title)
	}

	if config.Priority >= -2 && config.Priority <= 1 {
		data.Set("priority", strconv.FormatInt(int64(config.Priority), 10))
	}

	response := new(string)
	if err := webclient.PostUrl(hookURL, data, response); err != nil {
		service.Logf("Response: %q", *response)
		return err
	}

	return nil
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{}
	service.pkr = format.NewPropKeyResolver(service.config)
	if err := service.config.setURL(&service.pkr, configURL); err != nil {
		return err
	}

	return nil
}
