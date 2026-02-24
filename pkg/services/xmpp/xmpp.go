// Package xmpp implements XMPP (Jabber) as a shoutrrr service
package xmpp

import (
	"fmt"
	"net/url"

	xmppclient "github.com/xmppo/go-xmpp"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Service sends notifications via XMPP (Jabber)
type Service struct {
	standard.Standard
	config *Config
	pkr    format.PropKeyResolver
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{}
	service.pkr = format.NewPropKeyResolver(service.config)
	if err := service.pkr.SetDefaultProps(service.config); err != nil {
		return err
	}
	return service.config.setURL(&service.pkr, configURL)
}

// Send a notification message to XMPP recipients
func (service *Service) Send(message string, params *types.Params) error {
	config := service.config

	if err := service.pkr.UpdateConfigFromParams(config, params); err != nil {
		return fmt.Errorf("failed to update config from params: %w", err)
	}

	client, err := dial(config)
	if err != nil {
		return fmt.Errorf("failed to connect to XMPP server: %w", err)
	}
	defer client.Close()

	for _, receiver := range config.Receiver {
		if receiver == "" {
			continue
		}
		if _, err := client.Send(xmppclient.Chat{
			Remote: receiver,
			Type:   "chat",
			Text:   message,
		}); err != nil {
			return fmt.Errorf("failed to send XMPP message to %q: %w", receiver, err)
		}
	}

	return nil
}

func dial(config *Config) (*xmppclient.Client, error) {
	host := fmt.Sprintf("%s:%d", config.Host, config.Port)
	options := xmppclient.Options{
		Host:     host,
		User:     config.Username,
		Password: config.Password,
		NoTLS:    !config.TLS,
		StartTLS: config.StartTLS,
	}

	client, err := options.NewClient()
	if err != nil {
		return nil, err
	}

	return client, nil
}
