//go:build !go1.18
// +build !go1.18

package confwriter

import (
	"reflect"
	"sort"
)

func mapKeys(m interface{}) ([]string, int) {
	maxLen := 0

	v := reflect.ValueOf(m)

	keys := make([]string, v.Len())

	for i, key := range v.MapKeys() {
		name := key.String()
		if len(name) > maxLen {
			maxLen = len(name)
		}
		keys[i] = name
	}

	sort.Strings(keys)

	return keys, maxLen
}
