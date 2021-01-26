package opsgenie

import (
	"fmt"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/types"
)

// AlertPayload represents the payload being sent to the OpsGenie API
//
// See: https://docs.opsgenie.com/docs/alert-api#create-alert
//
// Some fields contain complex values like arrays and objects.
// Because `params` are strings only we cannot pass in slices
// or maps. Instead we "preserve" the JSON in those fields. That
// way we can pass in complex types as JSON like so:
//
//	service.Send("An example alert message", &types.Params{
//		"alias":       "Life is too short for no alias",
//		"description": "Every alert needs a description",
//		"responders":  `[{"id":"4513b7ea-3b91-438f-b7e4-e3e54af9147c","type":"team"},{"name":"NOC","type":"team"}]`,
//		"visibleTo":   `[{"id":"4513b7ea-3b91-438f-b7e4-e3e54af9147c","type":"team"},{"name":"rocket_team","type":"team"}]`,
//		"details":     `{"key1": "value1", "key2": "value2"}`,
//	})
type AlertPayload struct {
	Message     string            `json:"message"`
	Alias       string            `json:"alias,omitempty"`
	Description string            `json:"description,omitempty"`
	Responders  []Entity          `json:"responders,omitempty"`
	VisibleTo   []Entity          `json:"visibleTo,omitempty"`
	Actions     []string          `json:"actions,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	Details     map[string]string `json:"details,omitempty"`
	Entity      string            `json:"entity,omitempty"`
	Source      string            `json:"source,omitempty"`
	Priority    string            `json:"priority,omitempty"`
	User        string            `json:"user,omitempty"`
	Note        string            `json:"note,omitempty"`
}

// TODO: Duplicated code, see: formatter.go
func deserializeSlice(str string) []string {
	return strings.Split(str, ",")
}

// TODO: Duplicated code, see: formatter.go
func deserializeMap(str string) (map[string]string, error) {
	result := make(map[string]string)

	pairs := strings.Split(str, ",")
	for _, pair := range pairs {
		elems := strings.Split(pair, ":")
		if len(elems) != 2 {
			return result, fmt.Errorf("field format is not supported")
		}
		key := elems[0]
		value := elems[1]

		result[key] = value
	}

	return result, nil
}

// NewAlertPayload instantiates AlertPayload
func NewAlertPayload(message string, config *Config, params *types.Params) (AlertPayload, error) {
	if params == nil {
		params = &types.Params{}
	}

	result := AlertPayload{
		Message: message,
		// Populate with values from the query string as defaults
		Alias:       config.Alias,
		Description: config.Description,
		Responders:  config.Responders,
		VisibleTo:   config.VisibleTo,
		Actions:     config.Actions,
		Tags:        config.Tags,
		Details:     config.Details,
		Entity:      config.Entity,
		Source:      config.Source,
		Priority:    config.Priority,
		User:        config.User,
		Note:        config.Note,
	}

	// Populate with values from runtime parameters
	for key, value := range *params {
		var err error

		switch key {
		case "alias":
			result.Alias = value
		case "description":
			result.Description = value
		case "responders":
			result.Responders, err = deserializeEntities(value)
		case "visibleTo":
			result.VisibleTo, err = deserializeEntities(value)
		case "actions":
			result.Actions = deserializeSlice(value)
		case "tags":
			result.Tags = deserializeSlice(value)
		case "details":
			result.Details, err = deserializeMap(value)
		case "entity":
			result.Entity = value
		case "source":
			result.Source = value
		case "priority":
			result.Priority = value
		case "user":
			result.User = value
		case "note":
			result.Note = value
		default:
			return result, fmt.Errorf("unknown config key: %q", key)
		}

		if err != nil {
			return result, err
		}

	}

	return result, nil
}
