package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util"
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
	hookURL = "https://discord.com/api/webhooks"
	// Only search this many runes for a good split position
	maxSearchRunes = 100
)

// Send a notification message to discord
func (service *Service) Send(message string, params *types.Params) (err error) {
	if service.config.JSON {
		postURL := CreateAPIURLFromConfig(service.config)
		err = doSend([]byte(message), postURL)
	} else {
		items, omitted := CreateItemsFromPlain(message, service.config.SplitLines)
		err = service.sendItems(items, params, omitted)
	}

	if err != nil {
		err = fmt.Errorf("failed to send discord notification: %v", err)
	}

	return
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

	var payloadBytes []byte
	payloadBytes, err = json.Marshal(payload)
	if err != nil {
		return err
	}

	postURL := CreateAPIURLFromConfig(&config)
	return doSend(payloadBytes, postURL)
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

func doSend(payload []byte, postURL string) error {
	res, err := http.Post(postURL, "application/json", bytes.NewBuffer(payload))

	if res == nil && err == nil {
		err = fmt.Errorf("unknown error")
	}

	if err == nil && res.StatusCode != http.StatusNoContent {
		err = fmt.Errorf("response status code %s", res.Status)
	}

	return err
}
