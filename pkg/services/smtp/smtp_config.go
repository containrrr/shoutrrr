package smtp

import (
	"errors"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
	"net/url"
	"strconv"
	"strings"
)

// Config is the configuration needed to send e-mail notifications over SMTP
type Config struct {
	Host        string   `desc:"SMTP server hostname or IP address"`
	Username    string   `desc:"authentication username"`
	Password    string   `desc:"authentication password or hash"`
	Port        uint16   `desc:"SMTP server port, common ones are 25, 465, 587 or 2525" default:"25"`
	FromAddress string   `desc:"e-mail address that the mail are sent from"`
	FromName    string   `desc:"name of the sender" optional:"yes"`
	ToAddresses []string `desc:"list of recipient e-mails separated by \",\" (comma)"`
	Subject     string   `desc:"the subject of the sent mail" tpl:"subject" default:"Shoutrrr Notification"`
	Auth        authType `desc:"SMTP authentication method"`
	UseStartTLS bool     `desc:"attempt to use SMTP StartTLS encryption" default:"Yes"`
	UseHTML     bool     `desc:"whether the message being sent is in HTML" default:"No"`
}

// GetURL returns a URL representation of it's current field values
func (config *Config) GetURL() *url.URL {

	return &url.URL{
		User:       url.UserPassword(config.Username, config.Password),
		Host:       fmt.Sprintf("%s:%d", config.Host, config.Port),
		Path:       "/",
		Scheme:     Scheme,
		ForceQuery: true,
		RawQuery:   format.BuildQuery(config),
	}

}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *Config) SetURL(url *url.URL) error {
	hostParts := strings.Split(url.Host, ":")
	host := hostParts[0]
	port, err := strconv.ParseUint(hostParts[1], 10, 16)
	if err != nil {
		return err
	}
	password, _ := url.User.Password()

	config.Username = url.User.Username()
	config.Password = password
	config.Host = host
	config.Port = uint16(port)

	for key, vals := range url.Query() {
		if err := config.Set(key, vals[0]); err != nil {
			return err
		}
	}

	if len(config.FromAddress) < 1 {
		return errors.New("fromAddress missing from config URL")
	}

	if len(config.ToAddresses) < 1 {
		return errors.New("toAddress missing from config URL")
	}

	return nil
}

// QueryFields returns the fields that are part of the Query of the service URL
func (config *Config) QueryFields() []string {
	return []string {
	"fromAddress",
	"fromName",
	"toAddresses",
	"auth",
	"subject",
	"startTls",
	"useHTML",
	}
}

// Get returns the value of a Query field
func (config *Config) Get(key string) (string, error) {
	switch key {
	case "fromAddress":
		return config.FromAddress, nil
	case "fromName":
		return config.FromName, nil
	case "toAddresses":
		return strings.Join(config.ToAddresses, ","), nil
	case "auth":
		return config.Auth.String(), nil
	case "subject":
		return config.Subject, nil
	case "startTls":
		return format.PrintBool(config.UseStartTLS), nil
	case "useHTML":
		return format.PrintBool(config.UseHTML), nil
	}
	return "", fmt.Errorf("invalid query key \"%s\"", key)
}

// Set updates the value of a Query field
func (config *Config) Set(key string, value string) error {
	switch key {
	case "fromAddress":
		config.FromAddress = value
	case "fromName":
		config.FromName = value
	case "toAddresses":
		config.ToAddresses = strings.Split(value, ",")
	case "auth":
		config.Auth = parseAuth(value)
	case "subject":
		config.Subject = value
	case "startTls":
		config.UseStartTLS, _ = format.ParseBool(value, true)
	case "useHTML":
		config.UseHTML, _ = format.ParseBool(value, false)
	default:
		return fmt.Errorf("invalid query key \"%s\"", key)
	}
	return nil
}

// CreateConfigFromURL creates a Config struct given a valid discord notification url
func (plugin *Service) CreateConfigFromURL(url *url.URL) (*Config, error) {

	config := &Config{}
	err := config.SetURL(url)

	return config, err
}

// Enums returns the fields that should use a corresponding EnumFormatter to Print/Parse their values
func (config Config) Enums() map[string]types.EnumFormatter {
	return map[string]types.EnumFormatter{
		"Auth": authTypes.Enum,
	}
}

// Scheme is the identifying part of this service's configuration URL
const Scheme = "smtp"