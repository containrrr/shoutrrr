package smtp

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Config is the configuration needed to send e-mail notifications over SMTP
type Config struct {
	Host        string    `desc:"SMTP server hostname or IP address"`
	Username    string    `desc:"authentication username"`
	Password    string    `desc:"authentication password or hash"`
	Port        uint16    `desc:"SMTP server port, common ones are 25, 465, 587 or 2525" default:"25"`
	FromAddress string    `desc:"e-mail address that the mail are sent from" key:"fromaddress"`
	FromName    string    `desc:"name of the sender" optional:"yes" key:"fromname"`
	ToAddresses []string  `desc:"list of recipient e-mails separated by \",\" (comma)" key:"toaddresses"`
	Subject     string    `desc:"the subject of the sent mail" key:"subject,title" default:"Shoutrrr Notification"`
	Auth        authType  `desc:"SMTP authentication method" key:"auth"`
	Encryption  encMethod `desc:"Encryption method" default:"Auto" key:"encryption"`
	UseStartTLS bool      `desc:"attempt to use SMTP StartTLS encryption" default:"Yes" key:"starttls"`
	UseHTML     bool      `desc:"whether the message being sent is in HTML" default:"No" key:"usehtml"`
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
		User:       url.UserPassword(config.Username, config.Password),
		Host:       fmt.Sprintf("%s:%d", config.Host, config.Port),
		Path:       "/",
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

	if port, err := strconv.ParseUint(url.Port(), 10, 16); err == nil {
		config.Port = uint16(port)
	}

	for key, vals := range url.Query() {
		if err := resolver.Set(key, vals[0]); err != nil {
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

// Clone returns a copy of the config
func (config *Config) Clone() Config {
	clone := *config
	clone.ToAddresses = make([]string, len(config.ToAddresses))
	copy(clone.ToAddresses, config.ToAddresses)
	return clone
}

// Enums returns the fields that should use a corresponding EnumFormatter to Print/Parse their values
func (config Config) Enums() map[string]types.EnumFormatter {
	return map[string]types.EnumFormatter{
		"Auth":       AuthTypes.Enum,
		"Encryption": EncMethods.Enum,
	}
}

// Scheme is the identifying part of this service's configuration URL
const Scheme = "smtp"
