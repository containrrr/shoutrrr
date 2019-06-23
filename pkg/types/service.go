package types

import (
	"log"
	"net/url"
)

// Service is the public common interface for all notification services
type Service interface {
	Initialize(serviceURL *url.URL, logger *log.Logger) error
	Send(message string, params *map[string]string) error

	// Queue methods
	Enqueuef(format string, v ...interface{})
	Enqueue(message string)
	Flush(params *map[string]string)
}