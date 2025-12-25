package ntfy

import (
	"net/url"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Config for use within the ntfy service
type Config struct {
	Title    string   `key:"title"       default:""          desc:"Message title"`
	Host     string   `url:"host"        default:"ntfy.sh"   desc:"Server hostname and port"`
	Topic    string   `url:"path"        required:""         desc:"Target topic name"`
	Password string   `url:"password"    optional:""         desc:"Auth password"`
	Username string   `url:"user"        optional:""         desc:"Auth username"`
	Scheme   string   `key:"scheme"      default:"https"     desc:"Server protocol, http or https"`
	Tags     []string `key:"tags"        optional:""         desc:"List of tags that may or not map to emojis"`
	Priority priority `key:"priority"    default:"default"   desc:"Message priority with 1=min, 3=default and 5=max"`
	Actions  []string `key:"actions"     optional:"" sep:";" desc:"Custom user action buttons for notifications, see https://docs.ntfy.sh/publish/#action-buttons"`
	Click    string   `key:"click"       optional:""         desc:"Website opened when notification is clicked"`
	Attach   string   `key:"attach"      optional:""         desc:"URL of an attachment, see attach via URL"`
	Filename string   `key:"filename"    optional:""         desc:"File name of the attachment"`
	Delay    string   `key:"delay,at,in" optional:""         desc:"Timestamp or duration for delayed delivery, see https://docs.ntfy.sh/publish/#scheduled-delivery"`
	Email    string   `key:"email"       optional:""         desc:"E-mail address for e-mail notifications"`
	Icon     string   `key:"icon"        optional:""         desc:"URL to use as notification icon"`
	Cache    bool     `key:"cache"       default:"yes"       desc:"Cache messages"`
	Firebase bool     `key:"firebase"    default:"yes"       desc:"Send to firebase"`
	Template bool     `key:"template"    default:"no"        desc:"Use message and title as template"`
	Message  string   `key:"message"     optional:""         desc:"Message, to be used only then template is set to true"`
}

// Enums implements types.ServiceConfig
func (*Config) Enums() map[string]types.EnumFormatter {
	return map[string]types.EnumFormatter{
		"Priority": Priority.Enum,
	}
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

// GetAPIURL returns the API URL corresponding to the passed endpoint based on the configuration
func (config *Config) GetAPIURL() string {

	path := config.Topic
	if !strings.HasPrefix(config.Topic, "/") {
		path = "/" + path
	}

	var creds *url.Userinfo
	if config.Password != "" {
		creds = url.UserPassword(config.Username, config.Password)
	}

	apiURL := url.URL{
		Scheme: config.Scheme,
		Host:   config.Host,
		Path:   path,
		User:   creds,
	}
	return apiURL.String()
}

func (config *Config) getURL(resolver types.ConfigQueryResolver) *url.URL {
	return &url.URL{
		User:       url.UserPassword(config.Username, config.Password),
		Host:       config.Host,
		Scheme:     Scheme,
		ForceQuery: true,
		Path:       config.Topic,
		RawQuery:   format.BuildQuery(resolver),
	}

}

func (config *Config) setURL(resolver types.ConfigQueryResolver, url *url.URL) error {

	password, _ := url.User.Password()
	config.Password = password
	config.Username = url.User.Username()
	config.Host = url.Host
	config.Topic = strings.TrimPrefix(url.Path, "/")

	// Escape raw `;` in queries
	url.RawQuery = strings.ReplaceAll(url.RawQuery, ";", "%3b")

	for key, vals := range url.Query() {
		if err := resolver.Set(key, vals[0]); err != nil {
			return err
		}
	}

	return nil
}

// Scheme is the identifying part of this service's configuration URL
const (
	Scheme = "ntfy"
)
