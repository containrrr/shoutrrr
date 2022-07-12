package telegram

import (
	"errors"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/conf"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

const (
	apiFormat = "https://%s/bot%s/%s"
	maxlength = 4096
)

// Service sends notifications to a given telegram chat
type Service struct {
	standard.Standard
	config *Config
}

// Send notification to Telegram
func (service *Service) Send(message string, params *types.Params) error {
	if len(message) > maxlength {
		return errors.New("Message exceeds the max length")
	}

	config := *service.config

	if err := conf.UpdateFromParams(&config, params); err != nil {
		return err
	}

	return service.sendMessageForChatIDs(message, &config)
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{}
	return conf.Init(service.config, configURL)
}

func (service *Service) sendMessageForChatIDs(message string, config *Config) error {
	for _, chat := range service.config.Chats {
		if err := sendMessageToAPI(message, chat, config); err != nil {
			return err
		}
	}
	return nil
}

// GetConfig returns the Config for the service
func (service *Service) GetConfig() *Config {
	return service.config
}

func sendMessageToAPI(message string, chat string, config *Config) error {
	client := &Client{token: config.token(), apiHost: config.apiHost()}
	payload := createSendMessagePayload(message, chat, config)
	_, err := client.SendMessage(&payload)
	return err
}
