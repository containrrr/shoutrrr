package types

import (
	"net/http"
	"net/url"
)

// Service is the public common interface for all notification services
type Service interface {
	Sender
	Templater
	Initialize(serviceURL *url.URL, logger StdLogger) error
	SetLogger(logger StdLogger)
}

// HTTPService is the common interface for services that use a http.Client to send notifications
type HTTPService interface {
	HTTPClient() *http.Client
}
