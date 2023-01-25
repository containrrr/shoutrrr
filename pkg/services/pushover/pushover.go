package pushover

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/format"
	"net/http"
	"net/url"
	"strconv"
	"strings"

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

// Send a notification message to Pushover
func (service *Service) Send(message string, params *types.Params) error {
	config := service.config
	if err := service.pkr.UpdateConfigFromParams(config, params); err != nil {
		return err
	}

	device := strings.Join(config.Devices, ",")
	if err := service.sendToDevice(device, message, config); err != nil {
		return fmt.Errorf("failed to send notifications to pushover devices: %w", err)
	}

	return nil
}

func (service *Service) sendToDevice(device string, message string, config *Config) error {

	data := url.Values{}
	data.Set("device", device)
	data.Set("user", config.User)
	data.Set("token", config.Token)
	data.Set("message", message)

	if len(config.Title) > 0 {
		data.Set("title", config.Title)
	}

	if config.Priority >= -2 && config.Priority <= 1 {
		data.Set("priority", strconv.FormatInt(int64(config.Priority), 10))
	}

	res, err := http.Post(
		hookURL,
		contentType,
		strings.NewReader(data.Encode()))

	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send notification to pushover device %q, response status %q", device, res.Status)
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
