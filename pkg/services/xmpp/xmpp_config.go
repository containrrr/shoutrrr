package xmpp

import (
	"fmt"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util"
)

// Config is the configuration needed to send XMPP notifications
type Config struct {
	Host     string   `desc:"XMPP server hostname or IP address" url:"Host"`
	Username string   `desc:"XMPP username (local part or full JID)" url:"User"`
	Password string   `desc:"XMPP password" url:"Pass"`
	Port     uint16   `desc:"XMPP server port" default:"5222" key:"port"`
	Receiver []string `desc:"List of receiver JIDs separated by \",\" (comma)" key:"receiver,to"`
	StartTLS bool     `desc:"Use STARTTLS to negotiate TLS encryption" default:"Yes" key:"starttls"`
	TLS      bool     `desc:"Use implicit TLS (connects directly over TLS)" default:"No" key:"tls"`
}

// Enums implements types.ServiceConfig
func (*Config) Enums() map[string]types.EnumFormatter {
	return map[string]types.EnumFormatter{}
}

// GetURL returns a URL representation of its current field values
func (config *Config) GetURL() *url.URL {
	resolver := format.NewPropKeyResolver(config)
	return config.getURL(&resolver)
}

// SetURL updates a ServiceConfig from a URL representation of its field values
func (config *Config) SetURL(url *url.URL) error {
	resolver := format.NewPropKeyResolver(config)
	return config.setURL(&resolver, url)
}

func (config *Config) getURL(resolver types.ConfigQueryResolver) *url.URL {
	return &url.URL{
		User:       util.URLUserPassword(config.Username, config.Password),
		Host:       fmt.Sprintf("%s:%d", config.Host, config.Port),
		Scheme:     Scheme,
		ForceQuery: true,
		RawQuery:   format.BuildQuery(resolver),
	}
}

func (config *Config) setURL(resolver types.ConfigQueryResolver, url *url.URL) error {
	password, _ := url.User.Password()
	config.Username = url.User.Username()
	config.Password = password
	config.Host = url.Hostname()

	if url.Port() != "" {
		if _, err := fmt.Sscanf(url.Port(), "%d", &config.Port); err != nil {
			return fmt.Errorf("invalid port: %w", err)
		}
	}

	for key, vals := range url.Query() {
		if err := resolver.Set(key, vals[0]); err != nil {
			return err
		}
	}

	if !hasNonEmptyReceiver(config.Receiver) {
		return fmt.Errorf("receiver missing from config URL")
	}

	return nil
}

// hasNonEmptyReceiver returns true if the list has at least one non-empty receiver JID
func hasNonEmptyReceiver(receivers []string) bool {
	for _, r := range receivers {
		if r != "" {
			return true
		}
	}
	return false
}

// Scheme is the identifying part of this service's configuration URL
const Scheme = "xmpp"
