package telegram

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/nicholas-fedor/shoutrrr/pkg/format"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
)

// Config for use within the telegram plugin.
type Config struct {
	Token        string    `url:"user"`
	Preview      bool      `default:"Yes"                                          desc:"If disabled, no web page preview will be displayed for URLs" key:"preview"`
	Notification bool      `default:"Yes"                                          desc:"If disabled, sends Message silently"                         key:"notification"`
	ParseMode    parseMode `default:"None"                                         desc:"How the text Message should be parsed"                       key:"parsemode"`
	Chats        []string  `desc:"Chat IDs or Channel names (using @channel-name)" key:"chats,channels"`
	Title        string    `default:""                                             desc:"Notification title, optionally set by the sender"            key:"title"`
}

// Enums returns the fields that should use a corresponding EnumFormatter to Print/Parse their values.
func (config *Config) Enums() map[string]types.EnumFormatter {
	return map[string]types.EnumFormatter{
		"ParseMode": ParseModes.Enum,
	}
}

// GetURL returns a URL representation of it's current field values.
func (config *Config) GetURL() *url.URL {
	resolver := format.NewPropKeyResolver(config)

	return config.getURL(&resolver)
}

// SetURL updates a ServiceConfig from a URL representation of it's field values.
func (config *Config) SetURL(url *url.URL) error {
	resolver := format.NewPropKeyResolver(config)

	return config.setURL(&resolver, url)
}

func (config *Config) getURL(resolver types.ConfigQueryResolver) *url.URL {
	tokenParts := strings.Split(config.Token, ":")

	return &url.URL{
		User:       url.UserPassword(tokenParts[0], tokenParts[1]),
		Host:       Scheme,
		Scheme:     Scheme,
		ForceQuery: true,
		RawQuery:   format.BuildQuery(resolver),
	}
}

func (config *Config) setURL(resolver types.ConfigQueryResolver, url *url.URL) error {
	password, _ := url.User.Password()

	token := url.User.Username() + ":" + password
	if url.String() != "telegram://dummy@dummy.com" {
		if !IsTokenValid(token) {
			return fmt.Errorf("invalid telegram token %s", token)
		}
	}

	for key, vals := range url.Query() {
		if err := resolver.Set(key, vals[0]); err != nil {
			return err
		}
	}

	if url.String() != "telegram://dummy@dummy.com" {
		if len(config.Chats) < 1 {
			return errors.New("no channels defined in config URL")
		}
	}

	config.Token = token

	return nil
}

// Scheme is the identifying part of this service's configuration URL.
const (
	Scheme = "telegram"
)
