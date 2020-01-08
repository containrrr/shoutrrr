package pushbullet

import (
	"encoding/json"

	"github.com/containrrr/shoutrrr/pkg/types"
)

// JSON used within the Slack service
type JSON struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Body  string `json:"body"`

	Email      string `json:"email"`
	ChannelTag string `json:"channel_tag"`
	DeviceIden string `json:"device_iden"`
}

// CreateJSONPayload compatible with the slack webhook api
func CreateJSONPayload(target string, targetType TargetType, config *Config, message string, params *types.Params) ([]byte, error) {
	baseMessage := JSON{
		Type:  "note",
		Title: getTitle(params),
		Body:  message,
	}

	switch targetType {
	case EmailTarget:
		return CreateEmailPayload(config, target, baseMessage)
	case ChannelTarget:
		return CreateChannelPayload(config, target, baseMessage)
	case DeviceTarget:
		return CreateDevicePayload(config, target, baseMessage)
	}
	return json.Marshal(baseMessage)
}

//CreateChannelPayload from a base message
func CreateChannelPayload(config *Config, target string, partialPayload JSON) ([]byte, error) {
	partialPayload.ChannelTag = target[1:]
	return json.Marshal(partialPayload)
}

//CreateDevicePayload from a base message
func CreateDevicePayload(config *Config, target string, partialPayload JSON) ([]byte, error) {
	partialPayload.DeviceIden = target
	return json.Marshal(partialPayload)
}

//CreateEmailPayload from a base message
func CreateEmailPayload(config *Config, target string, partialPayload JSON) ([]byte, error) {
	partialPayload.Email = target
	return json.Marshal(partialPayload)
}
