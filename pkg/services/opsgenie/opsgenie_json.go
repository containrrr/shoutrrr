package opsgenie

import (
	"encoding/json"
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
	Message     string          `json:"message"`
	Alias       string          `json:"alias,omitempty"`
	Description string          `json:"description,omitempty"`
	Responders  json.RawMessage `json:"responders,omitempty"`
	VisibleTo   json.RawMessage `json:"visibleTo,omitempty"`
	Actions     json.RawMessage `json:"actions,omitempty"`
	Tags        json.RawMessage `json:"tags,omitempty"`
	Details     json.RawMessage `json:"details,omitempty"`
	Entity      string          `json:"entity,omitempty"`
	Source      string          `json:"source,omitempty"`
	Priority    string          `json:"priority,omitempty"`
	User        string          `json:"user,omitempty"`
	Note        string          `json:"note,omitempty"`
}

func (j AlertPayload) setStringValue(variable *string, key string, params *types.Params) {
	paramValue, ok := (*params)[key]
	if ok {
		*variable = paramValue
	}
}

func (j AlertPayload) setRawMessageValue(variable *json.RawMessage, key string, params *types.Params) {
	paramValue, ok := (*params)[key]
	if ok {
		*variable = json.RawMessage(paramValue)
	}
}

// NewAlertPayload instantiates AlertPayload
func NewAlertPayload(message string, config *Config, params *types.Params) AlertPayload {
	if params == nil {
		params = &types.Params{}
	}

	result := AlertPayload{
		Message: message,
		// Populate with values from the query string as defaults
		Alias:       config.Alias,
		Description: config.Description,
		Responders:  json.RawMessage(config.Responders),
		VisibleTo:   json.RawMessage(config.VisibleTo),
		Actions:     json.RawMessage(config.Actions),
		Tags:        json.RawMessage(config.Tags),
		Details:     json.RawMessage(config.Details),
		Entity:      config.Entity,
		Source:      config.Source,
		Priority:    config.Priority,
		User:        config.User,
		Note:        config.Note,
	}

	// If specified, overwrite the values from the query string
	result.setStringValue(&result.Alias, "alias", params)
	result.setStringValue(&result.Description, "description", params)
	result.setRawMessageValue(&result.Responders, "responders", params)
	result.setRawMessageValue(&result.VisibleTo, "visibleTo", params)
	result.setRawMessageValue(&result.Actions, "actions", params)
	result.setRawMessageValue(&result.Tags, "tags", params)
	result.setRawMessageValue(&result.Details, "details", params)
	result.setStringValue(&result.Entity, "entity", params)
	result.setStringValue(&result.Source, "source", params)
	result.setStringValue(&result.Priority, "priority", params)
	result.setStringValue(&result.User, "user", params)
	result.setStringValue(&result.Note, "note", params)

	return result
}
