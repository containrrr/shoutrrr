package types

import (
	"net/url"
)


// Service is the common interface for all notification services
type Service interface {
	Send(serviceURL *url.URL, message string, opts ServiceOpts) error
	GetConfig() ServiceConfig
}