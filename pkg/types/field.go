package types

import "sort"

// Field is a Key/Value pair used for extra data in log messages
type Field struct {
	Key   string
	Value string
}

// FieldsFromMap creates a Fields slice from a map, optionally sorting keys
func FieldsFromMap(fieldMap map[string]string, sorted bool) []Field {
	keys := make([]string, 0, len(fieldMap))
	fields := make([]Field, 0, len(fieldMap))

	for key := range fieldMap {
		keys = append(keys, key)
	}

	if sorted {
		sort.Strings(keys)
	}

	for i, key := range keys {
		fields[i].Key = key
		fields[i].Value = fieldMap[key]
	}

	return fields
}
