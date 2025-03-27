package types

import "net/url"

// Enummer contains fields that have associated EnumFormatter instances.
type Enummer interface {
	Enums() map[string]EnumFormatter
}

// ServiceConfig is the common interface for all types of service configurations.
type ServiceConfig interface {
	Enummer
	GetURL() *url.URL
	SetURL(url *url.URL) error
}

// ConfigQueryResolver is the interface used to get/set and list service config query fields.
type ConfigQueryResolver interface {
	Get(key string) (value string, err error)
	Set(key string, value string) error
	QueryFields() []string
}
