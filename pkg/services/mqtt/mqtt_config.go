package mqtt

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
	"net/url"
	"strings"
	"strconv"
)

// Config for use within the mqtt
type Config struct {
	Host       	 string    `key:"broker" default:"" desc:"MQTT broker server hostname or IP address"`
	Port         int64       `key:"port" default:"1883" desc:"TCP Port"`
	Topic        string    `key:"topic" default:"" desc:"Topic where the message is sent"`
	DisableTLS   bool      `key:"disabletls" default:"Yes"`
	ParseMode    parseMode `key:"parsemode" default:"None" desc:"How the text message should be parsed"`
}

// Enums returns the fields that should use a corresponding EnumFormatter to Print/Parse their values
func (config *Config) Enums() map[string]types.EnumFormatter {
	return map[string]types.EnumFormatter{
		"ParseMode": parseModes.Enum,
	}
}

// GetURL returns a URL representation of it's current field values
func (config *Config) GetURL() *url.URL {
	resolver := format.NewPropKeyResolver(config)
	return config.getURL(&resolver)
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *Config) SetURL(url *url.URL) error {
	resolver := format.NewPropKeyResolver(config)
	return config.setURL(&resolver, url)
}

func (config *Config) getURL(resolver types.ConfigQueryResolver) *url.URL {
	
	return &url.URL{
		Host:       fmt.Sprintf("%s:%d", config.Host, config.Port),
		Scheme:     Scheme,
		ForceQuery: true,
		RawQuery:   format.BuildQuery(resolver),
	}

}

func Split(r rune) bool {
    return r == '/' || r == '?' || r == ':'
}

func getTopic(r rune) bool {
    return r == '='
}

func (config *Config) setURL(resolver types.ConfigQueryResolver, url *url.URL) error {

	u := strings.FieldsFunc(url.String(), Split)
	topic := strings.FieldsFunc(url.String(), getTopic)

	port, err := strconv.ParseInt(u[2], 10, 64)

	if err != nil {
		return err
	}

	if len(u) > 4 {
		config.Host = u[1]
		config.Port = port
		config.Topic = topic[1]
	}

	return nil
}

// Scheme is the identifying part of this service's configuration URL
const (
	Scheme = "tcp"
)
