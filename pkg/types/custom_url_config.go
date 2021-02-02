package types

import "net/url"

// CustomURLService is the interface that needs to be implemented to support custom URLs in services
type CustomURLService interface {
	Service
	GetConfigURLFromCustom(customURL *url.URL) (serviceURL *url.URL, err error)
}
