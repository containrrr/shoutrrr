package pushbullet

import (
	"log"
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

const (
	serviceURL = "https://api.pushbullet.com/v2/pushes"
	// Scheme is the scheme part of the service configuration URL
	Scheme = "pushbullet"
)

var _ types.Service = &Service{}

// Send ...
func (service *Service) Send(message string, params *types.Params) error {
	config := service.config
	for _, target := range config.Targets {
		if err := doSend(config.Token, target, message); err != nil {
			return err
		}
	}
	return nil
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

func doSend(token string, target string, message string) error {
	_, err := getTargetType(target)
	if err != nil {
		return err
	}

	/*
			payload format:
			{
				"type": "note",
		        "title": title,
				"body": message,
				"x": target // replace x with email, channel_tag or device_iden based on target type
			}
	*/

	return nil
}

func getTargetType(target string) (TargetType, error) {
	matchesEmail, err := regexp.MatchString(`.*@.*\..*`, target)

	if matchesEmail && err == nil {
		return EmailTarget, nil
	} else if string(target[0]) == "#" {
		return ChannelTarget, nil
	} else {
		return DeviceTarget, nil
	}
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
