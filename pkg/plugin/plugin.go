package plugin

import (
	"fmt"
	"log"
	"net/url"
)

type Plugin interface {
	Send(serviceUrl url.URL, message string, opts PluginOpts) error
	GetConfig() PluginConfig
}

type PluginConfig interface {
	Get(string) (string, error)
	Set(string, string) error
	QueryFields() []string
	GetURL() url.URL
	SetURL(url.URL) error
	Enums() map[string]EnumFormatter
}

func FormatQuery(c PluginConfig) string {
	query := ""
	fields := c.QueryFields()
	format := "%s=%s"
	for index, key := range fields {
		value, _ := c.Get(key)
		if index == 1 {
			format = "&%s=%s"
		}
		query += fmt.Sprintf(format, key, value)
	}

	return query
}

func SetConfigQuery(c PluginConfig, values url.Values) (PluginConfig, error)  {

	for key, vals := range values {

		value := vals[0]

		if len(vals) > 1 {
			log.Printf("warning: %s additional value ignored!: %s\n", key, vals[1])
		}

		if err := c.Set(key, value); err != nil {
			return nil, err
		}

		if true {
			log.Printf("Query \"%s\" => \"%s\"\n", key, value)
		}

	}
	return c, nil
}