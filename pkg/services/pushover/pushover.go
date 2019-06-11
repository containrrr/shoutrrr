package pushover

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
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
type Service struct{
	standard.Standard
}

// Send a notification message to Pushover
func (service *Service) Send(url *netUrl.URL, message string, params *map[string]string) error {
	config := Config{}
	config.SetURL(url)
	data := netUrl.Values{}
	data.Set("device", config.Devices[0])
	data.Set("user", config.User)
	data.Set("token", config.Token)
	data.Set("message", message)
	service.Logln(data.Encode())

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
func (service *Service) GetConfig() types.ServiceConfig {
	return &Config{}
}