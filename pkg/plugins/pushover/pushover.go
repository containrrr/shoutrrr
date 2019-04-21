package pushover

import (
	"errors"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/plugins"
	"net/http"
	netUrl "net/url"
	"strings"
)

type PushoverConfig struct {
	Token   string
	User    string
	Devices []string
}

const (
	hookUrl     = "https://api.pushover.net/1/messages.json"
	contentType = "application/x-www-form-urlencoded"
)

type PushoverPlugin struct{}

func (plugin *PushoverPlugin) Send(url string, message string) error {
	config, _ := CreateConfigFromUrl(url)
	data := netUrl.Values{}
	data.Set("device", config.Devices[0])
	data.Set("user", config.User)
	data.Set("token", config.Token)
	data.Set("message", message)
	fmt.Println(data.Encode())

	res, err := http.Post(
		hookUrl,
		contentType,
		strings.NewReader(data.Encode()))
	if res.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("failed to send notification to pushover, response status code %s", res.Status))
	}
	return err
}

func CreateConfigFromUrl(url string) (PushoverConfig, error) {
	args, err := plugins.ExtractArguments(url)
	if err != nil {
		return PushoverConfig{}, err
	}
	if len(args) < 2 {
		return PushoverConfig{}, errors.New("the minimum amount of arguments for pushover is 2")
	}
	return PushoverConfig{
		Token:   args[0],
		User:    args[1],
		Devices: args[2:],
	}, nil
}
