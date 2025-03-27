package dedupe

import "slices"

// RemoveDuplicates from a slice of strings.
func RemoveDuplicates(src []string) []string {
	unique := make([]string, 0, len(src))
	for _, s := range src {
		found := slices.Contains(unique, s)
		if !found {
			unique = append(unique, s)
		}
	}

	return unique
}
