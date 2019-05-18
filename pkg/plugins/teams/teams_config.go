package teams

import (
	"errors"
	. "github.com/containrrr/shoutrrr/pkg/plugins"
)

type TeamsConfig struct {
	Token TeamsToken
}

func (plugin *TeamsPlugin) CreateConfigFromURL(url string) (*TeamsConfig, error) {
	arguments, err := ExtractArguments(url);
	if err != nil {
		return nil, err
	} else if !isTokenValid(arguments) {
		return nil, errors.New("invalid service url. malformed tokens")
	}
	return &TeamsConfig{
		Token:TeamsToken{
			A: arguments[0],
			B: arguments[1],
			C: arguments[2],
		},
	}, nil
}