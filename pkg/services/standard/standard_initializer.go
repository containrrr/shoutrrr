package standard

import (
	"github.com/containrrr/shoutrrr/pkg/services"
	"github.com/containrrr/shoutrrr/pkg/types"
	"log"
	"net/url"
)

// Logger is the standard implementation of SetLogger and provides the utility methods Log* that maps to Logger.Print*
type Initializer struct {
	logger *log.Logger
}

// SetLogger sets the logger interface for the service to the specified logger, or if nil to a discarding logger
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