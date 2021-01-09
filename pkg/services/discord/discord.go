package discord

import (
	"bytes"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types/rich"
	"github.com/containrrr/shoutrrr/pkg/util"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Service providing Discord as a notification service
type Service struct {
	standard.Standard
	config *Config
	pkr    format.PropKeyResolver
}

const (
	hookURL             = "https://discordapp.com/api/webhooks"
	maxContentLength    = 2000
	maxTotalEmbedLength = 6000
	// Technically, the max is 10, but we use the first one for meta
	maxEmbedCount = 9
	// Only search this many runes for a good split position
	maxSearchRunes = 100
)

// Send a notification message to discord
func (service *Service) Send(message string, params *types.Params) error {

	config := *service.config
	if err := service.pkr.UpdateConfigFromParams(&config, params); err != nil {
		return err
	}

	var payload []byte
	var err error
	if config.JSON {
		payload = []byte(message)
	} else {
		items, omitted := CreateItemsFromPlain(&config, message)

		payload, err = CreatePayloadFromItems(items, config.Title, config.LevelColors(), omitted)
		if err != nil {
			return err
		}
	}

	postURL := CreateAPIURLFromConfig(service.config)

	return doSend(payload, postURL)
}

func CreateItemsFromPlain(config *Config, plain string) (items []rich.MessageItem, omitted int) {
	items = make([]rich.MessageItem, 0, maxEmbedCount)
	omitted = 0

	var lines []string
	if config.SplitLines {
		totalLength := 0
		for l, line := range strings.Split(plain, "\n") {
			if l < maxEmbedCount && totalLength < maxTotalEmbedLength {
				runes := []rune(line)
				maxLen := maxContentLength
				if totalLength+maxLen > maxTotalEmbedLength {
					maxLen = maxTotalEmbedLength - totalLength
				}
				if len(runes) > maxLen {
					// Trim and add ellipsis
					runes = runes[:maxLen-6]
					line = string(runes) + " [...]"
				}

				items = append(items, rich.MessageItem{
					Text: line,
				})

				totalLength += len(runes)

			} else {
				omitted += len(line)
			}
		}
	} else {
		lines, omitted = PartitionString(plain, maxContentLength, maxSearchRunes, maxEmbedCount, maxTotalEmbedLength)
		for _, line := range lines {
			items = append(items, rich.MessageItem{
				Text: line,
			})
		}
	}

	return items, omitted
}

func PartitionString(input string, chunkSize int, distance int, maxCount int, maxTotal int) ([]string, int) {
	runes := []rune(input)
	chunkOffset := 0
	maxTotal = util.Min(len(runes), maxTotal)

	var chunks []string
	for i := 0; i < maxCount; i++ {
		// If no suitable split point is found, use the chunkSize
		chunkEnd := chunkOffset + chunkSize
		// ... and start next chunk directly after this one
		nextChunkStart := chunkEnd
		if chunkEnd > maxTotal {
			// The chunk is smaller than the limit, no need to search
			chunkEnd = maxTotal
			nextChunkStart = maxTotal
		} else {
			for r := 0; r < distance; r++ {
				rp := chunkEnd - r
				if runes[rp] == '\n' || runes[rp] == ' ' {
					// Suitable split point found
					chunkEnd = rp
					// Since the split is on a whitespace, skip it in the next chunk
					nextChunkStart = chunkEnd + 1
					break
				}
			}
		}

		chunks = append(chunks, string(runes[chunkOffset:chunkEnd]))

		chunkOffset = nextChunkStart
		if chunkOffset >= maxTotal {
			break
		}
	}

	return chunks, len(runes) - chunkOffset
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger *log.Logger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{}
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
		config.Channel,
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

	if err != nil {
		return fmt.Errorf("failed to send discord notification: %v", err)
	}

	return nil
}
