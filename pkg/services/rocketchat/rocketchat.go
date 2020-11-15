package rocketchat

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Service sends notifications to a pre-configured channel or user
type Service struct {
	standard.Standard
	config *Config
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

// Send a notification message to Rocket.chat
func (service *Service) Send(message string, params *types.Params) error {
	var res *http.Response
	var err error
	config := service.config
	apiURL := buildURL(config)
	json, _ := CreateJSONPayload(config, message, params)
	res, err = http.Post(apiURL, "application/json", bytes.NewReader(json))
	if err != nil {
		return fmt.Errorf("Error while posting to URL: %v\nHOST: %s\nPORT: %s\n", err, config.Host, config.Port)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send notification to service, response status code %s", res.Status)
	}
	return err
}

func buildURL(config *Config) string {
	if config.Port != "" {
		return fmt.Sprintf("https://%s:%s/hooks/%s/%s", config.Host, config.Port, config.TokenA, config.TokenB)
	} else {
		return fmt.Sprintf("https://%s/hooks/%s/%s", config.Host, config.TokenA, config.TokenB)
	}
}
