package services

import (
	"io/ioutil"
	"log"
)

// ServiceOpts contains various properties that can be set by the consumer to alter the behaviour of services
type ServiceOpts struct {
	verbose bool
	logger *log.Logger
	props map[string]string
}

// Verbose marks whether the consumer want a service to output to the logger
func (svc *ServiceOpts) Verbose() bool {
	return svc.verbose
}

// Logger is the logging interface to be used by a service
func (svc *ServiceOpts) Logger() *log.Logger {
	return svc.logger
}

// Props is a collection of strings that should be used for substitution in fields tagged with `template:"yes"`
func (svc *ServiceOpts) Props() map[string]string {
	return svc.props
}

// GetDefaultOpts creates a default ServiceOpts struct that discards any output written to it's logger
func GetDefaultOpts() *ServiceOpts {
	return &ServiceOpts{
		verbose: false,
		logger: DiscardLogger,
	}
}

// CreateServiceOpts creates a ServiceOpts struct
func CreateServiceOpts(logger *log.Logger, verbose bool, props map[string]string) *ServiceOpts {
	return &ServiceOpts{
		verbose,
		logger,
		props,
	}
}

// DiscardLogger is a logger that discards any output written to it
var DiscardLogger = log.New(ioutil.Discard, "", 0)