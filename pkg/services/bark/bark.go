package bark

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/util/jsonclient"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Service sends notifications Bark
type Service struct {
	standard.Standard
	config *Config
	pkr    format.PropKeyResolver
}

// Send a notification message to Bark
func (service *Service) Send(message string, params *types.Params) error {
	config := service.config

	if err := service.pkr.UpdateConfigFromParams(config, params); err != nil {
		return err
	}

	if err := service.sendAPI(config, message); err != nil {
		return fmt.Errorf("failed to send bark notification: %v", err)
	}

	return nil
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{}
	service.pkr = format.NewPropKeyResolver(service.config)

	if err := service.pkr.SetDefaultProps(service.config); err != nil {
		return err
	}

	return service.config.setURL(&service.pkr, configURL)

}

func (service *Service) sendAPI(config *Config, message string) error {
	response := apiResponse{}
	request := PushPayload{
		Body:      message,
		DeviceKey: config.DeviceKey,
		Title:     config.Title,
		Category:  config.Category,
		Copy:      config.Copy,
		Sound:     config.Sound,
		Group:     config.Group,
		Badge:     &config.Badge,
	}
	jsonClient := jsonclient.NewClient()

	if err := jsonClient.Post(config.GetAPIURL("push"), &request, &response); err != nil {
		if jsonClient.ErrorResponse(err, &response) {
			return fmt.Errorf("server response: %v", response.Message)
		}
		return err
	}

	if response.Code != http.StatusOK {
		return fmt.Errorf("unknown error")
	}

	return nil
}
