package standard

import (
"errors"
)

// QuerylessConfig implements the ServiceConfig interface for services that does not use Query fields
type QuerylessConfig struct {}


// QueryFields returns an empty list of Query fields
func (qc *QuerylessConfig) QueryFields() []string {
	return []string{}
}

// Get is a dummy function that will return an error if called
func (qc *QuerylessConfig) Get(string) (string, error) {
	return "", errors.New("service config does not support Get")
}

// Set is a dummy function that will return an error if called
func (qc *QuerylessConfig) Set(string, string) error {
	return errors.New("service config does not support Set")
}