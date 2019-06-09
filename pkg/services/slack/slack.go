package slack

import (
    "bytes"
    "errors"
    "fmt"
    "github.com/containrrr/shoutrrr/pkg/types"
    "net/http"
    "net/url"
)

// Service sends notifications to a pre-configured channel or user
type Service struct {}

const (
    apiURL    = "https://hooks.slack.com/services"
    maxlength = 1000
)


// Send a notification message to Slack
func (plugin *Service) Send(serviceURL *url.URL, message string, opts types.ServiceOpts) error {
    config, err := CreateConfigFromURL(serviceURL)
    if err != nil {
        return err
    }
    if err := validateToken(config.Token); err != nil {
        return err
    }
    if len(message) > maxlength {
        return errors.New("message exceeds max length")
    }

    return plugin.doSend(config, message)
}

// GetConfig returns an empty ServiceConfig for this Service
func (plugin *Service) GetConfig() types.ServiceConfig {
    return &Config{}
}

func (plugin *Service) doSend(config *Config, message string) error {
    url := plugin.getURL(config)
    json, _ := CreateJSONPayload(config, message)
    res, err := http.Post(url, "application/json", bytes.NewReader(json))

    if res.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to send notification to service, response status code %s", res.Status)
    }
    return err
}

func (plugin *Service) getURL(config *Config) string {
    return fmt.Sprintf(
        "%s/%s/%s/%s",
        apiURL,
        config.Token.A,
        config.Token.B,
        config.Token.C)
}