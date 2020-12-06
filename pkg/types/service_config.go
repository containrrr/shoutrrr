package types

import "net/url"

// ServiceConfig is the common interface for all types of service configurations
type ServiceConfig interface {
	GetURL() *url.URL
	SetURL(*url.URL) error
	Enums() map[string]EnumFormatter
}

type ConfigQueryResolver interface {
	Get(string) (string, error)
	Set(string, string) error
	QueryFields() []string
}
