package pushover

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/plugin"
	"net/http"
	netUrl "net/url"
	"strings"
)


const (
	hookURL     = "https://api.pushover.net/1/messages.json"
	contentType = "application/x-www-form-urlencoded"
)

// Plugin providing the notification service Pushover
type Plugin struct{}

// Send a notification message to Pushover
func (plugin *Plugin) Send(url netUrl.URL, message string, opts plugin.PluginOpts) error {
	config := Config{}
	config.SetURL(url)
	data := netUrl.Values{}
	data.Set("device", config.Devices[0])
	data.Set("user", config.User)
	data.Set("token", config.Token)
	data.Set("message", message)
	opts.Logger.Println(data.Encode())

	res, err := http.Post(
		hookURL,
		contentType,
		strings.NewReader(data.Encode()))
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send notification to pushover, response status code %s", res.Status)
	}
	return err
}

func (plugin *Plugin) GetConfig() plugin.PluginConfig {
	return &Config{}
}