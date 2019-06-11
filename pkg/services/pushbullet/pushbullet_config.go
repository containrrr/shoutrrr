package pushbullet

import "github.com/containrrr/shoutrrr/pkg/services/standard"

// Config ...
type Config struct {
	standard.QuerylessConfig
	Targets []string
	Token   string
}

var (
	minimumArguments = 2
)

// CreateConfigFromURL ...
func CreateConfigFromURL(url string) (*Config, error) {
	return &Config {}, nil
}