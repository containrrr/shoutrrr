package slack

import (
    "bytes"
    "errors"
    "fmt"
    "net/http"
)

// SlackPlugin sends notifications to a pre-configured channel or user
type SlackPlugin struct {}

const (
    url = "https://hooks.slack.com/services"
    maxlength = 1000
)



func (slack *SlackPlugin) Send(url string, message string) error {
    config, err := CreateConfigFromUrl(url)
    if err != nil {
        return err
    }
    if err := validateToken(config.Token); err != nil {
        return err
    }
    if len(message) > maxlength {
        return errors.New("message exceeds max length")
    }

    return slack.doSend(config, message)
}

func (slack *SlackPlugin) doSend(config *SlackConfig, message string) error {
    url := slack.getUrl(config)
    json, _ := CreateJsonPayload(config, message)
    res, err := http.Post(url, "application/json", bytes.NewReader(json))

    if res.StatusCode != http.StatusOK {
        return errors.New(fmt.Sprintf("failed to send notification to slack, response status code %s", res.Status))
    }
    return err
}

func (slack *SlackPlugin) getUrl(config *SlackConfig) string {
    return fmt.Sprintf(
        "%s/%s/%s/%s",
        url,
        config.Token.A,
        config.Token.B,
        config.Token.C)
}