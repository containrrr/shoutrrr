package teams

import (
	"errors"
	"github.com/containrrr/shoutrrr/pkg/plugins"
)

// Config for use within the teams plugin
type Config struct {
	Token Token
}

// CreateConfigFromURL for use within the teams plugin
func (plugin *Plugin) CreateConfigFromURL(url string) (*Config, error) {
	arguments, err := plugins.ExtractArguments(url);
	if err != nil {
		return nil, err
	} else if !isTokenValid(arguments) {
		return nil, errors.New("invalid service url. malformed tokens")
	}
	return &Config{
		Token: Token{
			A: arguments[0],
			B: arguments[1],
			C: arguments[2],
		},
	}, nil
}