package pagerduty

import (
	"fmt"
	"reflect"
	"strconv"
)

const (
	tagUrl     = "url"
	tagDefault = "default"
)

func setUrlDefaults(config *Config) {
	cfg := reflect.TypeOf(*config)
	values := getDefaultUrlValues(cfg)

	for fieldName, defVal := range values {
		field := reflect.ValueOf(config).Elem().FieldByName(fieldName)

		switch field.Kind() {
		case reflect.String:
			field.SetString(defVal)
		case reflect.Int:
			intVal, err := strconv.Atoi(defVal)
			if err != nil {
				fmt.Errorf("enable to convert %q to an int: %w", defVal, err)
			}
			field.SetInt(int64(intVal))
		case reflect.Uint16:
			intVal, err := strconv.Atoi(defVal)
			if err != nil {
				fmt.Errorf("enable to convert %q to an int: %w", defVal, err)
			}
			field.SetUint(uint64(intVal))
		}
	}
}

// getDefaultUrlValues finds field names tagged with `url` and returns a map of the field name to their default value.
func getDefaultUrlValues(cfg reflect.Type) map[string]string {
	fieldToValue := make(map[string]string)

	for i := 0; i < cfg.NumField(); i++ {
		field := cfg.Field(i)

		_, ok := field.Tag.Lookup(tagUrl)
		if ok {
			defaultVal := field.Tag.Get(tagDefault)
			if defaultVal != "" {
				fieldToValue[field.Name] = defaultVal
			}
		}
	}
	return fieldToValue
}
