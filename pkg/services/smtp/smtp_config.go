package smtp

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util"
)

// Config is the configuration needed to send e-mail notifications over SMTP
type Config struct {
	Host        string    `desc:"SMTP server hostname or IP address" url:"Host"`
	Username    string    `desc:"SMTP server username" default:"" url:"User"`
	Password    string    `desc:"SMTP server password or hash (for OAuth2)" default:"" url:"Pass"`
	Port        uint16    `desc:"SMTP server port, common ones are 25, 465, 587 or 2525" default:"25" url:"Port"`
	FromAddress string    `desc:"E-mail address that the mail are sent from" key:"fromaddress,from"`
	FromName    string    `desc:"Name of the sender" optional:"yes" key:"fromname"`
	ToAddresses []string  `desc:"List of recipient e-mails separated by \",\" (comma)" key:"toaddresses,to"`
	Subject     string    `desc:"The subject of the sent mail" key:"subject,title" default:"Shoutrrr Notification"`
	Auth        authType  `desc:"SMTP authentication method" key:"auth" default:"Unknown"`
	Encryption  encMethod `desc:"Encryption method" default:"Auto" key:"encryption"`
	UseStartTLS bool      `desc:"Whether to use StartTLS encryption" default:"Yes" key:"usestarttls,starttls"`
	UseHTML     bool      `desc:"Whether the message being sent is in HTML" default:"No" key:"usehtml"`
	ClientHost  string    `desc:"The client host name sent to the SMTP server during HELLO phase. If set to \"auto\" it will use the OS hostname" key:"clienthost" default:"localhost"`
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
		User:       util.URLUserPassword(config.Username, config.Password),
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

// FixEmailTags replaces parsed spaces (+) in e-mail addresses with '+'
func (config *Config) FixEmailTags() {
	config.FromAddress = strings.ReplaceAll(config.FromAddress, " ", "+")
	for i, adr := range config.ToAddresses {
		config.ToAddresses[i] = strings.ReplaceAll(adr, " ", "+")
	}
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
