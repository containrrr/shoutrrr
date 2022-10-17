package rocketchat

import (
	"bytes"
	"fmt"
	"io/ioutil"
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
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
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
		return fmt.Errorf("Error while posting to URL: %w\nHOST: %s\nPORT: %s", err, config.Host, config.Port)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		resBody, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("notification failed: %d %s", res.StatusCode, resBody)
	}
	return err
}

func buildURL(config *Config) string {
	if config.Port != "" {
		return fmt.Sprintf("https://%s:%s/hooks/%s/%s", config.Host, config.Port, config.TokenA, config.TokenB)
	}

	return fmt.Sprintf("https://%s/hooks/%s/%s", config.Host, config.TokenA, config.TokenB)
}
