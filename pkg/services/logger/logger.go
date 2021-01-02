package logger

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Service is the Logger service struct
type Service struct {
	standard.Standard
	config *Config
}

// Send a notification message to log
func (service *Service) Send(message string, params *types.Params) error {
	if params == nil {
		params = &types.Params{}
	}
	(*params)["message"] = message
	return service.doSend(params)
}

func (service *Service) doSend(params *types.Params) error {
	msg := (*params)["message"]
	if tpl, found := service.GetTemplate("message"); found {
		wc := &strings.Builder{}
		if err := tpl.Execute(wc, params); err != nil {
			return fmt.Errorf("failed to write template to log: %s", err)
		}
		msg = wc.String()
	}
	service.Log(msg)
	return nil
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(_ *url.URL, logger *log.Logger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{}
	return nil
}
