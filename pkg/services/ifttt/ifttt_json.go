package ifttt

import (
	"encoding/json"

	"github.com/nicholas-fedor/shoutrrr/pkg/types"
)

const (
	ValueFieldOne   = 1 // Represents Value1 field
	ValueFieldTwo   = 2 // Represents Value2 field
	ValueFieldThree = 3 // Represents Value3 field
)

// jsonPayload is the actual notification payload.
type jsonPayload struct {
	Value1 string `json:"value1"`
	Value2 string `json:"value2"`
	Value3 string `json:"value3"`
}

// createJSONToSend creates a jsonPayload payload to be sent to the IFTTT webhook API.
func createJSONToSend(config *Config, message string, params *types.Params) ([]byte, error) {
	payload := jsonPayload{
		Value1: config.Value1,
		Value2: config.Value2,
		Value3: config.Value3,
	}

	if params != nil {
		if value, found := (*params)["value1"]; found {
			payload.Value1 = value
		}

		if value, found := (*params)["value2"]; found {
			payload.Value2 = value
		}

		if value, found := (*params)["value3"]; found {
			payload.Value3 = value
		}
	}

	switch config.UseMessageAsValue {
	case ValueFieldOne:
		payload.Value1 = message
	case ValueFieldTwo:
		payload.Value2 = message
	case ValueFieldThree:
		payload.Value3 = message
	}

	return json.Marshal(payload)
}
