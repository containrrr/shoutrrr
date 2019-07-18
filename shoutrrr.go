package shoutrrr

import (
	"log"

	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/containrrr/shoutrrr/pkg/types"
)

var routing = router.ServiceRouter{}

// SetLogger sets the logger that the services will use to write progress logs
func SetLogger(logger *log.Logger) {
	routing.SetLogger(logger)
}

// Send notifications using a supplied url and message
func Send(rawURL string, message string) error {
	service, err := routing.Locate(rawURL)
	if err != nil {
		return err
	}

	return service.Send(message, nil)
}

// CreateSender returns a notification sender configured according to the supplied URL
func CreateSender(rawURL string) (types.Service, error) {
	return routing.Locate(rawURL)
}