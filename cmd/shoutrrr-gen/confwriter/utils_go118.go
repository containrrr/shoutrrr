//go:build go1.18
// +build go1.18

package confwriter

import (
	"sort"
)

func mapKeys[V any](m map[string]V) ([]string, int) {
	maxLen := 0
	keys := make([]string, 0, len(m))
	for name := range m {
		if len(name) > maxLen {
			maxLen = len(name)
		}
		keys = append(keys, name)
	}
	sort.Strings(keys)
	return keys, maxLen
}
