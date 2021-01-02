package format

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// BuildQuery converts the fields of a config object to a delimited query string
func BuildQuery(cqr types.ConfigQueryResolver) string {
	query := ""
	format := "%s=%s"
	fields := cqr.QueryFields()
	for index, key := range fields {
		value, _ := cqr.Get(key)
		if index == 1 {
			format = "&%s=%s"
		}
		query += fmt.Sprintf(format, key, value)
	}

	return query
}
