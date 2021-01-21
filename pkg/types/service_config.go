package types

import "net/url"

// ServiceConfig is the common interface for all types of service configurations
type ServiceConfig interface {
	GetURL() *url.URL
	SetURL(*url.URL) error
	Enums() map[string]EnumFormatter
}

// ConfigQueryResolver is the interface used to get/set and list service config query fields
type ConfigQueryResolver interface {
	Get(string) (string, error)
	Set(string, string) error
	QueryFields() []string
}
