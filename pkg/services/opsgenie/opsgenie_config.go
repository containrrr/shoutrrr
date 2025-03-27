package opsgenie

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/nicholas-fedor/shoutrrr/pkg/format"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
)

const defaultPort = 443

// Config for use within the opsgenie service.
type Config struct {
	APIKey      string            `desc:"The OpsGenie API key"                                                                                   url:"path"`
	Host        string            `default:"api.opsgenie.com"                                                                                    desc:"The OpsGenie API host. Use 'api.eu.opsgenie.com' for EU instances" url:"host"`
	Port        uint16            `default:"443"                                                                                                 desc:"The OpsGenie API port."                                            url:"port"`
	Alias       string            `desc:"Client-defined identifier of the alert"                                                                 key:"alias"                                                              optional:"true"`
	Description string            `desc:"Description field of the alert"                                                                         key:"description"                                                        optional:"true"`
	Responders  []Entity          `desc:"Teams, users, escalations and schedules that the alert will be routed to send notifications"            key:"responders"                                                         optional:"true"`
	VisibleTo   []Entity          `desc:"Teams and users that the alert will become visible to without sending any notification"                 key:"visibleTo"                                                          optional:"true"`
	Actions     []string          `desc:"Custom actions that will be available for the alert"                                                    key:"actions"                                                            optional:"true"`
	Tags        []string          `desc:"Tags of the alert"                                                                                      key:"tags"                                                               optional:"true"`
	Details     map[string]string `desc:"Map of key-value pairs to use as custom properties of the alert"                                        key:"details"                                                            optional:"true"`
	Entity      string            `desc:"Entity field of the alert that is generally used to specify which domain the Source field of the alert" key:"entity"                                                             optional:"true"`
	Source      string            `desc:"Source field of the alert"                                                                              key:"source"                                                             optional:"true"`
	Priority    string            `desc:"Priority level of the alert. Possible values are P1, P2, P3, P4 and P5"                                 key:"priority"                                                           optional:"true"`
	Note        string            `desc:"Additional note that will be added while creating the alert"                                            key:"note"                                                               optional:"true"`
	User        string            `desc:"Display name of the request owner"                                                                      key:"user"                                                               optional:"true"`
	Title       string            `default:""                                                                                                    desc:"notification title, optionally set by the sender"                  key:"title"`
}

// Enums returns an empty map because the OpsGenie service doesn't use Enums.
func (config Config) Enums() map[string]types.EnumFormatter {
	return map[string]types.EnumFormatter{}
}

// GetURL is the public version of getURL that creates a new PropKeyResolver when accessed from outside the package.
func (config *Config) GetURL() *url.URL {
	resolver := format.NewPropKeyResolver(config)

	return config.getURL(&resolver)
}

// Private version of GetURL that can use an instance of PropKeyResolver
// instead of rebuilding it's model from reflection.
func (config *Config) getURL(resolver types.ConfigQueryResolver) *url.URL {
	host := ""
	if config.Port > 0 {
		host = fmt.Sprintf("%s:%d", config.Host, config.Port)
	} else {
		host = config.Host
	}

	result := &url.URL{
		Host:     host,
		Path:     fmt.Sprintf("/%s", config.APIKey),
		Scheme:   Scheme,
		RawQuery: format.BuildQuery(resolver),
	}

	return result
}

// SetURL updates a ServiceConfig from a URL representation of it's field values.
func (config *Config) SetURL(url *url.URL) error {
	resolver := format.NewPropKeyResolver(config)

	return config.setURL(&resolver, url)
}

// Private version of SetURL that can use an instance of PropKeyResolver
// instead of rebuilding it's model from reflection.
func (config *Config) setURL(resolver types.ConfigQueryResolver, url *url.URL) error {
	config.Host = url.Hostname()

	if url.String() != "opsgenie://dummy@dummy.com" {
		if len(url.Path) > 0 {
			config.APIKey = url.Path[1:]
		} else {
			return fmt.Errorf("API key missing from config URL path")
		}
	}

	if url.Port() != "" {
		port, err := strconv.ParseUint(url.Port(), 10, 16)
		if err != nil {
			return err
		}

		config.Port = uint16(port)
	} else {
		config.Port = 443
	}

	for key, vals := range url.Query() {
		if err := resolver.Set(key, vals[0]); err != nil {
			return err
		}
	}

	return nil
}

const (
	// Scheme is the identifying part of this service's configuration URL.
	Scheme = "opsgenie"
)
