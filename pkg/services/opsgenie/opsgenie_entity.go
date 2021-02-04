package opsgenie

import (
	"fmt"
	"regexp"
	"strings"
)

// Entity represents either a user or a team
//
// The different variations are:
//
// { "id":"4513b7ea-3b91-438f-b7e4-e3e54af9147c", "type":"team" }
// { "name":"rocket_team", "type":"team" }
// { "id":"bb4d9938-c3c2-455d-aaab-727aa701c0d8", "type":"user" }
// { "username":"trinity@opsgenie.com", "type":"user" }
type Entity struct {
	Type     string `json:"type"`
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Username string `json:"username,omitempty"`
}

// SetFromProp deserializes an entity
func (e *Entity) SetFromProp(propValue string) error {
	elements := strings.Split(propValue, ":")

	if len(elements) != 2 {
		return fmt.Errorf("invalid entity, should have two elments separated by colon: %q", propValue)
	}
	e.Type = elements[0]
	identifier := elements[1]
	isID, err := isOpsGenieID(identifier)
	if err != nil {
		return fmt.Errorf("invalid entity, cannot parse id/name: %q", identifier)
	}

	if isID {
		e.ID = identifier
	} else if e.Type == "team" {
		e.Name = identifier
	} else if e.Type == "user" {
		e.Username = identifier
	} else {
		return fmt.Errorf("invalid entity, unexpected entity type: %q", e.Type)
	}

	return nil
}

// GetPropValue serializes an entity
func (e *Entity) GetPropValue() (string, error) {
	identifier := ""

	if e.ID != "" {
		identifier = e.ID
	} else if e.Name != "" {
		identifier = e.Name
	} else if e.Username != "" {
		identifier = e.Username
	} else {
		return "", fmt.Errorf("invalid entity, should have either ID, name or username")
	}

	return fmt.Sprintf("%s:%s", e.Type, identifier), nil
}

func serializeEntities(entities []Entity) (string, error) {
	entityStrings := []string{}

	for _, entity := range entities {
		entityStr, err := entity.GetPropValue()
		if err != nil {
			return "", err
		}
		entityStrings = append(entityStrings, entityStr)
	}

	return strings.Join(entityStrings, ","), nil
}

func deserializeEntities(str string) ([]Entity, error) {
	result := []Entity{}

	entityStrings := strings.Split(str, ",")
	for _, entityStr := range entityStrings {
		entity := Entity{}
		if err := entity.SetFromProp(entityStr); err != nil {
			return result, err
		}
		result = append(result, entity)
	}

	return result, nil
}

// Detects OpsGenie IDs in the form 4513b7ea-3b91-438f-b7e4-e3e54af9147c
func isOpsGenieID(str string) (bool, error) {
	return regexp.MatchString(`^[0-9a-f]{8}\-[0-9a-f]{4}\-[0-9a-f]{4}\-[0-9a-f]{4}\-[0-9a-f]{12}$`, str)
}
