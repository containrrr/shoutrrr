package shoutrrr

import (
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
	"log"
	"net/url"
)

type Service struct {
	standard.Standard
	config *Config
}

func (s *Service) Send(message string, params *types.Params) error {
	panic("implement me")
}

func (s *Service) Initialize(serviceURL *url.URL, logger *log.Logger) error {
	s.SetLogger(logger)
	s.config = &Config{}

	return nil
}

const (
	Schema = "shoutrrr"
)
