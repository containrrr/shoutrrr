package router

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/services/discord"
	"github.com/containrrr/shoutrrr/pkg/services/ifttt"
	"github.com/containrrr/shoutrrr/pkg/services/pushover"
	"github.com/containrrr/shoutrrr/pkg/services/slack"
	"github.com/containrrr/shoutrrr/pkg/services/smtp"
	"github.com/containrrr/shoutrrr/pkg/services/teams"
	"github.com/containrrr/shoutrrr/pkg/services/telegram"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// ServiceRouter is responsible for routing a message to a specific notification service using the notification URL
type ServiceRouter struct {
	logger *log.Logger
}

// SetLogger sets the logger that the services will use to write progress logs
func (router *ServiceRouter) SetLogger(logger *log.Logger) {
	router.logger = logger
}

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

	service, err := router.Locate(rawURL)
	if err != nil {
		return err
	}

	return service.Send(message, nil)
}

var serviceMap = map[string]func() types.Service{
	"discord":  func() types.Service { return &discord.Service{} },
	"pushover": func() types.Service { return &pushover.Service{} },
	"slack":    func() types.Service { return &slack.Service{} },
	"teams":    func() types.Service { return &teams.Service{} },
	"telegram": func() types.Service { return &telegram.Service{} },
	"smtp":     func() types.Service { return &smtp.Service{} },
	"ifttt":    func() types.Service { return &ifttt.Service{} },
}

func (router *ServiceRouter) initService(rawURL string) (types.Service, error) {
	scheme, configURL, err := router.ExtractServiceName(rawURL)
	if err != nil {
		return nil, err
	}

	serviceFactory, valid := serviceMap[strings.ToLower(scheme)]
	if !valid {
		return nil, fmt.Errorf("unknown service scheme '%s'", scheme)
	}

	service := serviceFactory()

	err = service.Initialize(configURL, router.logger)
	if err != nil {
		return service, err
	}

	return service, nil
}

// Locate returns the service implementation that corresponds to the given service URL
func (router *ServiceRouter) Locate(rawURL string) (types.Service, error) {
	service, err := router.initService(rawURL)
	return service, err
}
