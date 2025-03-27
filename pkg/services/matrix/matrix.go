package matrix

import (
	"fmt"
	"net/url"

	"github.com/nicholas-fedor/shoutrrr/pkg/format"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/standard"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
)

// Scheme is the identifying part of this service's configuration URL.
const Scheme = "matrix"

// Service providing Matrix as a notification service.
type Service struct {
	standard.Standard
	Config *Config
	client *client
	pkr    format.PropKeyResolver
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service.
func (s *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	s.SetLogger(logger)
	s.Config = &Config{}
	s.pkr = format.NewPropKeyResolver(s.Config)

	if err := s.Config.setURL(&s.pkr, configURL); err != nil {
		return err
	}

	if configURL.String() != "matrix://dummy@dummy.com" {
		s.client = newClient(s.Config.Host, s.Config.DisableTLS, logger)
		if s.Config.User != "" {
			return s.client.login(s.Config.User, s.Config.Password)
		}

		s.client.useToken(s.Config.Password)
	}

	return nil
}

// GetID returns the service identifier.
func (service *Service) GetID() string {
	return Scheme
}

// Send notification.
func (s *Service) Send(message string, params *types.Params) error {
	config := *s.Config
	if err := s.pkr.UpdateConfigFromParams(&config, params); err != nil {
		return err
	}

	if s.client == nil {
		return fmt.Errorf("client not initialized; cannot send message")
	}

	errors := s.client.sendMessage(message, s.Config.Rooms)
	if len(errors) > 0 {
		for _, err := range errors {
			s.Logf("error sending message: %w", err)
		}

		return fmt.Errorf("%v error(s) sending message, with initial error: %v", len(errors), errors[0])
	}

	return nil
}
