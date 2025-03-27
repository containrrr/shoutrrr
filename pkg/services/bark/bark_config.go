package bark

import (
	"net/url"
	"strings"

	"github.com/nicholas-fedor/shoutrrr/pkg/format"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/standard"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
)

// Config for use within the bark service.
type Config struct {
	standard.EnumlessConfig
	Title     string `default:""                      desc:"Notification title, optionally set by the sender"           key:"title"`
	Host      string `desc:"Server hostname and port" url:"host"`
	Path      string `default:"/"                     desc:"Server path"                                                url:"path"`
	DeviceKey string `desc:"The key for each device"  url:"password"`
	Scheme    string `default:"https"                 desc:"Server protocol, http or https"                             key:"scheme"`
	Sound     string `default:""                      desc:"Value from https://github.com/Finb/Bark/tree/master/Sounds" key:"sound"`
	Badge     int64  `default:"0"                     desc:"The number displayed next to App icon"                      key:"badge"`
	Icon      string `default:""                      desc:"An url to the icon, available only on iOS 15 or later"      key:"icon"`
	Group     string `default:""                      desc:"The group of the notification"                              key:"group"`
	URL       string `default:""                      desc:"Url that will jump when click notification"                 key:"url"`
	Category  string `default:""                      desc:"Reserved field, no use yet"                                 key:"category"`
	Copy      string `default:""                      desc:"The value to be copied"                                     key:"copy"`
}

// GetURL returns a URL representation of it's current field values.
func (config *Config) GetURL() *url.URL {
	resolver := format.NewPropKeyResolver(config)

	return config.getURL(&resolver)
}

// SetURL updates a ServiceConfig from a URL representation of it's field values.
func (config *Config) SetURL(url *url.URL) error {
	resolver := format.NewPropKeyResolver(config)

	return config.setURL(&resolver, url)
}

// GetAPIURL returns the API URL corresponding to the passed endpoint based on the configuration.
func (config *Config) GetAPIURL(endpoint string) string {
	path := strings.Builder{}
	if !strings.HasPrefix(config.Path, "/") {
		path.WriteByte('/')
	}

	_, _ = path.WriteString(config.Path)

	if !strings.HasSuffix(path.String(), "/") {
		path.WriteByte('/')
	}

	path.WriteString(endpoint)

	apiURL := url.URL{
		Scheme: config.Scheme,
		Host:   config.Host,
		Path:   path.String(),
	}

	return apiURL.String()
}

func (config *Config) getURL(resolver types.ConfigQueryResolver) *url.URL {
	return &url.URL{
		User:       url.UserPassword("", config.DeviceKey),
		Host:       config.Host,
		Scheme:     Scheme,
		ForceQuery: true,
		Path:       config.Path,
		RawQuery:   format.BuildQuery(resolver),
	}
}

func (config *Config) setURL(resolver types.ConfigQueryResolver, url *url.URL) error {
	password, _ := url.User.Password()
	config.DeviceKey = password
	config.Host = url.Host
	config.Path = url.Path

	for key, vals := range url.Query() {
		if err := resolver.Set(key, vals[0]); err != nil {
			return err
		}
	}

	return nil
}

// Scheme is the identifying part of this service's configuration URL.
const (
	Scheme = "bark"
)
