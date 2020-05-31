package pushover

import (
	"errors"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
	"net/url"
	"strconv"
	"strings"
)

// Config for the Pushover notification service service
type Config struct {
	Token    string
	User     string
	Devices  []string
	Priority int8
	Title    string
}

// QueryFields returns the fields that are part of the Query of the service URL
func (config *Config) QueryFields() []string {
	return []string{
		"devices",
		"priority",
		"title",
	}
}

// Enums returns the fields that should use a corresponding EnumFormatter to Print/Parse their values
func (config *Config) Enums() map[string]types.EnumFormatter {
	return map[string]types.EnumFormatter{}
}

// Get returns the value of a Query field
func (config *Config) Get(key string) (string, error) {
	switch key {
	case "devices":
		return strings.Join(config.Devices, ","), nil
	case "priority":
		return strconv.FormatInt(int64(config.Priority), 10), nil
	case "title":
		return config.Title, nil
	}

	return "", fmt.Errorf("invalid query key \"%s\"", key)
}

// Set updates the value of a Query field
func (config *Config) Set(key string, value string) error {
	switch key {
	case "devices":
		config.Devices = strings.Split(value, ",")
	case "priority":
		priority, err := strconv.ParseInt(value, 10, 8)
		if err == nil {
			config.Priority = int8(priority)
		}
		return err
	case "title":
		config.Title = value
	default:
		return fmt.Errorf("invalid query key \"%s\"", key)
	}
	return nil
}

// GetURL returns a URL representation of it's current field values
func (config *Config) GetURL() *url.URL {

	return &url.URL{
		User:       url.UserPassword("Token", config.Token),
		Host:       config.User,
		Scheme:     Scheme,
		ForceQuery: true,
		RawQuery:   format.BuildQuery(config),
	}

}

// SetURL updates a ServiceConfig from a URL representation of it's field values
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
		return errors.New(string(UserMissing))
	}

	if len(config.Token) < 1 {
		return errors.New(string(TokenMissing))
	}

	return nil
}

// Scheme is the identifying part of this service's configuration URL
const Scheme = "pushover"
