package slack

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Service sends notifications to a pre-configured channel or user
type Service struct {
	standard.Standard
	config *Config
}

const (
	apiURL    = "https://hooks.slack.com/services"
	maxlength = 1000
)

// Send a notification message to Slack
func (service *Service) Send(message string, params *types.Params) error {
	config := service.config

	if err := ValidateToken(config.Token); err != nil {
		return err
	}
	if len(message) > maxlength {
		return errors.New("message exceeds max length")
	}

	return service.doSend(config, message)
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

func (service *Service) doSend(config *Config, message string) error {
	apiURL := service.getURL(config)
	json, _ := CreateJSONPayload(config, message)
	res, err := http.Post(apiURL, "application/json", bytes.NewReader(json))

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send notification to service, response status code %s", res.Status)
	}
	return err
}

func (service *Service) getURL(config *Config) string {
	return fmt.Sprintf("%s/%s", apiURL, strings.Join(config.Token, "/"))
}
