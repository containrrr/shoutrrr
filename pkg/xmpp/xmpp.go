package xmpp

import (
	"github.com/containrrr/shoutrrr/pkg/format"
	"log"
	"net/url"

	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Service sends notifications via XMPP
type Service struct {
	standard.Standard
	pkr    format.PropKeyResolver
	client *xmpp.Client
	router *xmpp.Router
	config *Config
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger *log.Logger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{
		Port:    5222,
		Subject: "Shoutrrr Notification",
	}

	service.pkr = format.NewPropKeyResolver(service.config)
	if err := service.config.setURL(&service.pkr, configURL); err != nil {
		return err
	}

	config := service.config.getClientConfig()

	client, err := xmpp.NewClient(*config, service.router)
	if err != nil {
		return err
	}

	service.client = client

	return nil
}

// Send a notification message to the configured recipient
func (service *Service) Send(message string, params *types.Params) error {
	if err := service.client.Connect(); err != nil {
		return err
	}

	config := service.config
	if err := service.pkr.UpdateConfigFromParams(config, params); err != nil {
		return err
	}

	msg := stanza.Message{
		Subject: config.Subject,
		Body:    message,
		Attrs: stanza.Attrs{
			To: service.config.ToAddress,
		},
	}

	return service.client.Send(msg)

}
