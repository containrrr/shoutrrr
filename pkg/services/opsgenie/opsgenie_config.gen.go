// Code generated by "shoutrrr-gen "; DO NOT EDIT.
package opsgenie

import (
	"fmt"
	"net/url"
	_ "strings"

	"github.com/containrrr/shoutrrr/pkg/conf"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// (‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾)
//  )  Props                          (
// (___________________________________)

type Config struct {
	APIKey      string            `url:"path1" `
	Actions     []string          `key:"actions" `
	Alias       string            `key:"alias" `
	Description string            `key:"description" `
	Details     map[string]string `key:"details" `
	Entity      string            `key:"entity" `
	Host        string            `url:"host" `
	Note        string            `key:"note" `
	Port        int64             `url:"port" `
	Priority    string            `key:"priority" `
	Responders  []Entity          `key:"responders" `
	Source      string            `key:"source" `
	Tags        []string          `key:"tags" `
	Title       string            `key:"title" `
	User        string            `key:"user" `
	VisibleTo   []Entity          `key:"visibleto" `
}

type configProp int

const (
	propAPIKey      configProp = 0
	propActions     configProp = 1
	propAlias       configProp = 2
	propDescription configProp = 3
	propDetails     configProp = 4
	propEntity      configProp = 5
	propHost        configProp = 6
	propNote        configProp = 7
	propPort        configProp = 8
	propPriority    configProp = 9
	propResponders  configProp = 10
	propSource      configProp = 11
	propTags        configProp = 12
	propTitle       configProp = 13
	propUser        configProp = 14
	propVisibleTo   configProp = 15
	propCount                  = 16
)

var propInfo = types.ConfigPropInfo{
	PropNames: []string{
		"APIKey",
		"Actions",
		"Alias",
		"Description",
		"Details",
		"Entity",
		"Host",
		"Note",
		"Port",
		"Priority",
		"Responders",
		"Source",
		"Tags",
		"Title",
		"User",
		"VisibleTo",
	},

	// Note that propKeys may not align with propNames, as a property can have no or multiple keys
	Keys: []string{
		"actions",
		"alias",
		"description",
		"details",
		"entity",
		"note",
		"priority",
		"responders",
		"source",
		"tags",
		"title",
		"user",
		"visibleto",
	},

	DefaultValues: []string{
		"",
		"",
		"",
		"",
		"",
		"",
		"api.opsgenie.com",
		"",
		"443",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
	},

	PrimaryKeys: []int{
		-1,
		0,
		1,
		2,
		3,
		4,
		-1,
		5,
		-1,
		6,
		7,
		8,
		9,
		10,
		11,
		12,
	},

	KeyPropIndexes: map[string]int{
		"actions":     1,
		"alias":       2,
		"description": 3,
		"details":     4,
		"entity":      5,
		"note":        7,
		"priority":    9,
		"responders":  10,
		"source":      11,
		"tags":        12,
		"title":       13,
		"user":        14,
		"visibleto":   15,
	},
}

func (_ *Config) PropInfo() *types.ConfigPropInfo {
	return &propInfo
}

// (‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾)
//  )  GetURL                         (
// (___________________________________)

// GetURL returns a URL representation of it's current field values
func (config *Config) GetURL() *url.URL {
	return &url.URL{
		// Userinfo fields are not used for configuration
		Host:     conf.FormatHost(config.Host, config.Port),
		Path:     conf.JoinPath(string(config.APIKey)),
		RawQuery: conf.QueryValues(config).Encode(),
		Scheme:   Scheme,
	}
}

// (‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾)
//  )  SetURL                         (
// (___________________________________)

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *Config) SetURL(configURL *url.URL) error {
	if lc, ok := (interface{})(config).(types.ConfigWithLegacyURLSupport); ok {
		configURL = lc.UpdateLegacyURL(configURL)
	}
	updates := make(map[int]string, propCount)
	if port := configURL.Port(); port != "" {
		updates[int(propPort)] = port
	}
	updates[int(propHost)] = configURL.Hostname()

	pathParts := conf.SplitPath(configURL.Path)
	if len(pathParts) > 0 {
		updates[int(propAPIKey)] = pathParts[0]
	}
	if len(pathParts) > 1 {
		return fmt.Errorf("too many path items: %v, expected 1", len(pathParts))
	}

	for key, value := range configURL.Query() {

		if propIndex, found := propInfo.PropIndexFor(key); found {
			updates[propIndex] = value[0]
		} else if key != "title" {
			return fmt.Errorf("invalid key %q", key)
		}
	}

	err := config.Update(updates)
	if err != nil {
		return err
	}

	if config.APIKey == "" {
		return fmt.Errorf("apiKey missing from config URL")
	}

	if !conf.ValueMatchesPattern(config.Priority, "(P[1-5])?") {
		return fmt.Errorf("value %v for priority does not match the expected format", config.Priority)
	}

	return nil
}

// (‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾)
//  )  Enums / Options                (
// (___________________________________)

func (config *Config) Enums() map[string]types.EnumFormatter {
	return map[string]types.EnumFormatter{}
}

// Update updates the Config from a map of it's properties
func (config *Config) Update(updates map[int]string) error {
	var last_err error
	for index, value := range updates {
		switch configProp(index) {
		case propAPIKey:
			if val, err := conf.ParseTextValue(value); err != nil {
				last_err = err
			} else {
				config.APIKey = val
			}
		case propActions:
			if val, err := conf.ParseListValue(value, ","); err != nil {
				last_err = err
			} else {
				config.Actions = val
			}
		case propAlias:
			if val, err := conf.ParseTextValue(value); err != nil {
				last_err = err
			} else {
				config.Alias = val
			}
		case propDescription:
			if val, err := conf.ParseTextValue(value); err != nil {
				last_err = err
			} else {
				config.Description = val
			}
		case propDetails:
			if val, err := conf.ParseMapValue(value, ",", ":"); err != nil {
				last_err = err
			} else {
				config.Details = val
			}
		case propEntity:
			if val, err := conf.ParseTextValue(value); err != nil {
				last_err = err
			} else {
				config.Entity = val
			}
		case propHost:
			if val, err := conf.ParseTextValue(value); err != nil {
				last_err = err
			} else {
				config.Host = val
			}
		case propNote:
			if val, err := conf.ParseTextValue(value); err != nil {
				last_err = err
			} else {
				config.Note = val
			}
		case propPort:
			if val, err := conf.ParseNumberValue(value, 10); err != nil {
				last_err = err
			} else {
				config.Port = val
			}
		case propPriority:
			if val, err := conf.ParseTextValue(value); err != nil {
				last_err = err
			} else {
				config.Priority = val
			}
		case propResponders:
			if val, err := parseEntityItems(value); err != nil {
				last_err = err
			} else {
				config.Responders = val
			}
		case propSource:
			if val, err := conf.ParseTextValue(value); err != nil {
				last_err = err
			} else {
				config.Source = val
			}
		case propTags:
			if val, err := conf.ParseListValue(value, ","); err != nil {
				last_err = err
			} else {
				config.Tags = val
			}
		case propTitle:
			if val, err := conf.ParseTextValue(value); err != nil {
				last_err = err
			} else {
				config.Title = val
			}
		case propUser:
			if val, err := conf.ParseTextValue(value); err != nil {
				last_err = err
			} else {
				config.User = val
			}
		case propVisibleTo:
			if val, err := parseEntityItems(value); err != nil {
				last_err = err
			} else {
				config.VisibleTo = val
			}
		default:
			return fmt.Errorf("invalid key")
		}
		if last_err != nil {
			return fmt.Errorf("failed to set value for %v: %v", propInfo.PropNames[index], last_err)
		}
	}
	return nil
}

// Update updates the Config from a map of it's properties
func (config *Config) PropValue(prop int) string {
	switch configProp(prop) {
	case propAPIKey:
		return conf.FormatTextValue(config.APIKey)
	case propActions:
		return conf.FormatListValue(config.Actions, ",")
	case propAlias:
		return conf.FormatTextValue(config.Alias)
	case propDescription:
		return conf.FormatTextValue(config.Description)
	case propDetails:
		return conf.FormatMapValue(config.Details, ",", ":")
	case propEntity:
		return conf.FormatTextValue(config.Entity)
	case propHost:
		return conf.FormatTextValue(config.Host)
	case propNote:
		return conf.FormatTextValue(config.Note)
	case propPort:
		return conf.FormatNumberValue(config.Port, 10)
	case propPriority:
		return conf.FormatTextValue(config.Priority)
	case propResponders:
		return conf.FormatListValue(formatEntityItems(config.Responders), ",")
	case propSource:
		return conf.FormatTextValue(config.Source)
	case propTags:
		return conf.FormatListValue(config.Tags, ",")
	case propTitle:
		return conf.FormatTextValue(config.Title)
	case propUser:
		return conf.FormatTextValue(config.User)
	case propVisibleTo:
		return conf.FormatListValue(formatEntityItems(config.VisibleTo), ",")
	default:
		return ""
	}
}
