package dedupe

// RemoveDuplicates from a slice of strings
func RemoveDuplicates(src []string) []string {
	unique := make([]string, 0, len(src))

	for _, s := range src {
		found := false
		for _, us := range unique {
			if us == s {
				found = true
				break
			}
		}
		if !found {
			unique = append(unique, s)
		}
	}

	return unique
}
