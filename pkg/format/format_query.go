package format

import (
	"net/url"
	"strings"

	t "github.com/containrrr/shoutrrr/pkg/types"
)

// BuildQuery converts the fields of a config object to a delimited query string
func BuildQuery(cqr t.ConfigQueryResolver) string {
	return BuildQueryWithCustomFields(cqr, url.Values{}).Encode()
}

// BuildQueryWithCustomFields converts the fields of a config object to a delimited query string,
// escaping any custom fields that share the same key as a config prop using a "__" prefix
func BuildQueryWithCustomFields(cqr t.ConfigQueryResolver, query url.Values) url.Values {
	fields := cqr.QueryFields()
	skipEscape := len(query) < 1

	pkr, isPkr := cqr.(*PropKeyResolver)

	for _, key := range fields {
		if !skipEscape {
			// Escape any webhook query keys using the same name as service props
			escValues := query[key]
			if len(escValues) > 0 {
				query.Del(key)
				query[EscapeKey(key)] = escValues
			}
		}

		if isPkr && !pkr.KeyIsPrimary(key) {
			continue
		}
		value, err := cqr.Get(key)

		if err != nil || isPkr && pkr.IsDefault(key, value) {
			continue
		}

		query.Set(key, value)
	}

	return query
}

// SetConfigPropsFromQuery iterates over all the config prop keys and sets the config prop to the corresponding
// query value based on the key.
// SetConfigPropsFromQuery returns a non-nil url.Values query with all config prop keys removed, even if any of
// them could not be used to set a config field, and with any escaped keys unescaped.
// The error returned is the first error that occurred, subsequent errors are just discarded.
func SetConfigPropsFromQuery(cqr t.ConfigQueryResolver, query url.Values) (url.Values, error) {
	var firstError error
	if query == nil {
		return url.Values{}, nil
	}
	for _, key := range cqr.QueryFields() {
		// Retrieve the service-related prop value
		values := query[key]
		if len(values) > 0 {
			if err := cqr.Set(key, values[0]); err != nil && firstError == nil {
				firstError = err
			}
		}
		// Remove it from the query Values
		query.Del(key)

		// If an escaped version of the key exist, unescape it
		escKey := EscapeKey(key)
		escValues := query[escKey]
		if len(escValues) > 0 {
			query.Del(escKey)
			query[key] = escValues
		}
	}
	return query, firstError
}

// EscapeKey adds the KeyPrefix to custom URL query keys that conflict with service config prop keys
func EscapeKey(key string) string {
	return KeyPrefix + key
}

// UnescapeKey removes the KeyPrefix from custom URL query keys that conflict with service config prop keys
func UnescapeKey(key string) string {
	return strings.TrimPrefix(key, KeyPrefix)
}

// KeyPrefix is the prefix prepended to custom URL query keys that conflict with service config prop keys,
// consisting of two underscore characters ("__")
const KeyPrefix = "__"
