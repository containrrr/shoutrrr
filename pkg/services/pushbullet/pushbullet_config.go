package pushbullet

// Config ...
type Config struct {
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