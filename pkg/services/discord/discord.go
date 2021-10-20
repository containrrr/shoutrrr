package discord

import (
	"fmt"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util"
	"github.com/containrrr/shoutrrr/pkg/util/jsonclient"
)

// Service providing Discord as a notification service
type Service struct {
	standard.Standard
	config *Config
	pkr    format.PropKeyResolver
}

var limits = types.MessageLimit{
	ChunkSize:      2000,
	TotalChunkSize: 6000,
	ChunkCount:     10,
}

const (
	hookURL = "https://discordapp.com/api/webhooks"
	// Only search this many runes for a good split position
	maxSearchRunes = 100
)

// EmptyConfig returns an empty types.ServiceConfig for the service
func (service *Service) EmptyConfig() types.ServiceConfig {
	return &Config{}
}

// Send a notification message to discord
func (service *Service) Send(message string, params *types.Params) error {

	if service.config.JSON {
		postURL := CreateAPIURLFromConfig(service.config)
		return doSend([]byte(message), postURL)
	}

	items, omitted := CreateItemsFromPlain(message, service.config.SplitLines)
	return service.sendItems(items, params, omitted)
}

// SendItems sends items with additional meta data and richer appearance
func (service *Service) SendItems(items []types.MessageItem, params *types.Params) error {
	return service.sendItems(items, params, 0)
}

func (service *Service) sendItems(items []types.MessageItem, params *types.Params, omitted int) error {
	var err error

	config := *service.config
	if err = service.pkr.UpdateConfigFromParams(&config, params); err != nil {
		return err
	}

	var payload WebhookPayload
	payload, err = CreatePayloadFromItems(items, config.Title, config.LevelColors(), omitted)
	if err != nil {
		return err
	}

	payload.Username = config.Username
	payload.AvatarURL = config.Avatar

	postURL := CreateAPIURLFromConfig(&config)
	return doSend(payload, postURL)
}

// CreateItemsFromPlain creates a set of MessageItems that is compatible with Discords webhook payload
func CreateItemsFromPlain(plain string, splitLines bool) (items []types.MessageItem, omitted int) {
	if splitLines {
		return util.MessageItemsFromLines(plain, limits)
	}

	return util.PartitionMessage(plain, limits, maxSearchRunes)
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{}
	service.pkr = format.NewPropKeyResolver(service.config)

	if err := service.pkr.SetDefaultProps(service.config); err != nil {
		return err
	}

	if err := service.config.SetURL(configURL); err != nil {
		return err
	}

	return nil
}

// CreateAPIURLFromConfig takes a discord config object and creates a post url
func CreateAPIURLFromConfig(config *Config) string {
	return fmt.Sprintf(
		"%s/%s/%s",
		hookURL,
		config.WebhookID,
		config.Token)
}

type payloadResponse struct{}

func doSend(payload interface{}, postURL string) error {
	if err := jsonclient.Post(postURL, payload, nil); err != nil {
		return fmt.Errorf("failed to send discord notification: %v", err)
	}

	return nil
}
