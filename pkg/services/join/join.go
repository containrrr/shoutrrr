package join

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	hookURL     = "https://joinjoaomgcd.appspot.com/_ah/api/messaging/v1/sendPush"
	contentType = "text/plain"
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

	title, found := (*params)["title"]
	if !found {
		title = config.Title
	}

	icon, found := (*params)["icon"]
	if !found {
		icon = config.Icon
	}

	devices := strings.Join(config.Devices, ",")

	return service.sendToDevices(devices, message, title, icon)
}

func (service *Service) sendToDevices(devices string, message string, title string, icon string) error {
	config := service.config

	apiUrl, err := url.Parse(hookURL)
	if err != nil {
		return err
	}

	data := url.Values{}
	data.Set("deviceIds", devices)
	data.Set("apikey", config.APIKey)
	data.Set("text", message)

	if len(title) > 0 {
		data.Set("title", title)
	}

	if len(icon) > 0 {
		data.Set("icon", icon)
	}

	apiUrl.RawQuery = data.Encode()

	res, err := http.Post(
		apiUrl.String(),
		contentType,
		nil)

	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send notification to join devices %q, response status %q", devices, res.Status)
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
