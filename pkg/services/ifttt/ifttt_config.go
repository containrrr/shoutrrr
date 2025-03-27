package ifttt

import (
	"errors"
	"net/url"

	"github.com/nicholas-fedor/shoutrrr/pkg/format"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/standard"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
)

const (
	// Scheme is the identifying part of this service's configuration URL.
	Scheme = "ifttt"

	// Config constants for IFTTT webhook values.
	DefaultMessageValue = 2 // Default value field (1-3) for the notification message.
	DisabledValue       = 0 // Value to disable title assignment.
	MinValueField       = 1 // Minimum valid value field (Value1).
	MaxValueField       = 3 // Maximum valid value field (Value3).
	MinLength           = 1 // Minimum length for required fields like Events and WebHookID.
)

// Config is the configuration needed to send IFTTT notifications.
type Config struct {
	standard.EnumlessConfig
	WebHookID         string   `required:"true" url:"host"`
	Events            []string `key:"events"    required:"true"`
	Value1            string   `key:"value1"    optional:""`
	Value2            string   `key:"value2"    optional:""`
	Value3            string   `key:"value3"    optional:""`
	UseMessageAsValue uint8    `default:"2"     desc:"Sets the corresponding value field to the notification message" key:"messagevalue"`
	UseTitleAsValue   uint8    `default:"0"     desc:"Sets the corresponding value field to the notification title"   key:"titlevalue"`
	Title             string   `default:""      desc:"Notification title, optionally set by the sender"               key:"title"`
}

// GetURL returns a URL representation of its current field values.
func (config *Config) GetURL() *url.URL {
	resolver := format.NewPropKeyResolver(config)

	return config.getURL(&resolver)
}

// SetURL updates a ServiceConfig from a URL representation of its field values.
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
	if config.UseMessageAsValue == DisabledValue {
		config.UseMessageAsValue = DefaultMessageValue
	}

	config.WebHookID = url.Hostname()

	for key, vals := range url.Query() {
		if err := resolver.Set(key, vals[0]); err != nil {
			return err
		}
	}

	if config.UseMessageAsValue > MaxValueField || config.UseMessageAsValue < MinValueField {
		return errors.New("invalid value for messagevalue: only values 1-3 are supported")
	}

	if config.UseTitleAsValue > MaxValueField {
		return errors.New("invalid value for titlevalue: only values 1-3 or 0 (for disabling) are supported")
	}

	if config.UseTitleAsValue != DisabledValue && config.UseTitleAsValue == config.UseMessageAsValue {
		return errors.New("titlevalue cannot use the same number as messagevalue")
	}

	if url.String() != "ifttt://dummy@dummy.com" {
		if len(config.Events) < MinLength {
			return errors.New("events missing from config URL")
		}

		if len(config.WebHookID) < MinLength {
			return errors.New("webhook ID missing from config URL")
		}
	}

	return nil
}
