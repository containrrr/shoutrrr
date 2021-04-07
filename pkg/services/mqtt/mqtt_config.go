package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strconv"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Config for use within the mqtt
type Config struct {
	Host       string `key:"host" default:"" desc:"MQTT broker server hostname or IP address"`
	Port       uint16 `key:"port" default:"8883" desc:"SMTP server port, common ones are 8883, 1883"`
	Topic      string `key:"topic" default:"" desc:"Topic where the message is sent"`
	ClientID   string `key:"clientid" default:"" desc:"client's id from the message is sent"`
	Username   string `key:"username" default:"" desc:"username for auth"`
	Password   string `key:"password" default:"" desc:"password for auth"`
	DisableTLS bool   `key:"disabletls" default:"No"`
}

// Enums returns the fields that should use a corresponding EnumFormatter to Print/Parse their values
func (config *Config) Enums() map[string]types.EnumFormatter {
	return map[string]types.EnumFormatter{}
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
		Host:       fmt.Sprintf("%s:%d", config.Host, config.Port),
		Scheme:     Scheme,
		ForceQuery: true,
		RawQuery:   format.BuildQuery(resolver),
	}

}

func (config *Config) setURL(resolver types.ConfigQueryResolver, url *url.URL) error {

	config.Host = url.Hostname()

	if port, err := strconv.ParseUint(url.Port(), 10, 16); err == nil {
		config.Port = uint16(port)
	}

	for key, vals := range url.Query() {
		if err := resolver.Set(key, vals[0]); err != nil {
			return err
		}
	}

	return nil
}

// MqttURL returns a string that is synchronized with the config props
func (config *Config) MqttURL() string {
	MqttHost := config.Host
	MqttPort := config.Port
	scheme := DefaultWebhookScheme
	if config.DisableTLS {
		scheme = Scheme[:4]
	}
	return fmt.Sprintf("%s://%s:%d", scheme, MqttHost, MqttPort)
}

// GetClientConfig returns the client options
func (config *Config) GetClientConfig(postURL string) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()

	opts.AddBroker(postURL)

	if len(config.ClientID) > 0 {
		opts.SetClientID(config.ClientID)
	}

	if len(config.Username) > 0 {
		opts.SetUsername(config.Username)
	}

	if len(config.Password) > 0 {
		opts.SetPassword(config.Password)
	}

	if !config.DisableTLS {
		tlsConfig := config.GetTLSConfig()
		opts.SetTLSConfig(tlsConfig)
	}

	return opts
}

// GetTLSConfig returns the configuration with the certificates for TLS
func (config *Config) GetTLSConfig() *tls.Config {
	certpool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("ca.crt")

	if err != nil {
		log.Fatalln(err.Error())
	}
	certpool.AppendCertsFromPEM(ca)

	clientKeyPair, err := tls.LoadX509KeyPair("client.crt", "client.key")
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		RootCAs:            certpool,
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{clientKeyPair},
	}
}

const (
	// Scheme is the identifying part of this service's configuration URL
	Scheme = "mqtt"
	// DefaultWebhookScheme is the scheme used for webhook URLs unless overridden
	DefaultWebhookScheme = "mqtts"
)
