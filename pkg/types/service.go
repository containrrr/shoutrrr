package types

import (
	"log"
	"net/url"
)

// Service is the public common interface for all notification services
type Service interface {
	Sender
	Templater
	Initialize(serviceURL *url.URL, logger *log.Logger) error
}