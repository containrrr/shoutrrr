package matrix

import (
	"fmt"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/conf"
	"github.com/containrrr/shoutrrr/pkg/pkr"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	t "github.com/containrrr/shoutrrr/pkg/types"
)

// Scheme is the identifying part of this service's configuration URL
const Scheme = "matrix"

// Service providing Matrix as a notification service
type Service struct {
	standard.Standard
	config *Config
	client *client
	pkr    pkr.PropKeyResolver
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (s *Service) Initialize(configURL *url.URL, logger t.StdLogger) error {
	s.SetLogger(logger)
	s.config = &Config{}
	if err := conf.Init(s.config, configURL); err != nil {
		return err
	}

	s.client = newClient(s.config.Host, s.config.DisableTLS, logger)
	if s.config.User != "" {
		return s.client.login(s.config.User, s.config.Password)
	}

	s.client.useToken(s.config.Password)
	return nil
}

// Send notification
func (s *Service) Send(message string, params *t.Params) error {
	config := *s.config
	if err := conf.UpdateFromParams(&config, params); err != nil {
		return err
	}

	errors := s.client.sendMessage(message, s.config.Rooms)

	if len(errors) > 0 {
		for _, err := range errors {
			s.Logf("error sending message: %v", err)
		}
		return fmt.Errorf("%v error(s) sending message, with initial error: %v", len(errors), errors[0])
	}

	return nil
}
