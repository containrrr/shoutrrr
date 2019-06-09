package pushover

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/types"
	"net/http"
	netUrl "net/url"
	"strings"
)


const (
	hookURL     = "https://api.pushover.net/1/messages.json"
	contentType = "application/x-www-form-urlencoded"
)

// Service providing the notification service Pushover
type Service struct{}

// Send a notification message to Pushover
func (plugin *Service) Send(url *netUrl.URL, message string, opts types.ServiceOpts) error {
	config := Config{}
	config.SetURL(url)
	data := netUrl.Values{}
	data.Set("device", config.Devices[0])
	data.Set("user", config.User)
	data.Set("token", config.Token)
	data.Set("message", message)
	opts.Logger().Println(data.Encode())

	res, err := http.Post(
		hookURL,
		contentType,
		strings.NewReader(data.Encode()))
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send notification to pushover, response status code %s", res.Status)
	}
	return err
}

// GetConfig returns an empty ServiceConfig for this Service
func (plugin *Service) GetConfig() types.ServiceConfig {
	return &Config{}
}