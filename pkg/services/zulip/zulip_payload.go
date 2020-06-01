package zulip

import (
	"net/url"
)

// CreatePayload compatible with the zulip api
func CreatePayload(config *Config, message string) url.Values {
	form := url.Values{}
	form.Set("type", "stream")
	form.Set("to", config.Stream)
	form.Set("content", message)

	if config.Topic != "" {
		form.Set("topic", config.Topic)
	}

	return form
}
