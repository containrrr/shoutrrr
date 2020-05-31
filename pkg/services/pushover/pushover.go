package pushover

import (
	"fmt"
	"log"
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
}

// Send a notification message to Pushover
func (service *Service) Send(message string, params *types.Params) error {
	config := service.config
	if params == nil {
		params = &types.Params{}
	}
	errors := make([]error, 0)

	title, found := (*params)["subject"]
	if !found {
		title = config.Title
	}

	priority, found := (*params)["priority"]
	if !found {
		priority = strconv.FormatInt(int64(config.Priority), 10)
	}

	for _, device := range config.Devices {
		if err := service.sendToDevice(device, message, title, priority); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to send notifications to pushover devices: %v", errors)
	}

	return nil
}

func (service *Service) sendToDevice(device string, message string, title string, priority string) error {
	config := service.config

	data := url.Values{}
	data.Set("device", device)
	data.Set("user", config.User)
	data.Set("token", config.Token)
	data.Set("message", message)

	if len(title) > 0 {
		data.Set("title", title)
	}

	if len(priority) > 0 {
		data.Set("priority", priority)
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
func (service *Service) Initialize(configURL *url.URL, logger *log.Logger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{}
	if err := service.config.SetURL(configURL); err != nil {
		return err
	}

	return nil
}
