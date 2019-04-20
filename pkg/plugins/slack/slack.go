package slack

import (
    "errors"
    "fmt"
    "net/http"
    "strings"
)

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
    slack.getUrl(config)

    return nil
}

func (slack *SlackPlugin) doSend(config *SlackConfig, message string) error {
    url := slack.getUrl(config)
    _, err := http.Post(url, "application/json", strings.NewReader(message))
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