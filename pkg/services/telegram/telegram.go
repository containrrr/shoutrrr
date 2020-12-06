package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/format"
	"log"
	"net/http"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

const (
	apiBase   = "https://api.telegram.org/bot"
	maxlength = 4096
)

// Service sends notifications to a given telegram chat
type Service struct {
	standard.Standard
	config *Config
	pkr    format.PropKeyResolver
}

// Send notification to Telegram
func (service *Service) Send(message string, _ *types.Params) error {
	if len(message) > maxlength {
		return errors.New("message exceeds the max length")
	}

	return service.sendMessageForChatIDs(message)
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger *log.Logger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{}
	service.pkr = format.NewPropKeyResolver(service.config)
	if err := service.config.setURL(&service.pkr, configURL); err != nil {
		return err
	}

	return nil
}

func (service *Service) sendMessageForChatIDs(message string) error {
	for _, channel := range service.config.Channels {
		if err := sendMessageToAPI(message, channel, service.config.Token); err != nil {
			return err
		}
	}
	return nil
}

// GetConfig returns the Config for the service
func (service *Service) GetConfig() *Config {
	return service.config
}

func sendMessageToAPI(message string, channel string, apiToken string) error {
	postURL := fmt.Sprintf("%s%s/sendMessage", apiBase, apiToken)
	jsonData, err := json.Marshal(
		JSON{
			Text: message,
			ID:   channel,
		})
	if err != nil {
		return err
	}

	res, err := http.Post(postURL, "application/jsonData", bytes.NewBuffer(jsonData))
	if err == nil && res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send notification to \"%s\", response status code %s", channel, res.Status)
	}
	return err
}
