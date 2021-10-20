package join

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/util/jsonclient"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

const (
	hookURL     = "https://joinjoaomgcd.appspot.com/_ah/api/messaging/v1/sendPush"
	contentType = "text/plain"
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

	apiURL, err := url.Parse(hookURL)
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

	apiURL.RawQuery = data.Encode()

	err = jsonclient.Post(apiURL.String(), nil, nil)

	if err = jsonclient.Post(apiURL.String(), nil, nil); err != nil {
		return fmt.Errorf("failed to send notification to join devices %q: %v", devices, err)
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
