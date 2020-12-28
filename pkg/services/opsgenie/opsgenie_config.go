package opsgenie

import (
	"fmt"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
)

// Config for use within the opsgenie service
type Config struct {
	ApiKey string `desc:"The OpsGenie API key"`
	Host   string `desc:"The OpsGenie API host. Use 'api.opsgenie.com' for US and 'api.eu.opsgenie.com' for EU instances"`
	standard.EnumlessConfig
	Alias       string
	Description string
	Responders  string
	VisibleTo   string
	Actions     string
	Tags        string
	Details     string
	Entity      string
	Source      string
	Priority    string
	Note        string
	User        string
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

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *Config) SetURL(url *url.URL) error {
	config.Host = url.Hostname() + ":" + url.Port()
	config.ApiKey = url.Path[1:]

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
