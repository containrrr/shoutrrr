package slack

import (
    "bytes"
    "errors"
    "fmt"
    "net/http"
)

// Plugin sends notifications to a pre-configured channel or user
type Plugin struct {}

const (
    url = "https://hooks.slack.com/services"
    maxlength = 1000
)


// Send a notification message to Slack
func (plugin *Plugin) Send(url string, message string) error {
    config, err := CreateConfigFromURL(url)
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

func (plugin *Plugin) doSend(config *Config, message string) error {
    url := plugin.getURL(config)
    json, _ := CreateJSONPayload(config, message)
    res, err := http.Post(url, "application/json", bytes.NewReader(json))

    if res.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to send notification to plugin, response status code %s", res.Status)
    }
    return err
}

func (plugin *Plugin) getURL(config *Config) string {
    return fmt.Sprintf(
        "%s/%s/%s/%s",
        url,
        config.Token.A,
        config.Token.B,
        config.Token.C)
}