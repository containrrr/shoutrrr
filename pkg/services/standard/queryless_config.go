package standard

import (
	"errors"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// QuerylessConfig implements the ServiceConfig interface for services that does not use Query fields
type QuerylessConfig struct {}


// QueryFields returns an empty list of Query fields
func (qc *QuerylessConfig) QueryFields() []string {
	return []string{}
}

// Enums returns an empty map
func (qc *QuerylessConfig) Enums() map[string]types.EnumFormatter {
	return map[string]types.EnumFormatter{}
}

// Get is a dummy function that will return an error if called
func (qc *QuerylessConfig) Get(string) (string, error) {
	return "", errors.New("service config does not support Get")
}

// Set is a dummy function that will return an error if called
func (qc *QuerylessConfig) Set(string, string) error {
	return errors.New("service config does not support Get")
}