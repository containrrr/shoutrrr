package slack

import (
	"errors"
	"regexp"
	"strings"
)

type SlackConfig struct {
	Botname string
	Token   Token
}

const (
	defaultUser = "Shoutrrr"
)

func CreateConfigFromUrl(url string ) (*SlackConfig, error) {
	arguments, err := ExtractArguments(url)
	if err != nil {
		return nil, err
	}
	if len(arguments) < 3 {
		return nil, errors.New(string(NotEnoughArguments))
	}

	if len(arguments) < 4 {
		return &SlackConfig{
			Botname: defaultUser,
			Token: Token{
				A: arguments[0],
				B: arguments[1],
				C: arguments[2],
			},
		}, nil
	}

	return &SlackConfig{
		Botname: arguments[0],
		Token: Token{
			A: arguments[1],
			B: arguments[2],
			C: arguments[3],
		},
	}, nil
}

func ExtractArguments(url string) ([]string, error) {
	regex, err := regexp.Compile("^[a-zA-Z]+://(.*)$")
	if err != nil {
		return nil, errors.New("could not compile regex")
	}
	match := regex.FindStringSubmatch(url)
	if len(match[1]) <= 0 {
		return nil, errors.New("could not extract any arguments")
	}
	return strings.Split(match[1], "/"), nil
}