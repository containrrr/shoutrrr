package types

import "sort"

type Field struct {
	Key string
	Value string
}

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