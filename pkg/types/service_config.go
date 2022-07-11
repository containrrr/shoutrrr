package types

import (
	"net/url"
	"strings"
)

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

// GeneratedConfig is the interface for service configs created using shoutrr-gen
type GeneratedConfig interface {
	ServiceConfig
	PropInfo() *ConfigPropInfo
	Update(map[int]string) error
	PropValue(int) string
}

type ConfigWithLegacyURLSupport interface {
	UpdateLegacyURL(*url.URL) *url.URL
}

type CustomQueryConfig interface {
	CustomQueryVars() url.Values
}

type ConfigPropInfo struct {
	PropNames      []string
	DefaultValues  []string
	Keys           []string
	PrimaryKeys    []int
	KeyPropIndexes map[string]int
}

func (cpr *ConfigPropInfo) PropIndexFor(key string) (int, bool) {
	val, found := cpr.KeyPropIndexes[strings.ToLower(key)]
	return val, found
}
