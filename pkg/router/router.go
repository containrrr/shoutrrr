package router

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/services/discord"
	"github.com/containrrr/shoutrrr/pkg/services/pushover"
	"github.com/containrrr/shoutrrr/pkg/services/slack"
	"github.com/containrrr/shoutrrr/pkg/services/smtp"
	"github.com/containrrr/shoutrrr/pkg/services/teams"
	"github.com/containrrr/shoutrrr/pkg/services/telegram"
	"github.com/containrrr/shoutrrr/pkg/types"
	"net/url"
	"strings"
)


// ServiceRouter is responsible for routing a message to a specific notification service using the notification URL
type ServiceRouter struct {}

// ExtractServiceName from a notification URL
func (router *ServiceRouter) ExtractServiceName(rawURL string) (string, *url.URL, error) {
	serviceURL, err := url.Parse(rawURL)

	if err != nil {
		return "", &url.URL{}, err
	}
	return serviceURL.Scheme, serviceURL, nil
}


// Route a message to a specific notification service using the notification URL
func (router *ServiceRouter) Route(rawURL string, message string, opts types.ServiceOpts) error {
	svc, url, err := router.ExtractServiceName(rawURL)
	if err != nil {
		return err
	}

	service, err := router.Locate(svc)
	if err != nil {
		return err
	}

	return service.Send(url, message, opts)
}

var services = map[string]types.Service {
	"discord":	&discord.Service{},
	"pushover":	&pushover.Service{},
	"slack":	&slack.Service{},
	"teams":	&teams.Service{},
	"telegram":	&telegram.Service{},
	"smtp":	&smtp.Service{},
}

// Locate returns the service implementation that corresponds to the given scheme
func (router *ServiceRouter) Locate(serviceScheme string) (types.Service, error) {

	service, valid := services[strings.ToLower(serviceScheme)]
	if !valid {
		return nil, fmt.Errorf("unknown service scheme '%s'", serviceScheme)
	}

	return service, nil
}
