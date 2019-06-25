package shoutrrr

import (
	"log"

	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util/queue"
)

var routing = router.ServiceRouter{}

// SetLogger sets the logger that the services will use to write progress logs
func SetLogger(logger *log.Logger) {
	routing.SetLogger(logger)
}

// Send lets you send shoutrrr notifications using a supplied url and message
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

// CreateQueue returns a notification queued sender configured according to the supplied URL
func CreateQueue(rawURL string) (types.QueuedSender, error) {
	service, err := routing.Locate(rawURL)
	if err != nil {
		return nil, err
	}

	return queue.GetQueued(service), nil
}

// VerifyURL lets you verify that a configuration URL is valid
func VerifyURL(rawURL string) error {
	_, err := routing.Locate(rawURL)
	return err
}