package pushover

import (
	"errors"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/plugins"
	"net/http"
	netUrl "net/url"
	"strings"
)

// Config for the Pushover notification service plugin
type Config struct {
	Token   string
	User    string
	Devices []string
}

const (
	hookURL     = "https://api.pushover.net/1/messages.json"
	contentType = "application/x-www-form-urlencoded"
)

// Plugin providing the notification service Pushover
type Plugin struct{}

// Send a notification message to Pushover
func (plugin *Plugin) Send(url string, message string) error {
	config, _ := CreateConfigFromURL(url)
	data := netUrl.Values{}
	data.Set("device", config.Devices[0])
	data.Set("user", config.User)
	data.Set("token", config.Token)
	data.Set("message", message)
	fmt.Println(data.Encode())

	res, err := http.Post(
		hookURL,
		contentType,
		strings.NewReader(data.Encode()))
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send notification to pushover, response status code %s", res.Status)
	}
	return err
}

// CreateConfigFromURL to be used with the Pushover notification service plugin
func CreateConfigFromURL(url string) (Config, error) {
	args, err := plugins.ExtractArguments(url)
	if err != nil {
		return Config{}, err
	}
	if len(args) < 2 {
		return Config{}, errors.New("the minimum amount of arguments for pushover is 2")
	}
	return Config{
		Token:   args[0],
		User:    args[1],
		Devices: args[2:],
	}, nil
}
