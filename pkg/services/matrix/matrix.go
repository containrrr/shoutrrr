package matrix

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	t "github.com/containrrr/shoutrrr/pkg/types"
	"log"
	"net/url"
)

const Scheme = "matrix"

type Service struct {
	standard.Standard
	config *Config
	client *Client
	pkr    format.PropKeyResolver
}

func (s *Service) Initialize(configURL *url.URL, logger *log.Logger) error {
	s.SetLogger(logger)
	s.config = &Config{}

	s.pkr = format.NewPropKeyResolver(s.config)
	if err := s.config.setURL(&s.pkr, configURL); err != nil {
		return err
	}

	s.client = NewClient(s.config.Host, s.config.DisableTLS, logger)
	if s.config.User != "" {
		return s.client.Login(s.config.User, s.config.Password)
	}

	s.client.UseToken(s.config.Password)
	return nil
}

// Send notification
func (s *Service) Send(message string, params *t.Params) error {
	config := *s.config
	if err := s.pkr.UpdateConfigFromParams(&config, params); err != nil {
		return err
	}

	errors := s.client.SendMessage(message, s.config.Rooms)

	if len(errors) > 0 {
		for _, err := range errors {
			s.Logf("error sending message: %v", err)
		}
		return fmt.Errorf("%v error(s) sending message, with initial error: %v", len(errors), errors[0])
	}

	return nil
}
