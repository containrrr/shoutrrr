package types

import (
	"log"
	"net/url"
)


// Service is the common interface for all notification services
type Service interface {
	Send(message string, params *map[string]string) error
	NewConfig() ServiceConfig
	Initialize(config ServiceConfig, serviceURL *url.URL, logger *log.Logger) error
	ApplyTemplate(template string, params *map[string]string) (string, error)
}