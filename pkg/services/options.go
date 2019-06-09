package services

import (
	"io/ioutil"
	"log"
)

type ServiceOpts struct {
	verbose bool
	logger *log.Logger
	props map[string]string
}

func (svc *ServiceOpts) Verbose() bool {
	return svc.verbose
}

func (svc *ServiceOpts) Logger() *log.Logger {
	return svc.logger
}

func (svc *ServiceOpts) Props() map[string]string {
	return svc.props
}

func GetDefaultOpts() *ServiceOpts {
	return &ServiceOpts{
		verbose: false,
		logger: DiscardLogger,
	}
}

func CreateServiceOpts(logger *log.Logger, verbose bool, props map[string]string) *ServiceOpts {
	return &ServiceOpts{
		verbose,
		logger,
		props,
	}
}

var DiscardLogger = log.New(ioutil.Discard, "", 0)