package slack

import (
    NotifyFormat "github.com/containrrr/shoutrrr/pkg/format"
)

type SlackPlugin struct {}

type SlackConfig struct {
    Botname string
    Token SlackToken
}

const (
    name = "Slack"
    serviceUrl = "https://slack.com/"
    url = "https://hooks.slack.com/services"
    maxlength = 1000
    format = NotifyFormat.Markdown
)

type SlackErrorMessage string
const (
    TokenAMissing SlackErrorMessage = "First part of the API token is missing."
    TokenBMissing SlackErrorMessage = "Second part of the API token is missing."
    TokenCMissing SlackErrorMessage = "Third part of the API token is missing."
    TokenAMalformed SlackErrorMessage = "First part of the API token is malformed."
    TokenBMalformed SlackErrorMessage = "Second part of the API token is malformed."
    TokenCMalformed SlackErrorMessage = "Third part of the API token is malformed."
)

func (slack *SlackPlugin) Send(config SlackConfig, message string) error {

    if err := validateToken(config.Token); err != nil {
        return err
    }

    return nil
}
