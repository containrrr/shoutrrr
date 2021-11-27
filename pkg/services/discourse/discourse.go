package discourse

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/common/webclient"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
	"net/url"
)

// Service providing Discourse as a notification service
type Service struct {
	standard.Standard
	config *Config
	pkr    format.PropKeyResolver
	client webclient.WebClient
}

// Send a notification message to discourse
func (service *Service) Send(message string, params *types.Params) error {
	client := service.client
	config := *service.config
	if err := service.pkr.UpdateConfigFromParams(&config, params); err != nil {
		return err
	}

	postsURL := url.URL{
		Scheme: "https",
		Host:   config.Host,
		Path:   postsEndpoint,
	}

	payload := createPostPayload{
		Title:            config.Title,
		Raw:              message,
		TargetRecipients: config.Recipients,
		Archetype:        config.Type.Archetype(),
		CreatedAt:        "",
		EmbedUrl:         config.EmbedURL,
	}

	if config.Category != 0 {
		payload.Category = &config.Category
	}
	if config.Topic != 0 {
		payload.TopicId = &config.Topic
	}

	response := createPostResponse{}
	if err := client.Post(postsURL.String(), &payload, &response); err != nil {
		errorResponse := errorResponse{}
		if client.ErrorResponse(err, &errorResponse) {
			var firstErr error
			for _, msg := range errorResponse.Errors {
				service.Logf("API Error: %q", msg)
				if firstErr == nil {
					firstErr = fmt.Errorf("discourse API: %q", msg)
				}
			}
			if firstErr != nil {
				return firstErr
			}
		}

		return err
	}

	service.Logf("Created new post #%v in topic %q (%v)", response.PostNumber, response.TopicSlug, response.TopicId)

	return nil
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

	service.client = webclient.NewJSONClient()
	service.client.Headers().Set("Api-Username", service.config.Username)
	service.client.Headers().Set("Api-Key", service.config.APIKey)

	return nil
}
