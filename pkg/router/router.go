package router

import (
	"errors"
	"github.com/containrrr/shoutrrr/pkg/plugins/discord"
	"github.com/containrrr/shoutrrr/pkg/plugins/pushover"
	"github.com/containrrr/shoutrrr/pkg/plugins/slack"
	"github.com/containrrr/shoutrrr/pkg/plugins/teams"
	"github.com/containrrr/shoutrrr/pkg/plugins/telegram"
	"github.com/containrrr/shoutrrr/pkg/plugins/smtp"
	"regexp"
	"strings"
)


// ServiceRouter is responsible for routing a message to a specific notification service using the notification URL
type ServiceRouter struct {}

// ExtractServiceName from a notification URL
func (router *ServiceRouter) ExtractServiceName(url string) (string, error) {
	regex, err := regexp.Compile("^([a-zA-Z]+)://")
	if err != nil {
		return "", errors.New("could not compile regex")
	}
	match := regex.FindStringSubmatch(url)
	if len(match) <= 1 {
		return "", errors.New("could not find any service part")
	}
	return match[1], nil
}


// Route a message to a specific notification service using the notification URL
func (router *ServiceRouter) Route(url string, message string) error {
	svc, err := router.ExtractServiceName(url)
	if err != nil {
		return err
	}

	switch strings.ToLower(svc) {
	case "discord":
		return (&discord.Plugin{}).Send(url, message)
	case "pushover":
		return (&pushover.Plugin{}).Send(url, message)
	case "slack":
		return (&slack.Plugin{}).Send(url, message)
	case "teams":
		return (&teams.Plugin{}).Send(url, message)
	case "telegram":
		return (&telegram.Plugin{}).Send(url, message)
	case "smtp":
		return (&smtp.Plugin{}).Send(url, message)
	}
	return errors.New("unknown service")
}
