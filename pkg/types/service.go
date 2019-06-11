package types

import (
	"log"
	"net/url"
)


// Service is the common interface for all notification services
type Service interface {
	Send(serviceURL *url.URL, message string, params *map[string]string) error
	GetConfig() ServiceConfig
	SetLogger(logger *log.Logger)
	ApplyTemplate(template string, params *map[string]string) (string, error)
}