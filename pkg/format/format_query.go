package format

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// BuildQuery converts the fields of a config object to a delimited query string
func BuildQuery(c types.ServiceConfig) string {
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
