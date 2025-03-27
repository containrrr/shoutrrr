package bark

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/nicholas-fedor/shoutrrr/pkg/format"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/standard"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
	"github.com/nicholas-fedor/shoutrrr/pkg/util/jsonclient"
)

// Service sends notifications to Bark.
type Service struct {
	standard.Standard
	Config *Config // Changed from 'config' to 'Config'
	pkr    format.PropKeyResolver
}

// Send a notification message to Bark.
func (service *Service) Send(message string, params *types.Params) error {
	config := service.Config // Update reference

	if err := service.pkr.UpdateConfigFromParams(config, params); err != nil {
		return err
	}

	if err := service.sendAPI(config, message); err != nil {
		return fmt.Errorf("failed to send bark notification: %w", err)
	}

	return nil
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service.
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)
	service.Config = &Config{}                              // Update reference
	service.pkr = format.NewPropKeyResolver(service.Config) // Update reference

	_ = service.pkr.SetDefaultProps(service.Config) // Update reference

	return service.Config.setURL(&service.pkr, configURL) // Update reference
}

// GetID returns the service identifier.
func (service *Service) GetID() string {
	return Scheme
}

func (service *Service) sendAPI(config *Config, message string) error {
	response := APIResponse{}
	request := PushPayload{
		Body:      message,
		DeviceKey: config.DeviceKey,
		Title:     config.Title,
		Category:  config.Category,
		Copy:      config.Copy,
		Sound:     config.Sound,
		Group:     config.Group,
		Badge:     &config.Badge,
		Icon:      config.Icon,
		URL:       config.URL,
	}
	jsonClient := jsonclient.NewClient()

	if err := jsonClient.Post(config.GetAPIURL("push"), &request, &response); err != nil {
		if jsonClient.ErrorResponse(err, &response) {
			// apiResponse implements Error
			return &response
		}

		return err
	}

	if response.Code != http.StatusOK {
		return fmt.Errorf("unknown error")
	}

	return nil
}
