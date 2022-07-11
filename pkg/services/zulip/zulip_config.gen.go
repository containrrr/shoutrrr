// Code generated by "shoutrrr-gen "; DO NOT EDIT.
package zulip

import (
	"fmt"
	"net/url"
	_ "strings"

	"github.com/containrrr/shoutrrr/pkg/conf"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// (‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾)
//  )  Props                          (
// (___________________________________)

type Config struct {
	BotKey  string `url:"password" `
	BotMail string `url:"user" `
	Host    string `url:"host,port" `
	Stream  string `key:"stream" `
	Topic   string `key:"topic,title" `
}

type configProp int

const (
	propBotKey  configProp = 0
	propBotMail configProp = 1
	propHost    configProp = 2
	propStream  configProp = 3
	propTopic   configProp = 4
	propCount              = 5
)

var propInfo = types.ConfigPropInfo{
	PropNames: []string{
		"BotKey",
		"BotMail",
		"Host",
		"Stream",
		"Topic",
	},

	// Note that propKeys may not align with propNames, as a property can have no or multiple keys
	Keys: []string{
		"stream",
		"title",
		"topic",
	},

	DefaultValues: []string{
		"",
		"",
		"",
		"",
		"",
	},

	PrimaryKeys: []int{
		-1,
		-1,
		-1,
		0,
		2,
	},

	KeyPropIndexes: map[string]int{
		"stream": 3,
		"title":  4,
		"topic":  4,
	},
}

func (_ *Config) PropInfo() *types.ConfigPropInfo {
	return &propInfo
}

// (‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾)
//  )  GetURL                         (
// (___________________________________)

// GetURL returns a URL representation of it's current field values
func (config *Config) GetURL() *url.URL {
	return &url.URL{
		User:     conf.UserInfoOrNil(url.UserPassword(config.BotMail, config.BotKey)),
		Host:     config.Host,
		Path:     "",
		RawQuery: conf.QueryValues(config).Encode(),
		Scheme:   Scheme,
	}
}

// (‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾)
//  )  SetURL                         (
// (___________________________________)

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *Config) SetURL(configURL *url.URL) error {
	if lc, ok := (interface{})(config).(types.ConfigWithLegacyURLSupport); ok {
		configURL = lc.UpdateLegacyURL(configURL)
	}
	updates := make(map[int]string, propCount)
	updates[int(propHost)] = configURL.Host
	if pwd, found := configURL.User.Password(); found {
		updates[int(propBotKey)] = pwd
	}
	updates[int(propBotMail)] = configURL.User.Username()
	if configURL.Path != "" && configURL.Path != "/" {
		return fmt.Errorf("unexpected path in config URL: %v", configURL.Path)
	}

	for key, value := range configURL.Query() {

		if propIndex, found := propInfo.PropIndexFor(key); found {
			updates[propIndex] = value[0]
		} else if key != "title" {
			return fmt.Errorf("invalid key %q", key)
		}
	}

	err := config.Update(updates)
	if err != nil {
		return err
	}

	if config.BotKey == "" {
		return fmt.Errorf("botKey missing from config URL")
	}

	if config.BotMail == "" {
		return fmt.Errorf("botMail missing from config URL")
	}

	if config.Host == "" {
		return fmt.Errorf("host missing from config URL")
	}

	return nil
}

// (‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾)
//  )  Enums / Options                (
// (___________________________________)

func (config *Config) Enums() map[string]types.EnumFormatter {
	return map[string]types.EnumFormatter{}
}

// Update updates the Config from a map of it's properties
func (config *Config) Update(updates map[int]string) error {
	var last_err error
	for index, value := range updates {
		switch configProp(index) {
		case propBotKey:
			if val, err := conf.ParseTextValue(value); err != nil {
				last_err = err
			} else {
				config.BotKey = val
			}
		case propBotMail:
			if val, err := conf.ParseTextValue(value); err != nil {
				last_err = err
			} else {
				config.BotMail = val
			}
		case propHost:
			if val, err := conf.ParseTextValue(value); err != nil {
				last_err = err
			} else {
				config.Host = val
			}
		case propStream:
			if val, err := conf.ParseTextValue(value); err != nil {
				last_err = err
			} else {
				config.Stream = val
			}
		case propTopic:
			if val, err := conf.ParseTextValue(value); err != nil {
				last_err = err
			} else {
				config.Topic = val
			}
		default:
			return fmt.Errorf("invalid key")
		}
		if last_err != nil {
			return fmt.Errorf("failed to set value for %v: %v", propInfo.PropNames[index], last_err)
		}
	}
	return nil
}

// Update updates the Config from a map of it's properties
func (config *Config) PropValue(prop int) string {
	switch configProp(prop) {
	case propBotKey:
		return conf.FormatTextValue(config.BotKey)
	case propBotMail:
		return conf.FormatTextValue(config.BotMail)
	case propHost:
		return conf.FormatTextValue(config.Host)
	case propStream:
		return conf.FormatTextValue(config.Stream)
	case propTopic:
		return conf.FormatTextValue(config.Topic)
	default:
		return ""
	}
}
