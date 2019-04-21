package router

import (
	"errors"
	"github.com/containrrr/shoutrrr/pkg/plugins"
	"github.com/containrrr/shoutrrr/pkg/plugins/discord"
	"github.com/containrrr/shoutrrr/pkg/plugins/pushover"
	"github.com/containrrr/shoutrrr/pkg/plugins/slack"
	"github.com/containrrr/shoutrrr/pkg/plugins/telegram"
	"regexp"
	"strings"
)

type ServiceRouter struct {
}

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

func (router *ServiceRouter) RouteToPlugin(plugin plugins.Plugin, url string, message string) error {
	plugin.Send(url, message)
	return nil
}

func (router *ServiceRouter) Route(url string, message string) error {
	svc, err := router.ExtractServiceName(url)
	if err != nil {
		return err
	}

	var plugin plugins.Plugin

	switch strings.ToLower(svc) {
	case "slack":
		plugin = &slack.SlackPlugin{}
	case "telegram":
		plugin = &telegram.TelegramPlugin{}
	case "discord":
		plugin = &discord.DiscordPlugin{}
	case "pushover":
		plugin = &pushover.PushoverPlugin{}
	default:
		return errors.New("unknown service")
	}

	return router.RouteToPlugin(plugin, url, message)
}
