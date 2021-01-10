package pushbullet

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Service providing Pushbullet as a notification service
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

// Send ...
func (service *Service) Send(message string, params *types.Params) error {
	config := service.config
	for _, target := range config.Targets {
		if err := doSend(config, target, message, params); err != nil {
			return err
		}
	}
	return nil
}

// SendItems concatenates the items and sends them using Send
func (service *Service) SendItems(items []types.MessageItem, params *types.Params) error {
	return service.Send(types.ItemsToPlain(items), params)
}

func getTitle(params *types.Params) string {
	title := "Shoutrrr notification"
	if params != nil {
		valParams := *params
		title, ok := valParams["title"]
		if !ok {
			return title
		}
	}
	return title
}

func doSend(config *Config, target string, message string, params *types.Params) error {
	targetType, err := getTargetType(target)
	if err != nil {
		return err
	}

	apiURL := serviceURL
	json, _ := CreateJSONPayload(target, targetType, config, message, params)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", apiURL, bytes.NewReader(json))
	req.Header.Add("Access-Token", config.Token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send notification to service, response status code %s", res.Status)
	}

	if err != nil {
		return fmt.Errorf("error occurred while posting to pushbullet: %s", err.Error())
	}

	return nil
}

func getTargetType(target string) (TargetType, error) {
	matchesEmail, err := regexp.MatchString(`.*@.*\..*`, target)

	if matchesEmail && err == nil {
		return EmailTarget, nil
	}

	if len(target) > 0 && string(target[0]) == "#" {
		return ChannelTarget, nil
	}

	return DeviceTarget, nil
}

// TargetType ...
type TargetType int

const (
	// EmailTarget ...
	EmailTarget TargetType = 1
	// ChannelTarget ...
	ChannelTarget TargetType = 2
	// DeviceTarget ...
	DeviceTarget TargetType = 3
)
