package rocketchat

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"

	"github.com/nicholas-fedor/shoutrrr/pkg/services/standard"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
)

// Service sends notifications to a pre-configured channel or user.
type Service struct {
	standard.Standard
	Config *Config
	Client *http.Client // Add this field
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service.
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)

	service.Config = &Config{}
	if service.Client == nil {
		service.Client = http.DefaultClient // Default to standard client if not set
	}

	if err := service.Config.SetURL(configURL); err != nil {
		return err
	}

	return nil
}

// GetID returns the service identifier.
func (service *Service) GetID() string {
	return Scheme
}

// Send a notification message to Rocket.chat.
func (service *Service) Send(message string, params *types.Params) error {
	var res *http.Response

	var err error

	config := service.Config
	apiURL := buildURL(config)
	json, _ := CreateJSONPayload(config, message, params)

	res, err = service.Client.Post(apiURL, "application/json", bytes.NewReader(json))
	if err != nil {
		return fmt.Errorf("error while posting to URL: %w\nHOST: %s\nPORT: %s", err, config.Host, config.Port)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		resBody, _ := io.ReadAll(res.Body)

		return fmt.Errorf("notification failed: %d %s", res.StatusCode, resBody)
	}

	return nil
}

func buildURL(config *Config) string {
	base := config.Host
	if config.Port != "" {
		base = net.JoinHostPort(config.Host, config.Port)
	}

	return fmt.Sprintf("https://%s/hooks/%s/%s", base, config.TokenA, config.TokenB)
}
