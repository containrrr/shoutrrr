//go:generate go run ../../../cmd/shoutrrr-gen --lang go ../../../spec/telegram.yml

package telegram

// Scheme is the identifying part of this service's configuration URL
const (
	Scheme = "telegram"
)

func (config *Config) apiHost() string {
	if config.APIHost == "" || config.APIHost == "telegram" {
		return DEFAULT_API_HOST
	}
	return config.APIHost
}

func (config *Config) token() string {
	return config.BotID + ":" + config.Token
}
