//+build xmpp

package xmpp

import (
	"net/url"

	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"

	"github.com/containrrr/shoutrrr/pkg/format"
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
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{
		Port:    5222,
		Subject: "Shoutrrr Notification",
	}

	service.pkr = format.NewPropKeyResolver(service.config)
	if err := service.config.setURL(&service.pkr, configURL); err != nil {
		return err
	}

	service.router = xmpp.NewRouter()
	service.router.HandleFunc("message", func(s xmpp.Sender, p stanza.Packet) {
		msg, ok := p.(stanza.Message)
		if !ok {
			service.Logf("XMPP: ignoring unknown packet: %T", p)
			return
		}

		service.Logf("XMPP: message from %s: %s", msg.From, msg.Body)
	})

	config := service.config.getClientConfig()

	client, err := xmpp.NewClient(config, service.router, func(err error) {
		service.Log("XMPP:", err)
	})
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
