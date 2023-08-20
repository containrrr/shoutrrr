package generic

import (
	"net/url"
	"strings"
)

const extraPrefix = '$'
const headerPrefix = '@'
const caseOffset = 'a' - 'A'

func normalizedHeaderKey(key string) string {
	sb := strings.Builder{}
	sb.Grow(len(key) * 2)
	for i, c := range key {
		if 'A' <= c && c <= 'Z' {
			// Char is uppercase
			if i > 0 && key[i-1] != '-' {
				// Add missing dash
				sb.WriteRune('-')
			}
		} else if i == 0 || key[i-1] == '-' {
			// First char, or previous was dash
			c -= caseOffset
		}
		sb.WriteRune(c)
	}
	return sb.String()
}

func appendCustomQueryValues(query url.Values, headers map[string]string, extraData map[string]string) {
	for key, value := range headers {
		query.Set(string(headerPrefix)+key, value)
	}
	for key, value := range extraData {
		query.Set(string(extraPrefix)+key, value)
	}
}

func stripCustomQueryValues(query url.Values) (headers, extraData map[string]string) {
	headers = make(map[string]string)
	extraData = make(map[string]string)

	for key, values := range query {
		if key[0] == headerPrefix {
			headerKey := normalizedHeaderKey(key[1:])
			headers[headerKey] = values[0]
		} else if key[0] == extraPrefix {
			extraData[key[1:]] = values[0]
		} else {
			continue
		}
		delete(query, key)
	}
	return headers, extraData
}
