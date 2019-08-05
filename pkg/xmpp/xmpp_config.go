package xmpp

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"gosrc.io/xmpp"

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
	ServerHost string
	ToAddress  string
	Subject    string
}

// GetURL returns a URL representation of it's current field values
func (config *Config) GetURL() *url.URL {

	return &url.URL{
		User:     url.UserPassword(config.Username, config.Password),
		Host:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Path:     "/",
		Scheme:   Scheme,
		RawQuery: format.BuildQuery(config),
	}

}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *Config) SetURL(url *url.URL) error {

	password, _ := url.User.Password()

	config.Username = url.User.Username()
	config.Password = password
	config.Host = url.Hostname()

	if port, err := strconv.ParseUint(url.Port(), 10, 16); err == nil {
		config.Port = uint16(port)
	}

	for key, vals := range url.Query() {
		if err := config.Set(key, vals[0]); err != nil {
			return err
		}
	}

	if len(config.ToAddress) < 1 {
		return errors.New("toAddress missing from config URL")
	}

	return nil
}

// QueryFields returns the fields that are part of the Query of the service URL
func (config *Config) QueryFields() []string {
	return []string{
		"toAddress",
		"subject",
		"serverHost",
	}
}

// Get returns the value of a Query field
func (config *Config) Get(key string) (string, error) {
	switch strings.ToLower(key) {
	case "toaddress":
		return config.ToAddress, nil
	case "subject":
		return config.Subject, nil
	case "serverhost":
		return config.ServerHost, nil
	}
	return "", fmt.Errorf("invalid query key \"%s\"", key)
}

// Set updates the value of a Query field
func (config *Config) Set(key string, value string) error {
	switch strings.ToLower(key) {
	case "toaddress":
		config.ToAddress = value
	case "subject":
		config.Subject = value
	case "serverhost":
		config.ServerHost = value
	default:
		return fmt.Errorf("invalid query key \"%s\"", key)
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
