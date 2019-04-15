package router

import (
    "errors"
    "github.com/containrrr/shoutrrr/pkg/plugins/slack"
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

func (router *ServiceRouter) RouteToSlack(url string, message string) error {
    plugin := slack.SlackPlugin{}
    plugin.Send(url, message)
    return nil
}

func (router *ServiceRouter) Route(url string, message string) error {
    svc, err := router.ExtractServiceName(url)
    if err != nil {
        return err
    }

    if strings.ToLower(svc) == "slack" {
        err := router.RouteToSlack(url, message)
        return err
    } else {
        return errors.New("unknown service")
    }
}