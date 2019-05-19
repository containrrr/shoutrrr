package slack

import (
	"errors"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/plugins"
	)

// Config for the slack plugin
type Config struct {
	BotName string
	Token   Token
}

const (
	// DefaultUser for sending notifications to slack
	DefaultUser = "Shoutrrr"
)

// CreateConfigFromURL to use within the slack plugin
func CreateConfigFromURL(url string ) (*Config, error) {
	arguments, err := plugins.ExtractArguments(url)
	if err != nil {
		return nil, err
	}
	if len(arguments) < 3 {
		fmt.Println(arguments)
		return nil, errors.New(string(NotEnoughArguments))
	}

	if len(arguments) < 4 {
		return &Config{
			BotName: DefaultUser,
			Token: Token{
				A: arguments[0],
				B: arguments[1],
				C: arguments[2],
			},
		}, nil
	}

	return &Config{
		BotName: arguments[0],
		Token: Token{
			A: arguments[1],
			B: arguments[2],
			C: arguments[3],
		},
	}, nil
}