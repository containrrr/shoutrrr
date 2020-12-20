package ifttt

import (
	"errors"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
	"net/url"
)

const (
	// Scheme is the identifying part of this service's configuration URL
	Scheme = "ifttt"
)

// Config is the configuration needed to send IFTTT notifications
type Config struct {
	standard.EnumlessConfig
	WebHookID         string
	Events            []string `key:"events"`
	Value1            string   `key:"value1"`
	Value2            string   `key:"value2"`
	Value3            string   `key:"value3"`
	UseMessageAsValue uint8    `key:"messagevalue" desc:"" default:"2"`
}

// GetURL returns a URL representation of it's current field values
func (config *Config) GetURL() *url.URL {
	resolver := format.NewPropKeyResolver(config)
	return config.getURL(&resolver)
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *Config) SetURL(url *url.URL) error {
	resolver := format.NewPropKeyResolver(config)
	return config.setURL(&resolver, url)
}

func (config *Config) getURL(resolver types.ConfigQueryResolver) *url.URL {

	return &url.URL{
		Host:     config.WebHookID,
		Path:     "/",
		Scheme:   Scheme,
		RawQuery: format.BuildQuery(resolver),
	}
}

func (config *Config) setURL(resolver types.ConfigQueryResolver, url *url.URL) error {
	if config.UseMessageAsValue == 0 {
		config.UseMessageAsValue = 2
	}
	config.WebHookID = url.Hostname()

	for key, vals := range url.Query() {
		if err := resolver.Set(key, vals[0]); err != nil {
			return err
		}
	}

	if config.UseMessageAsValue > 3 || config.UseMessageAsValue < 1 {
		return errors.New("invalid value for messagevalue: only values 1-3 are supported")
	}

	if len(config.Events) < 1 {
		return errors.New("events missing from config URL")
	}

	if len(config.WebHookID) < 1 {
		return errors.New("webhook ID missing from config URL")
	}

	return nil
}
