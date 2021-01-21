package router

import (
	"fmt"
	"log"
	"net/url"
	"reflect"
	"strings"
	"time"

	t "github.com/containrrr/shoutrrr/pkg/types"
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
	results := router.SendAsync(message, params)

	for i := range router.services {
		errors[i] = <-results
	}

	return errors
}

// SendAsync sends the specified message using the routers underlying services
func (router *ServiceRouter) SendAsync(message string, params *t.Params) chan error {
	serviceCount := len(router.services)
	proxy := make(chan error, serviceCount)
	errors := make(chan error, serviceCount)

	if params == nil {
		params = &t.Params{}
	}
	for _, service := range router.services {
		go sendToService(service, proxy, router.Timeout, message, *params)
	}

	go func() {
		for i := 0; i < serviceCount; i++ {
			errors <- <-proxy
		}
		close(errors)
	}()

	return errors
}

func sendToService(service t.Service, results chan error, timeout time.Duration, message string, params t.Params) {
	result := make(chan error)

	// TODO: There really ought to be a better way to name the services
	pkg := reflect.TypeOf(service).Elem().PkgPath()
	serviceName := pkg[strings.LastIndex(pkg, "/")+1:]

	go func() { result <- service.Send(message, &params) }()

	select {
	case res := <-result:
		results <- res
	case <-time.After(timeout):
		results <- fmt.Errorf("failed to send using %v: timed out", serviceName)
	}
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

	scheme := serviceURL.Scheme
	schemeParts := strings.Split(scheme, "+")

	if len(schemeParts) > 1 {
		scheme = schemeParts[0]
	}

	return scheme, serviceURL, nil
}

// Route a message to a specific notification service using the notification URL
func (router *ServiceRouter) Route(rawURL string, message string) error {

	service, err := router.Locate(rawURL)
	if err != nil {
		return err
	}

	return service.Send(message, nil)
}

func (router *ServiceRouter) initService(rawURL string) (t.Service, error) {

	scheme, configURL, err := router.ExtractServiceName(rawURL)
	if err != nil {
		return nil, err
	}

	serviceFactory, valid := serviceMap[strings.ToLower(scheme)]
	if !valid {
		return nil, fmt.Errorf("unknown service scheme for URL '%s'", rawURL)
	}

	service := serviceFactory()

	if configURL.Scheme != scheme {
		router.logger.Println("Got custom URL:", configURL.String())
		customURLService, ok := service.(t.CustomURLService)
		if !ok {
			return nil, fmt.Errorf("custom URLs are not supported by '%s' service", scheme)
		}
		configURL, err = customURLService.GetConfigURLFromCustom(configURL)
		if err != nil {
			return nil, err
		}
		router.logger.Println("Converted service URL:", configURL.String())
	}

	err = service.Initialize(configURL, router.logger)
	if err != nil {
		return service, err
	}

	return service, nil
}

// NewService returns a new uninitialized service instance
func (router *ServiceRouter) NewService(service string) (t.Service, error) {
	serviceFactory, valid := serviceMap[strings.ToLower(service)]
	if !valid {
		return nil, fmt.Errorf("unknown service %q", service)
	}
	return serviceFactory(), nil
}

// ListServices returns the available services
func (router *ServiceRouter) ListServices() []string {
	services := make([]string, len(serviceMap))

	i := 0
	for key := range serviceMap {
		services[i] = key
		i++
	}

	return services
}

// Locate returns the service implementation that corresponds to the given service URL
func (router *ServiceRouter) Locate(rawURL string) (t.Service, error) {
	service, err := router.initService(rawURL)
	return service, err
}
