package opsgenie

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Config for use within the opsgenie service
type Config struct {
	ApiKey      string   `desc:"The OpsGenie API key"`
	Host        string   `desc:"The OpsGenie API host. Use 'api.eu.opsgenie.com' for EU instances" default:"api.opsgenie.com"`
	Port        uint16   `desc:"The OpsGenie API port." default:"443"`
	Alias       string   `key:"alias" desc:"Client-defined identifier of the alert" optional:"true"`
	Description string   `key:"description" desc:"Description field of the alert" optional:"true"`
	Responders  []Entity `key:"responders" desc:"Teams, users, escalations and schedules that the alert will be routed to send notifications" optional:"true"`
	VisibleTo   []Entity `key:"visibleTo" desc:"Teams and users that the alert will become visible to without sending any notification" optional:"true"`
	Actions     []string `key:"actions" desc:"Custom actions that will be available for the alert" optional:"true"`
	Tags        []string `key:"tags" desc:"Tags of the alert" optional:"true"`
	Details     string   `key:"details" desc:"Map of key-value pairs to use as custom properties of the alert" optional:"true"`
	Entity      string   `key:"entity" desc:"Entity field of the alert that is generally used to specify which domain the Source field of the alert" optional:"true"`
	Source      string   `key:"source" desc:"Source field of the alert" optional:"true"`
	Priority    string   `key:"priority" desc:"Priority level of the alert. Possible values are P1, P2, P3, P4 and P5" optional:"true"`
	Note        string   `key:"note" desc:"Additional note that will be added while creating the alert" optional:"true"`
	User        string   `key:"user" desc:"Display name of the request owner" optional:"true"`
}

func (config Config) Enums() map[string]types.EnumFormatter {
	return map[string]types.EnumFormatter{}
}

// Public version of GetURL that creates a new PropKeyResolver when accessed from outside the package
func (config *Config) GetURL() *url.URL {
	resolver := format.NewPropKeyResolver(config)
	return config.getURL(&resolver)
}

// Private version of GetURL that can use an instance of PropKeyResolver instead of rebuilding it's model from reflection
func (config *Config) getURL(resolver types.ConfigQueryResolver) *url.URL {
	host := ""
	if config.Port > 0 {
		host = fmt.Sprintf("%s:%d", config.Host, config.Port)
	} else {
		host = config.Host
	}

	result := &url.URL{
		Host:     host,
		Path:     fmt.Sprintf("/%s", config.ApiKey),
		Scheme:   Scheme,
		RawQuery: format.BuildQuery(resolver),
	}

	// Add OpsGenie specific data types
	q := result.Query()
	if len(config.Responders) > 0 {
		// TODO: Error handling
		responders, _ := serializeEntities(config.Responders)
		q.Set("responders", responders)
	}
	if len(config.VisibleTo) > 0 {
		// TODO: Error handling
		visibleTo, _ := serializeEntities(config.VisibleTo)
		q.Set("visibleTo", visibleTo)
	}
	result.RawQuery = q.Encode()

	return result
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *Config) SetURL(url *url.URL) error {
	resolver := format.NewPropKeyResolver(config)
	return config.setURL(&resolver, url)
}

// Private version of SetURL that can use an instance of PropKeyResolver instead of rebuilding it's model from reflection
func (config *Config) setURL(resolver types.ConfigQueryResolver, url *url.URL) error {
	config.Host = url.Hostname()
	config.ApiKey = url.Path[1:]

	if url.Port() != "" {
		port, err := strconv.ParseUint(url.Port(), 10, 16)
		if err != nil {
			return err
		}
		config.Port = uint16(port)
	}

	for key, vals := range url.Query() {
		var err error

		switch key {
		case "responders":
			config.Responders, err = deserializeEntities(vals[0])
		case "visibleTo":
			config.VisibleTo, err = deserializeEntities(vals[0])
		default:
			err = resolver.Set(key, vals[0])
		}

		if err != nil {
			return err
		}
	}

	return nil
}

const (
	// Scheme is the identifying part of this service's configuration URL
	Scheme = "opsgenie"
)
