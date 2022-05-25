package telegram

import (
	"fmt"
	"html"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/util"

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

	config := *service.config
	if err := service.pkr.UpdateConfigFromParams(&config, params); err != nil {
		return err
	}

	msgs, omitted := splitMessages(&config, message)
	if len(msgs) > 1 {
		service.Logf("the message was split into %d messages", len(msgs))
	}

	var firstErr error
	for _, msg := range msgs {
		if err := service.sendMessageForChatIDs(msg, &config); err != nil {
			service.Log(err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}

	if omitted > 0 {
		service.Logf("the message exceeded the total maximum amount of characters to send. %d characters were omitted", omitted)
	}

	if firstErr != nil {
		return fmt.Errorf("failed to send telegram notification: %v", firstErr)
	}

	return nil
}

func splitMessages(config *Config, message string) (messages []string, omitted int) {
	if config.ParseMode == ParseModes.None {
		// no parse mode has been provided, treat message as unescaped HTML
		message = html.EscapeString(message)
		config.ParseMode = ParseModes.HTML
	}

	// Remove the HTML overhead and title length from the maximum message length
	maxLen := maxlength - HTMLOverhead - len(config.Title)
	messageLimits := types.MessageLimit{
		ChunkSize:      maxLen,
		TotalChunkSize: maxLen * config.MaxSplitSends,
		ChunkCount:     config.MaxSplitSends,
	}

	items, omitted := util.PartitionMessage(message, messageLimits, 10)

	title := config.Title

	messages = make([]string, len(items))
	for i, item := range items {
		messages[i] = formatHTMLMessage(title, item.Text)
		if i == 0 {
			// Skip title for all but the first message
			title = ""
		}
	}
	return messages, omitted
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)

	service.config = &Config{
		Preview:      true,
		Notification: true,
	}
	config := service.config

	service.pkr = format.NewPropKeyResolver(config)
	pkr := &service.pkr

	_ = service.pkr.SetDefaultProps(config)

	if err := config.setURL(pkr, configURL); err != nil {
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
	client := &Client{token: config.Token}
	payload := createSendMessagePayload(message, chat, config)
	_, err := client.SendMessage(&payload)
	return err
}
