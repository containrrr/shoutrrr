package shoutrrr

import (
	"github.com/containrrr/shoutrrr/internal/meta"
	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/containrrr/shoutrrr/pkg/types"
)

var defaultRouter = router.ServiceRouter{}

// SetLogger sets the logger that the services will use to write progress logs
func SetLogger(logger types.StdLogger) {
	defaultRouter.SetLogger(logger)
}

// Send notifications using a supplied url and message
func Send(rawURL string, message string) error {
	service, err := defaultRouter.Locate(rawURL)
	if err != nil {
		return err
	}

	return service.Send(message, &types.Params{})
}

// CreateSender returns a notification sender configured according to the supplied URL
func CreateSender(rawURLs ...string) (*router.ServiceRouter, error) {
	return router.New(nil, rawURLs...)
}

// NewSender returns a notification sender, writing any log output to logger and configured
// to send to the services indicated by the supplied URLs
func NewSender(logger types.StdLogger, serviceURLs ...string) (*router.ServiceRouter, error) {
	return router.New(logger, serviceURLs...)
}

// Version returns the shoutrrr version
func Version() string {
	return meta.Version
}
