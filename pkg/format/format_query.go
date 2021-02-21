package format

import (
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/types"
)

// BuildQuery converts the fields of a config object to a delimited query string
func BuildQuery(cqr types.ConfigQueryResolver) string {
	query := url.Values{}
	fields := cqr.QueryFields()

	pkr, isPkr := cqr.(*PropKeyResolver)

	for _, key := range fields {
		if isPkr && !pkr.KeyIsPrimary(key) {
			continue
		}
		value, err := cqr.Get(key)

		if err == nil && value != "" {
			query.Set(key, value)
		}
	}

	return query.Encode()
}
