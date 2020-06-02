package router

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/containrrr/shoutrrr/pkg/services/discord"
	"github.com/containrrr/shoutrrr/pkg/services/gotify"
	"github.com/containrrr/shoutrrr/pkg/services/hangouts"
	"github.com/containrrr/shoutrrr/pkg/services/ifttt"
	"github.com/containrrr/shoutrrr/pkg/services/logger"
	"github.com/containrrr/shoutrrr/pkg/services/mattermost"
	"github.com/containrrr/shoutrrr/pkg/services/pushbullet"
	"github.com/containrrr/shoutrrr/pkg/services/pushover"
	"github.com/containrrr/shoutrrr/pkg/services/slack"
	"github.com/containrrr/shoutrrr/pkg/services/smtp"
	"github.com/containrrr/shoutrrr/pkg/services/teams"
	"github.com/containrrr/shoutrrr/pkg/services/telegram"
	"github.com/containrrr/shoutrrr/pkg/services/zulip"
	t "github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/xmpp"
)

// ServiceRouter is responsible for routing a message to a specific notification service using the notification URL
type ServiceRouter struct {
	logger   *log.Logger
	services []t.Service
	queue    []string
	Timeout  time.Duration
}

// New creates a new service router using the specified logger and service URLs
func New(logger *log.Logger, serviceURLs ...string) (*ServiceRouter, error) {
	router := ServiceRouter{
		logger:  logger,
		Timeout: 10 * time.Second,
	}
	for _, serviceURL := range serviceURLs {
		service, err := router.initService(serviceURL)
		if err != nil {
			return nil, fmt.Errorf("error initializing router services: %s", err)
		}
		router.services = append(router.services, service)
	}
	return &router, nil
}

// Send sends the specified message using the routers underlying services
func (router *ServiceRouter) Send(message string, params *t.Params) []error {
	if router == nil {
		return []error{fmt.Errorf("error sending message: no senders")}
	}

	serviceCount := len(router.services)
	errors := make([]error, serviceCount)
	results := make(chan error, serviceCount)

	if params == nil {
		params = &t.Params{}
	}
	for _, service := range router.services {
		go sendToService(service, results, router.Timeout, message, *params)
	}
	for i := range router.services {
		select {
		case res := <-results:
			errors[i] = res
		case <-time.After(10 * time.Second):
			fmt.Println("timeout 1")
		}
	}
	return errors
}

func sendToService(service t.Service, results chan error, timeout time.Duration, message string, params t.Params) {
	// TODO: There really ought to be a way to tell what service generated the error
	result := make(chan error, 1)

	go func() { result <- service.Send(message, &params) }()

	select {
	case res := <-result:
		results <- res
	case <-time.After(timeout):
		results <- fmt.Errorf("timed out")
	}
	close(result)
}

// Enqueue adds the message to an internal queue and sends it when Flush is invoked
func (router *ServiceRouter) Enqueue(message string, v ...interface{}) {
	if len(v) > 0 {
		message = fmt.Sprintf(message, v...)
	}
	router.queue = append(router.queue, message)
}

// Flush sends all messages that have been queued up as a combined message. This method should be deferred!
func (router *ServiceRouter) Flush(params *t.Params) {
	// Since this method is supposed to be deferred we just have to ignore errors
	_ = router.Send(strings.Join(router.queue, "\n"), params)
	router.queue = []string{}
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
func (router *ServiceRouter) Route(rawURL string, message string) error {

	service, err := router.Locate(rawURL)
	if err != nil {
		return err
	}

	return service.Send(message, nil)
}

var serviceMap = map[string]func() t.Service{
	"discord":    func() t.Service { return &discord.Service{} },
	"pushover":   func() t.Service { return &pushover.Service{} },
	"slack":      func() t.Service { return &slack.Service{} },
	"teams":      func() t.Service { return &teams.Service{} },
	"telegram":   func() t.Service { return &telegram.Service{} },
	"smtp":       func() t.Service { return &smtp.Service{} },
	"ifttt":      func() t.Service { return &ifttt.Service{} },
	"gotify":     func() t.Service { return &gotify.Service{} },
	"logger":     func() t.Service { return &logger.Service{} },
	"xmpp":       func() t.Service { return &xmpp.Service{} },
	"pushbullet": func() t.Service { return &pushbullet.Service{} },
	"mattermost": func() t.Service { return &mattermost.Service{} },
	"hangouts":   func() t.Service { return &hangouts.Service{} },
	"zulip":      func() t.Service { return &zulip.Service{} },
}

func (router *ServiceRouter) initService(rawURL string) (t.Service, error) {
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
func (router *ServiceRouter) Locate(rawURL string) (t.Service, error) {
	service, err := router.initService(rawURL)
	return service, err
}
