package pinglet

import (
	"errors"
	"net/url"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
)

var (
	// ErrTokenMissing is returned when no API key was supplied
	ErrTokenMissing = errors.New("no API key was supplied")

	// ErrNamespaceTopicMissing is returned when the namespace and/or topic are absent from the URL path
	ErrNamespaceTopicMissing = errors.New("both a namespace and a topic are required")
)

// Config for use within the pinglet service
type Config struct {
	Title     string            `key:"title"     default:""       desc:"Notification title, optionally set by the sender"`
	Token     string            `url:"user"                       desc:"Pinglet API key" required:""`
	Host      string            `url:"host,port"                  desc:"Server hostname (and optionally port)" required:""`
	Path      string            `url:"path1"     optional:""       desc:"Path prefix, e.g. a reverse-proxy mount point"`
	Namespace string            `url:"path2"                      desc:"The namespace the topic resides in" required:""`
	Topic     string            `url:"path3"                      desc:"The topic to publish to" required:""`
	Scheme    string            `key:"scheme"    default:"https"  desc:"Server protocol, http or https"`
	Priority  priority          `key:"priority"  default:"normal" desc:"Delivery priority: silent, normal or urgent"`
	Badges    map[string]string `key:"badges"    optional:""      desc:"Badge pills rendered on the notification card (max 3)"`
	Data      map[string]string `key:"data"      optional:""      desc:"Metadata key/value pairs shown on the detail sheet"`
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

// GetAPIURL returns the API URL that the notification is posted to
func (config *Config) GetAPIURL() string {
	apiURL := url.URL{
		Scheme: config.Scheme,
		Host:   config.Host,
		Path:   config.path(),
	}
	return apiURL.String()
}

// path assembles the request path from the optional prefix, namespace and topic
func (config *Config) path() string {
	prefix := ""
	if config.Path != "" {
		prefix = "/" + strings.Trim(config.Path, "/")
	}
	return prefix + "/" + config.Namespace + "/" + config.Topic
}

func (config *Config) getURL(resolver types.ConfigQueryResolver) *url.URL {
	return &url.URL{
		User:       url.User(config.Token),
		Host:       config.Host,
		Scheme:     Scheme,
		ForceQuery: true,
		Path:       config.path(),
		RawQuery:   format.BuildQuery(resolver),
	}
}

func (config *Config) setURL(resolver types.ConfigQueryResolver, url *url.URL) error {
	config.Token = url.User.Username()
	config.Host = url.Host

	// The last two path segments are the namespace and topic; anything before
	// them is treated as a path prefix (such as a reverse-proxy mount point).
	var parts []string
	for _, segment := range strings.Split(url.Path, "/") {
		if segment != "" {
			parts = append(parts, segment)
		}
	}
	if len(parts) < 2 {
		return ErrNamespaceTopicMissing
	}
	config.Topic = parts[len(parts)-1]
	config.Namespace = parts[len(parts)-2]
	config.Path = strings.Join(parts[:len(parts)-2], "/")

	for key, vals := range url.Query() {
		if err := resolver.Set(key, vals[0]); err != nil {
			return err
		}
	}

	if config.Token == "" {
		return ErrTokenMissing
	}

	return nil
}

const (
	// Scheme is the identifying part of this service's configuration URL
	Scheme = "pinglet"
)
