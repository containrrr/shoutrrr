package pushbullet

type Config struct {
	Targets []string
	Token   string
}

var (
	minimumArguments = 2
)

func CreateConfigFromURL(url string) (*Config, error) {
	return &Config {}, nil
}