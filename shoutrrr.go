package shoutrrr

import (
	"errors"
	"github.com/containrrr/shoutrrr/pkg/router"
	"os"
)

// Send lets you send shoutrrr notifications using a supplied url and message
func Send(url string, message string) error {
	routing := router.ServiceRouter{}
	return routing.Route(url, message)
}

// SendEnv lets you send shoutrrr notifications using an url stored in your env variables and a supplied message
func SendEnv(message string) error {
	envURL := os.Getenv("SHOUTRRR_URL")
	if envURL == "" {
		return errors.New("trying to use SendEnv but SHOUTRRR_URL is not set")
	}
	return Send(envURL, message)
}