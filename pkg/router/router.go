package router

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	t "github.com/containrrr/shoutrrr/pkg/types"
)

// ServiceRouter is responsible for routing a message to a specific notification service using the notification URL
type ServiceRouter struct {
	logger      t.StdLogger
	services    []t.Service
	serviceIds  []string
	schemeCount map[string]int
	queue       []string
	Timeout     time.Duration
}

// New creates a new service router using the specified logger and service URLs
func New(logger t.StdLogger, serviceURLs ...string) (*ServiceRouter, error) {
	router := ServiceRouter{
		logger:      logger,
		schemeCount: map[string]int{},
		Timeout:     10 * time.Second,
	}

	for _, serviceURL := range serviceURLs {
		if err := router.AddService(serviceURL); err != nil {
			return nil, fmt.Errorf("error initializing router services: %s", err)
		}
	}
	return &router, nil
}

// AddService initializes the specified service from its URL, and adds it if no errors occur
func (router *ServiceRouter) AddService(serviceURL string) error {
	service, scheme, err := router.initService(serviceURL)
	if err == nil {
		router.services = append(router.services, service)
		router.serviceIds = append(router.serviceIds, router.getNextServiceId(scheme))
	}
	return err
}

func (router *ServiceRouter) getNextServiceId(scheme string) string {
	if router.schemeCount == nil {
		router.schemeCount = map[string]int{}
	}
	schemeIndex := router.schemeCount[scheme] + 1
	router.schemeCount[scheme] = schemeIndex

	// Only append a sequence number for the second and subsequent instances
	if schemeIndex < 2 {
		return scheme
	}

	return scheme + fmt.Sprintf("%v", schemeIndex)
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

// SendItems sends the specified message items using the routers underlying services
func (router *ServiceRouter) SendItems(items []t.MessageItem, params t.Params) []error {
	if router == nil {
		return []error{fmt.Errorf("error sending message: no senders")}
	}

	// Fallback using old API for now
	message := strings.Builder{}
	for _, item := range items {
		message.WriteString(item.Text)
	}

	serviceCount := len(router.services)
	errors := make([]error, serviceCount)
	results := router.SendAsync(message.String(), &params)

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
	for i, service := range router.services {
		go sendToService(service, proxy, router.Timeout, message, *params, router.serviceIds[i])
	}

	go func() {
		for i := 0; i < serviceCount; i++ {
			errors <- <-proxy
		}
		close(errors)
	}()

	return errors
}

func sendToService(service t.Service, results chan error, timeout time.Duration, message string, params t.Params, serviceId string) {
	result := make(chan error)

	go func() {
		err := service.Send(message, &params)
		if err != nil {
			err = fmt.Errorf("%v: %v", serviceId, err)
		}
		result <- err
	}()

	select {
	case res := <-result:
		results <- res
	case <-time.After(timeout):
		results <- fmt.Errorf("failed to send using %v: timed out", serviceId)
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
func (router *ServiceRouter) SetLogger(logger t.StdLogger) {
	router.logger = logger
	for _, service := range router.services {
		service.SetLogger(logger)
	}
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

func (router *ServiceRouter) initService(rawURL string) (service t.Service, scheme string, err error) {

	var configURL *url.URL
	scheme, configURL, err = router.ExtractServiceName(rawURL)
	if err != nil {
		return
	}

	service, err = newService(scheme)
	if err != nil {
		return
	}

	if configURL.Scheme != scheme {
		router.log("Got custom URL:", configURL.String())
		customURLService, ok := service.(t.CustomURLService)
		if !ok {
			return nil, scheme, fmt.Errorf("custom URLs are not supported by '%s' service", scheme)
		}
		configURL, err = customURLService.GetConfigURLFromCustom(configURL)
		if err != nil {
			return
		}
		router.log("Converted service URL:", configURL.String())
	}

	err = service.Initialize(configURL, router.logger)
	return
}

// NewService returns a new uninitialized service instance
func (*ServiceRouter) NewService(serviceScheme string) (t.Service, error) {
	return newService(serviceScheme)
}

// newService returns a new uninitialized service instance
func newService(serviceScheme string) (t.Service, error) {
	serviceFactory, valid := serviceMap[strings.ToLower(serviceScheme)]
	if !valid {
		return nil, fmt.Errorf("unknown service %q", serviceScheme)
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

// ListAddedServices returns a list of the scheme identifiers of the added services
func (router *ServiceRouter) ListAddedServices() []string {
	return router.serviceIds
}

// Locate returns the service implementation that corresponds to the given service URL
func (router *ServiceRouter) Locate(rawURL string) (t.Service, error) {
	service, _, err := router.initService(rawURL)
	return service, err
}

func (router *ServiceRouter) log(v ...interface{}) {
	if router.logger == nil {
		return
	}
	router.logger.Println(v...)
}
