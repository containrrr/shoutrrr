package standard

import (
	"log"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/services"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Initializer is the standard implementation of Service.Initialize
type Initializer struct {
	logger *log.Logger
}

// Initialize sets the logger interface for the service to the specified logger and mutates config according to configURL
func (si *Initializer) Initialize(config types.ServiceConfig, configURL *url.URL, logger *log.Logger) error {
	if logger == nil {
		si.logger = services.DiscardLogger
	} else {
		si.logger = logger
	}

	if err := config.SetURL(configURL); err != nil {
		return err
	}

	return nil

}