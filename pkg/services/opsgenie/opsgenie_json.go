package opsgenie

import (
	"github.com/containrrr/shoutrrr/pkg/format"
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

func NewAlertPayload(message string, config *Config, params *types.Params) (AlertPayload, error) {
	if params == nil {
		params = &types.Params{}
	}

	// Defensive copy
	payloadFields := *config

	pkr := format.NewPropKeyResolver(&payloadFields)
	if value, found := (*params)["responders"]; found {
		responders, err := deserializeEntities(value)
		if err != nil {
			return AlertPayload{}, err
		}
		payloadFields.Responders = responders
		delete(*params, "responders")
	}
	if value, found := (*params)["visibleTo"]; found {
		visibleTo, err := deserializeEntities(value)
		if err != nil {
			return AlertPayload{}, err
		}
		payloadFields.VisibleTo = visibleTo
		delete(*params, "visibleTo")
	}
	if err := pkr.UpdateConfigFromParams(&payloadFields, params); err != nil {
		return AlertPayload{}, err
	}

	result := AlertPayload{
		Message:     message,
		Alias:       payloadFields.Alias,
		Description: payloadFields.Description,
		Responders:  payloadFields.Responders,
		VisibleTo:   payloadFields.VisibleTo,
		Actions:     payloadFields.Actions,
		Tags:        payloadFields.Tags,
		Details:     payloadFields.Details,
		Entity:      payloadFields.Entity,
		Source:      payloadFields.Source,
		Priority:    payloadFields.Priority,
		User:        payloadFields.User,
		Note:        payloadFields.Note,
	}
	return result, nil
}
