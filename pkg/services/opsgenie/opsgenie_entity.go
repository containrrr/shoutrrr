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

func serializeEntities(entities []Entity) (string, error) {
	entityStrings := []string{}

	for _, entity := range entities {
		identifier := ""

		if entity.ID != "" {
			identifier = entity.ID
		} else if entity.Name != "" {
			identifier = entity.Name
		} else if entity.Username != "" {
			identifier = entity.Username
		} else {
			return "", fmt.Errorf("invalid entity, should have either ID, name or username")
		}

		entityStr := fmt.Sprintf("%s:%s", entity.Type, identifier)
		entityStrings = append(entityStrings, entityStr)
	}

	return strings.Join(entityStrings, ","), nil
}

func deserializeEntities(str string) ([]Entity, error) {
	result := []Entity{}

	entityStrings := strings.Split(str, ",")
	for _, entityStr := range entityStrings {
		elements := strings.Split(entityStr, ":")

		if len(elements) != 2 {
			return result, fmt.Errorf("invalid entity, should have two elments separated by colon: %q", entityStr)
		}
		entityType := elements[0]
		identifier := elements[1]

		entity := Entity{
			Type: entityType,
		}

		isID, err := isOpsGenieID(identifier)
		if err != nil {
			return result, fmt.Errorf("invalid entity, cannot parse id/name: %q", identifier)
		}

		if isID {
			entity.ID = identifier
		} else if entityType == "team" {
			entity.Name = identifier
		} else if entityType == "user" {
			entity.Username = identifier
		} else {
			return result, fmt.Errorf("invalid entity, unexpected entity type: %q", entityType)
		}

		result = append(result, entity)
	}

	return result, nil
}

// Detects OpsGenie IDs in the form 4513b7ea-3b91-438f-b7e4-e3e54af9147c
func isOpsGenieID(str string) (bool, error) {
	return regexp.MatchString(`^[0-9a-f]{8}\-[0-9a-f]{4}\-[0-9a-f]{4}\-[0-9a-f]{4}\-[0-9a-f]{12}$`, str)
}
