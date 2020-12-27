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

	pkr, isPkr := cqr.(*PropKeyResolver)

	for index, key := range fields {
		if isPkr && !pkr.KeyIsPrimary(key) {
			continue
		}
		value, _ := cqr.Get(key)
		if index == 1 {
			format = "&%s=%s"
		}
		query += fmt.Sprintf(format, key, value)
	}

	return query
}
