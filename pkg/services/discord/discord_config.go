package discord

import (
	"errors"
	"net/url"

	"github.com/nicholas-fedor/shoutrrr/pkg/format"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/standard"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
)

// Config is the configuration needed to send discord notifications.
type Config struct {
	standard.EnumlessConfig
	WebhookID  string `url:"host"`
	Token      string `url:"user"`
	Title      string `default:""    key:"title"`
	Username   string `default:""    desc:"Override the webhook default username"                                                            key:"username"`
	Avatar     string `default:""    desc:"Override the webhook default avatar with specified URL"                                           key:"avatar,avatarurl"`
	Color      uint   `base:"16"     default:"0x50D9ff"                                                                                      desc:"The color of the left border for plain messages"   key:"color"`
	ColorError uint   `base:"16"     default:"0xd60510"                                                                                      desc:"The color of the left border for error messages"   key:"colorError"`
	ColorWarn  uint   `base:"16"     default:"0xffc441"                                                                                      desc:"The color of the left border for warning messages" key:"colorWarn"`
	ColorInfo  uint   `base:"16"     default:"0x2488ff"                                                                                      desc:"The color of the left border for info messages"    key:"colorInfo"`
	ColorDebug uint   `base:"16"     default:"0x7b00ab"                                                                                      desc:"The color of the left border for debug messages"   key:"colorDebug"`
	SplitLines bool   `default:"Yes" desc:"Whether to send each line as a separate embedded item"                                            key:"splitLines"`
	JSON       bool   `default:"No"  desc:"Whether to send the whole message as the JSON payload instead of using it as the 'content' field" key:"json"`
}

// LevelColors returns an array of colors with a MessageLevel index.
func (config *Config) LevelColors() (colors [types.MessageLevelCount]uint) {
	colors[types.Unknown] = config.Color
	colors[types.Error] = config.ColorError
	colors[types.Warning] = config.ColorWarn
	colors[types.Info] = config.ColorInfo
	colors[types.Debug] = config.ColorDebug

	return colors
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

func (config *Config) getURL(resolver types.ConfigQueryResolver) (u *url.URL) {
	u = &url.URL{
		User:       url.User(config.Token),
		Host:       config.WebhookID,
		Scheme:     Scheme,
		RawQuery:   format.BuildQuery(resolver),
		ForceQuery: false,
	}

	if config.JSON {
		u.Path = "/raw"
	}

	return u
}

// SetURL updates a ServiceConfig from a URL representation of it's field values.
func (config *Config) setURL(resolver types.ConfigQueryResolver, url *url.URL) error {
	config.WebhookID = url.Host
	config.Token = url.User.Username()

	if len(url.Path) > 0 {
		switch url.Path {
		case "/raw":
			config.JSON = true

			break
		default:
			return errors.New("illegal argument in config URL")
		}
	}

	if config.WebhookID == "" {
		return errors.New("webhook ID missing from config URL")
	}

	if len(config.Token) < 1 {
		return errors.New("token missing from config URL")
	}

	for key, vals := range url.Query() {
		if err := resolver.Set(key, vals[0]); err != nil {
			return err
		}
	}

	return nil
}

// Scheme is the identifying part of this service's configuration URL.
const Scheme = "discord"
