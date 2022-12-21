package ntfy

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/util/jsonclient"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Service sends notifications Ntfy
type Service struct {
	standard.Standard
	config *Config
	pkr    format.PropKeyResolver
}

// Send a notification message to Ntfy
func (service *Service) Send(message string, params *types.Params) error {
	config := service.config

	if err := service.pkr.UpdateConfigFromParams(config, params); err != nil {
		return err
	}

	if err := service.sendAPI(config, message); err != nil {
		return fmt.Errorf("failed to send ntfy notification: %w", err)
	}

	return nil
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{}
	service.pkr = format.NewPropKeyResolver(service.config)

	_ = service.pkr.SetDefaultProps(service.config)

	return service.config.setURL(&service.pkr, configURL)

}

func (service *Service) sendAPI(config *Config, message string) error {
	response := apiResponse{}
	request := message
	jsonClient := jsonclient.NewClient()
	jsonClient.Headers().Del("Content-Type")
	jsonClient.Headers().Add("Title", config.Title)
	jsonClient.Headers().Add("Priority", config.Priority.String())

	tags := strings.Join(config.Tags, ",")
	if tags != "" {
		jsonClient.Headers().Add("Tags", tags)
	}

	jsonClient.Headers().Add("Delay", config.Delay)

	actions := strings.Join(config.Actions, "; ")
	if actions != "" {
		jsonClient.Headers().Add("Actions", actions)
	}

	jsonClient.Headers().Add("Click", config.Click)
	jsonClient.Headers().Add("Attach", config.Attach)
	jsonClient.Headers().Add("Icon", config.Icon)
	jsonClient.Headers().Add("Filename", config.Filename)
	jsonClient.Headers().Add("Email", config.Email)

	if !config.Cache {
		jsonClient.Headers().Add("Cache", "no")
	}
	if !config.Firebase {
		jsonClient.Headers().Add("Firebase", "no")
	}

	if err := jsonClient.Post(config.GetAPIURL(), &request, &response); err != nil {
		if jsonClient.ErrorResponse(err, &response) {
			// apiResponse implements Error
			return &response
		}
		return err
	}

	return nil
}
