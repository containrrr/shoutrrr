package xmpp

import (
	"errors"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/types"
	"gosrc.io/xmpp"
	"net/url"
	"strconv"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
)

const (
	// Scheme is the identifying part of this service's configuration URL
	Scheme string = "xmpp"
)

// Config is the configuration needed to send notifications over XMPP
type Config struct {
	standard.EnumlessConfig
	Port       uint16
	Username   string
	Password   string
	Host       string
	ServerHost string `key:"serverhost"`
	ToAddress  string `key:"toaddress"`
	Subject    string `key:"subject"`
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
		User:     url.UserPassword(config.Username, config.Password),
		Host:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Path:     "/",
		Scheme:   Scheme,
		RawQuery: format.BuildQuery(resolver),
	}

}

func (config *Config) setURL(resolver types.ConfigQueryResolver, url *url.URL) error {

	password, _ := url.User.Password()

	config.Username = url.User.Username()
	config.Password = password
	config.Host = url.Hostname()

	if port, err := strconv.ParseUint(url.Port(), 10, 16); err == nil {
		config.Port = uint16(port)
	}

	for key, vals := range url.Query() {
		if err := resolver.Set(key, vals[0]); err != nil {
			return err
		}
	}

	if len(config.ToAddress) < 1 {
		return errors.New("toAddress missing from config URL")
	}

	return nil
}

func (config *Config) getClientConfig() *xmpp.Config {
	conf := xmpp.Config{
		Jid:      config.fromAddress(),
		Password: config.Password,
		Insecure: true,
	}

	if config.ServerHost != "" {
		conf.Address = fmt.Sprintf("%s:%d", config.ServerHost, config.Port)
	}

	return &conf
}

func (config *Config) fromAddress() string {
	return fmt.Sprintf("%s@%s", config.Username, config.Host)
}
