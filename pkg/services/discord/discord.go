package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/nicholas-fedor/shoutrrr/pkg/format"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/standard"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
	"github.com/nicholas-fedor/shoutrrr/pkg/util"
)

// Service providing Discord as a notification service.
type Service struct {
	standard.Standard
	Config *Config
	pkr    format.PropKeyResolver
}

// Message limit constants.
const (
	ChunkSize      = 2000 // Maximum size of a single message chunk
	TotalChunkSize = 6000 // Maximum total size of all chunks
	ChunkCount     = 10   // Maximum number of chunks allowed
	MaxSearchRunes = 100  // Maximum number of runes to search for split position
	HooksBaseURL   = "https://discord.com/api/webhooks"
)

var limits = types.MessageLimit{
	ChunkSize:      ChunkSize,
	TotalChunkSize: TotalChunkSize,
	ChunkCount:     ChunkCount,
}

// Send a notification message to discord.
func (service *Service) Send(message string, params *types.Params) error {
	var firstErr error

	if service.Config.JSON {
		postURL := CreateAPIURLFromConfig(service.Config)
		firstErr = doSend([]byte(message), postURL)
	} else {
		batches := CreateItemsFromPlain(message, service.Config.SplitLines)
		for _, items := range batches {
			if err := service.sendItems(items, params); err != nil {
				service.Log(err)

				if firstErr == nil {
					firstErr = err
				}
			}
		}
	}

	if firstErr != nil {
		return fmt.Errorf("failed to send discord notification: %v", firstErr)
	}

	return nil
}

// SendItems sends items with additional meta data and richer appearance.
func (service *Service) SendItems(items []types.MessageItem, params *types.Params) error {
	return service.sendItems(items, params)
}

func (service *Service) sendItems(items []types.MessageItem, params *types.Params) error {
	var err error

	config := *service.Config
	if err = service.pkr.UpdateConfigFromParams(&config, params); err != nil {
		return err
	}

	var payload WebhookPayload

	payload, err = CreatePayloadFromItems(items, config.Title, config.LevelColors())
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

// CreateItemsFromPlain creates a set of MessageItems that is compatible with Discords webhook payload.
func CreateItemsFromPlain(plain string, splitLines bool) (batches [][]types.MessageItem) {
	if splitLines {
		return util.MessageItemsFromLines(plain, limits)
	}

	for {
		items, omitted := util.PartitionMessage(plain, limits, MaxSearchRunes)
		batches = append(batches, items)

		if omitted == 0 {
			break
		}

		plain = plain[len(plain)-omitted:]
	}

	return
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service.
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)
	service.Config = &Config{}
	service.pkr = format.NewPropKeyResolver(service.Config)

	if err := service.pkr.SetDefaultProps(service.Config); err != nil {
		return err
	}

	if err := service.Config.SetURL(configURL); err != nil {
		return err
	}

	return nil
}

// GetID returns the service identifier.
func (service *Service) GetID() string {
	return Scheme
}

// CreateAPIURLFromConfig takes a discord config object and creates a post url.
func CreateAPIURLFromConfig(config *Config) string {
	return fmt.Sprintf(
		"%s/%s/%s",
		HooksBaseURL,
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
