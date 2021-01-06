package opsgenie

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
)

// Config for use within the opsgenie service
type Config struct {
	ApiKey string `desc:"The OpsGenie API key"`
	Host   string `desc:"The OpsGenie API host. Use 'api.eu.opsgenie.com' for EU instances" default:"api.opsgenie.com"`
	Port   uint16 `desc:"The OpsGenie API port." default:"443"`
	standard.EnumlessConfig
	Alias       string `desc:"Client-defined identifier of the alert" optional:"true"`
	Description string `desc:"Description field of the alert" optional:"true"`
	Responders  string `desc:"Teams, users, escalations and schedules that the alert will be routed to send notifications" optional:"true"`
	VisibleTo   string `desc:"Teams and users that the alert will become visible to without sending any notification" optional:"true"`
	Actions     string `desc:"Custom actions that will be available for the alert" optional:"true"`
	Tags        string `desc:"Tags of the alert" optional:"true"`
	Details     string `desc:"Map of key-value pairs to use as custom properties of the alert" optional:"true"`
	Entity      string `desc:"Entity field of the alert that is generally used to specify which domain the Source field of the alert" optional:"true"`
	Source      string `desc:"Source field of the alert" optional:"true"`
	Priority    string `desc:"Priority level of the alert. Possible values are P1, P2, P3, P4 and P5" optional:"true"`
	Note        string `desc:"Additional note that will be added while creating the alert" optional:"true"`
	User        string `desc:"Display name of the request owner" optional:"true"`
}

// QueryFields returns the fields that are part of the Query of the service URL
func (config *Config) QueryFields() []string {
	return []string{
		"alias",
		"responders",
		"description",
		"visibleTo",
		"actions",
		"tags",
		"details",
		"entity",
		"source",
		"priority",
		"note",
		"user",
	}
}

// Set updates the value of a Query field
func (config *Config) Set(key string, value string) error {
	switch key {
	case "alias":
		config.Alias = value
	case "description":
		config.Description = value
	case "responders":
		config.Responders = value
	case "visibleTo":
		config.VisibleTo = value
	case "actions":
		config.Actions = value
	case "tags":
		config.Tags = value
	case "details":
		config.Details = value
	case "entity":
		config.Entity = value
	case "source":
		config.Source = value
	case "priority":
		config.Priority = value
	case "note":
		config.Note = value
	case "user":
		config.User = value
	default:
		return fmt.Errorf("invalid query key \"%s\"", key)
	}

	return nil
}

// Get returns the values of a Query field
func (config *Config) Get(key string) (string, error) {
	switch key {
	case "alias":
		return config.Alias, nil
	case "description":
		return config.Description, nil
	case "responders":
		return config.Responders, nil
	case "visibleTo":
		return config.VisibleTo, nil
	case "actions":
		return config.Actions, nil
	case "tags":
		return config.Tags, nil
	case "details":
		return config.Details, nil
	case "entity":
		return config.Entity, nil
	case "source":
		return config.Source, nil
	case "priority":
		return config.Priority, nil
	case "note":
		return config.Note, nil
	case "user":
		return config.User, nil
	}

	return "", fmt.Errorf("invalid query key \"%s\"", key)
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *Config) SetURL(url *url.URL) error {
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
		if err := config.Set(key, vals[0]); err != nil {
			return err
		}
	}

	return nil
}

func (config *Config) GetURL() *url.URL {
	return &url.URL{
		Host:       fmt.Sprintf("%s:%d", config.Host, config.Port),
		Path:       fmt.Sprintf("/%s", config.ApiKey),
		Scheme:     Scheme,
		ForceQuery: true,
		RawQuery:   format.BuildQuery(config),
	}
}

const (
	// Scheme is the identifying part of this service's configuration URL
	Scheme = "opsgenie"
)
