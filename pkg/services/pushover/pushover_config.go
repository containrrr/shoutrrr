package pushover

import (
	"errors"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
	"net/url"
	"strings"
)

// Config for the Pushover notification service service
type Config struct {
	Token   string
	User    string
	Devices []string
}

func (config *Config) QueryFields() []string {
	return []string{
		"devices",
	}
}

func (config *Config) Enums() map[string]types.EnumFormatter {
	return map[string]types.EnumFormatter{}
}

func (config *Config) Get(key string) (string, error) {
	switch key {
	case "devices":
		return strings.Join(config.Devices, ","), nil
	}
	return "", fmt.Errorf("invalid query key \"%s\"", key)
}

func (config *Config) Set(key string, value string) error {
	switch key {
	case "devices":
		config.Devices = strings.Split(value, ",")
	default:
		return fmt.Errorf("invalid query key \"%s\"", key)
	}
	return nil
}

func (config *Config) GetURL() *url.URL {

	return &url.URL{
		User:       url.UserPassword("Token", config.Token),
		Host:       config.User,
		Scheme:     Scheme,
		ForceQuery: true,
		RawQuery:   format.BuildQuery(config),
	}

}

func (config *Config) SetURL(url *url.URL) error {

	password, _ := url.User.Password()

	config.User = url.Host
	config.Token = password

	for key, vals := range url.Query() {
		if err := config.Set(key, vals[0]); err != nil {
			return err
		}
	}

	if len(config.User) < 1 {
		return errors.New("user missing from config URL")
	}

	if len(config.Token) < 1 {
		return errors.New("token missing from config URL")
	}

	return nil
}

const Scheme = "pushover"