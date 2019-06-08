package pushbullet

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/plugins"
)

type Config struct {
	Targets []string
	Token   string
}

var (
	minimumArguments = 2
)

func CreateConfigFromURL(url string) (*Config, error) {
	arguments, err := plugins.ExtractArguments(url)
	if err != nil {
		return nil, err
	}
	if len(arguments) < minimumArguments {
		return nil, fmt.Errorf("pushbullet requires %d to work", minimumArguments)
	}

	return &Config {
		Token: arguments[0],
		Targets: arguments[1:],
	}, nil
}