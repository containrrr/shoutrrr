package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/format"
	"net/http"
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
	config *Config
	pkr    format.PropKeyResolver
}

// Send notification to Telegram
func (service *Service) Send(message string, params *types.Params) error {
	if len(message) > maxlength {
		return errors.New("message exceeds the max length")
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
	service.pkr = format.NewPropKeyResolver(service.config)
	if err := service.config.setURL(&service.pkr, configURL); err != nil {
		return err
	}

	return nil
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
	postURL := fmt.Sprintf(apiFormat, config.Token, "sendMessage")

	payload := createSendMessagePayload(message, chat, config)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	res, err := http.Post(postURL, "application/jsonData", bytes.NewBuffer(jsonData))
	if err == nil && res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send notification to \"%s\", response status code %s", chat, res.Status)
	}
	return err
}
