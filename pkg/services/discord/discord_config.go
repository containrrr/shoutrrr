//go:generate go run ../../../cmd/shoutrrr-gen --lang go
package discord

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/pkr"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// LegacyConfig is the configuration needed to send discord notifications
type LegacyConfig struct {
	standard.EnumlessConfig
	WebhookID  string `url:"host"`
	Token      string `url:"user"`
	Title      string `key:"title"      default:""`
	Username   string `key:"username"   default:""         desc:"Override the webhook default username"`
	Avatar     string `key:"avatar,avatarurl"     default:""         desc:"Override the webhook default avatar with specified URL"`
	Color      uint   `key:"color"      default:"0x50D9ff" desc:"The color of the left border for plain messages"   base:"16"`
	ColorError uint   `key:"colorError" default:"0xd60510" desc:"The color of the left border for error messages"   base:"16"`
	ColorWarn  uint   `key:"colorWarn"  default:"0xffc441" desc:"The color of the left border for warning messages" base:"16"`
	ColorInfo  uint   `key:"colorInfo"  default:"0x2488ff" desc:"The color of the left border for info messages"    base:"16"`
	ColorDebug uint   `key:"colorDebug" default:"0x7b00ab" desc:"The color of the left border for debug messages"   base:"16"`
	SplitLines bool   `key:"splitLines" default:"Yes"      desc:"Whether to send each line as a separate embedded item"`
	JSON       bool   `key:"json"       default:"No"       desc:"Whether to send the whole message as the JSON payload instead of using it as the 'content' field"`
}

// LevelColors returns an array of colors with a MessageLevel index
func (config *Config) LevelColors() (colors [types.MessageLevelCount]uint) {
	colors[types.Unknown] = uint(config.Color)
	colors[types.Error] = uint(config.ColorError)
	colors[types.Warning] = uint(config.ColorWarn)
	colors[types.Info] = uint(config.ColorInfo)
	colors[types.Debug] = uint(config.ColorDebug)

	return colors
}

// GetURL returns a URL representation of it's current field values
func (config *LegacyConfig) GetURL() *url.URL {
	resolver := pkr.NewPropKeyResolver(config)
	return config.getURL(&resolver)
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *LegacyConfig) SetURL(url *url.URL) error {
	resolver := pkr.NewPropKeyResolver(config)
	return config.setURL(&resolver, url)
}

func (config *LegacyConfig) getURL(resolver types.ConfigQueryResolver) (u *url.URL) {
	u = &url.URL{
		User:       url.User(config.Token),
		Host:       config.WebhookID,
		Scheme:     Scheme,
		RawQuery:   pkr.BuildQuery(resolver),
		ForceQuery: false,
	}

	if config.JSON {
		u.Path = "/raw"
	}

	return u
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *LegacyConfig) setURL(resolver types.ConfigQueryResolver, url *url.URL) error {

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

func (service *Service) GetLegacyConfig() types.ServiceConfig {
	return &LegacyConfig{}
}

type rawModeType string

func (config *Config) getRawMode() string {
	if config.JSON {
		return "raw"
	} else {
		return ""
	}
}

func (config *Config) setRawMode(v string) (rawModeType, error) {
	if v == "raw" {
		config.JSON = true
		return rawModeType(v), nil
	} else if v == "" {
		return rawModeType(""), nil
	}

	return "", fmt.Errorf("invalid value raw mode value %q", v)
}

// Scheme is the identifying part of this service's configuration URL
const Scheme = "discord"
