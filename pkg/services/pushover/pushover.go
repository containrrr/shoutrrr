package pushover

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
	"net/http"
	"net/url"
	"strings"
)


const (
	hookURL     = "https://api.pushover.net/1/messages.json"
	contentType = "application/x-www-form-urlencoded"
)

// Service providing the notification service Pushover
type Service struct{
	standard.Standard
	configURL *url.URL
}

// Send a notification message to Pushover
func (service *Service) Send(message string, params *map[string]string) error {
	config := Config{}
	config.SetURL(service.configURL)
	data := url.Values{}
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

// NewConfig returns an empty ServiceConfig for this Service
func (service *Service) NewConfig() types.ServiceConfig {
	return &Config{}
}