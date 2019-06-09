package router

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/plugin"
	"github.com/containrrr/shoutrrr/pkg/plugins/discord"
	"github.com/containrrr/shoutrrr/pkg/plugins/pushover"
	"github.com/containrrr/shoutrrr/pkg/plugins/slack"
	"github.com/containrrr/shoutrrr/pkg/plugins/smtp"
	"github.com/containrrr/shoutrrr/pkg/plugins/teams"
	"github.com/containrrr/shoutrrr/pkg/plugins/telegram"
	"net/url"
	"strings"
)


// ServiceRouter is responsible for routing a message to a specific notification service using the notification URL
type ServiceRouter struct {}

// ExtractServiceName from a notification URL
func (router *ServiceRouter) ExtractServiceName(rawUrl string) (string, url.URL, error) {
	if u, err := url.Parse(rawUrl); err != nil {
		return "", url.URL{}, err
	} else {
		return u.Scheme, *u, nil
	}
}


// Route a message to a specific notification service using the notification URL
func (router *ServiceRouter) Route(rawUrl string, message string, opts plugin.PluginOpts) error {
	svc, url, err := router.ExtractServiceName(rawUrl)
	if err != nil {
		return err
	}

	if service, err := router.Locate(svc); err != nil {
		return err
	} else {
		return service.Send(url, message, opts)
	}
}

var plugins = map[string]plugin.Plugin {
	"discord":	&discord.Plugin{},
	"pushover":	&pushover.Plugin{},
	"slack":	&slack.Plugin{},
	"teams":	&teams.Plugin{},
	"telegram":	&telegram.Plugin{},
	"smtp":	&smtp.Plugin{},
}

func (router *ServiceRouter) Locate(serviceScheme string) (plugin.Plugin, error) {

	service, valid := plugins[strings.ToLower(serviceScheme)]
	if !valid {
		return nil, fmt.Errorf("unknown service scheme '%s'", serviceScheme)
	}

	return service, nil
}
