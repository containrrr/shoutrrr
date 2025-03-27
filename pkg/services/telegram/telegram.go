package telegram

import (
	"errors"
	"net/url"

	"github.com/nicholas-fedor/shoutrrr/pkg/format"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/standard"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
)

const (
	apiFormat = "https://api.telegram.org/bot%s/%s"
	maxlength = 4096
)

// Service sends notifications to a given telegram chat.
type Service struct {
	standard.Standard
	Config *Config
	pkr    format.PropKeyResolver
}

// Send notification to Telegram.
func (service *Service) Send(message string, params *types.Params) error {
	if len(message) > maxlength {
		return errors.New("Message exceeds the max length")
	}

	config := *service.Config
	if err := service.pkr.UpdateConfigFromParams(&config, params); err != nil {
		return err
	}

	return service.sendMessageForChatIDs(message, &config)
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service.
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)
	service.Config = &Config{
		Preview:      true,
		Notification: true,
	}
	service.pkr = format.NewPropKeyResolver(service.Config)

	if err := service.Config.setURL(&service.pkr, configURL); err != nil {
		return err
	}

	return nil
}

// GetID returns the service identifier.
func (service *Service) GetID() string {
	return Scheme
}

func (service *Service) sendMessageForChatIDs(message string, config *Config) error {
	for _, chat := range service.Config.Chats {
		if err := sendMessageToAPI(message, chat, config); err != nil {
			return err
		}
	}

	return nil
}

// GetConfig returns the Config for the service.
func (service *Service) GetConfig() *Config {
	return service.Config
}

func sendMessageToAPI(message string, chat string, config *Config) error {
	client := &Client{token: config.Token}
	payload := createSendMessagePayload(message, chat, config)
	_, err := client.SendMessage(&payload)

	return err
}
