package router

import (
	"errors"
	"github.com/containrrr/shoutrrr/pkg/plugins/discord"
	"github.com/containrrr/shoutrrr/pkg/plugins/pushover"
	"github.com/containrrr/shoutrrr/pkg/plugins/slack"
	"github.com/containrrr/shoutrrr/pkg/plugins/teams"
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


func (router *ServiceRouter) Route(url string, message string) error {
	svc, err := router.ExtractServiceName(url)
	if err != nil {
		return err
	}

	switch strings.ToLower(svc) {
	case "discord":
		return (&discord.DiscordPlugin{}).Send(url, message)
	case "pushover":
		return (&pushover.PushoverPlugin{}).Send(url, message)
	case "slack":
		return (&slack.SlackPlugin{}).Send(url, message)
	case "teams":
		return (&teams.TeamsPlugin{}).Send(url, message)
	case "telegram":
		return (&telegram.TelegramPlugin{}).Send(url, message)
	}
	return errors.New("unknown service")
}
