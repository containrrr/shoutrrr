package types

import "net/url"

// Enummer contains fields that have associated EnumFormatter instances
type Enummer interface {
	Enums() map[string]EnumFormatter
}

// ServiceConfig is the common interface for all types of service configurations
type ServiceConfig interface {
	Enummer
	GetURL() *url.URL
	SetURL(*url.URL) error
}

// ConfigQueryResolver is the interface used to get/set and list service config query fields
type ConfigQueryResolver interface {
	Get(string) (string, error)
	Set(string, string) error
	QueryFields() []string
}
