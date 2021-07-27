package telegram

import (
	"errors"
	"github.com/containrrr/shoutrrr/pkg/common/webclient"
	"github.com/containrrr/shoutrrr/pkg/format"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

const (
	apiFormat = "https://api.telegram.org/bot%s/%s"
	maxlength = 4096
)

// Service sends notifications to a given telegram chat
type Service struct {
	standard.Standard
	webclient.ClientService
	config *Config
	pkr    format.PropKeyResolver
}

// Send notification to Telegram
func (service *Service) Send(message string, params *types.Params) error {
	if len(message) > maxlength {
		return errors.New("Message exceeds the max length")
	}

	config := *service.config
	if err := service.pkr.UpdateConfigFromParams(&config, params); err != nil {
		return err
	}

	return service.sendMessageForChatIDs(message, &config)
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{
		Preview:      true,
		Notification: true,
	}
	service.ClientService.Initialize()
	service.pkr = format.NewPropKeyResolver(service.config)
	if err := service.config.setURL(&service.pkr, configURL); err != nil {
		return err
	}

	return nil
}

func (service *Service) sendMessageForChatIDs(message string, config *Config) error {
	client := &Client{token: config.Token, WebClient: service.WebClient()}
	for _, chat := range service.config.Chats {
		payload := createSendMessagePayload(message, chat, config)
		if _, err := client.SendMessage(&payload); err != nil {
			return err
		}
	}
	return nil
}

// GetConfig returns the Config for the service
func (service *Service) GetConfig() *Config {
	return service.config
}
