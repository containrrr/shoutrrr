package ifttt

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
)

const (
	// Scheme is the identifying part of this service's configuration URL
	Scheme = "ifttt"
)

// Config is the configuration needed to send IFTTT notifications
type Config struct {
	standard.EnumlessConfig
	WebHookID         string
	Events            []string
	Value1            string
	Value2            string
	Value3            string
	UseMessageAsValue uint8 `desc:"" default:"2"`
}

// GetURL returns a URL representation of it's current field values
func (config *Config) GetURL() *url.URL {

	return &url.URL{
		Host:     config.WebHookID,
		Path:     "/",
		Scheme:   Scheme,
		RawQuery: format.BuildQuery(config),
	}

}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *Config) SetURL(url *url.URL) error {

	config.WebHookID = url.Hostname()

	for key, vals := range url.Query() {
		if err := config.Set(key, vals[0]); err != nil {
			return err
		}
	}

	if len(config.Events) < 1 {
		return errors.New("events missing from config URL")
	}

	if len(config.WebHookID) < 1 {
		return errors.New("webhook ID missing from config URL")
	}

	return nil
}

// QueryFields returns the fields that are part of the Query of the service URL
func (config *Config) QueryFields() []string {
	return []string{
		"events",
		"value1",
		"value2",
		"value3",
		"messagevalue",
	}
}

// Get returns the value of a Query field
func (config *Config) Get(key string) (string, error) {
	switch key {
	case "events":
		return strings.Join(config.Events, ","), nil
	case "value1":
		return config.Value1, nil
	case "value2":
		return config.Value2, nil
	case "value3":
		return config.Value3, nil
	case "messagevalue":
		return fmt.Sprintf("%d", config.UseMessageAsValue), nil
	}
	return "", fmt.Errorf("invalid query key \"%s\"", key)
}

// Set updates the value of a Query field
func (config *Config) Set(key string, value string) error {
	switch key {
	case "events":
		config.Events = strings.Split(value, ",")
	case "value1":
		config.Value1 = value
	case "value2":
		config.Value2 = value
	case "value3":
		config.Value3 = value
	case "messagevalue":
		val64, err := strconv.ParseUint(value, 10, 8)
		if err == nil && val64 > 3 {
			err = errors.New("only values 1-3 are supported")
		}
		if err != nil {
			return fmt.Errorf("invalid value \"%s\" for \"messagevalue\": %s", value, err)
		}
		config.UseMessageAsValue = uint8(val64)
	default:
		return fmt.Errorf("invalid query key \"%s\"", key)
	}
	return nil
}
