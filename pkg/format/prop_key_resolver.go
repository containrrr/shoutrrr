package format

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	t "github.com/containrrr/shoutrrr/pkg/types"
)

// PropKeyResolver implements the ConfigQueryResolver interface for services that uses key tags for query props
type PropKeyResolver struct {
	confValue reflect.Value
	keyFields map[string]FieldInfo
	keys      []string
}

// NewPropKeyResolver creates a new PropKeyResolver and initializes it using the provided config
func NewPropKeyResolver(config t.ServiceConfig) PropKeyResolver {

	configNode := GetConfigFormat(config)
	items := configNode.Items

	keyFields := make(map[string]FieldInfo, len(items))
	keys := make([]string, 0, len(items))
	for _, item := range items {
		field := *item.Field()
		for _, key := range field.Keys {
			key = strings.ToLower(key)
			if key != "" {
				keys = append(keys, key)
				keyFields[key] = field
			}
		}
	}

	sort.Strings(keys)

	confValue := reflect.ValueOf(config)
	if confValue.Kind() == reflect.Ptr {
		confValue = confValue.Elem()
	}

	return PropKeyResolver{
		keyFields: keyFields,
		confValue: confValue,
		keys:      keys,
	}
}

// QueryFields returns a list of tagged keys
func (pkr *PropKeyResolver) QueryFields() []string {
	return pkr.keys
}

// Get returns the value of a config property tagged with the corresponding key
func (pkr *PropKeyResolver) Get(key string) (string, error) {
	if field, found := pkr.keyFields[strings.ToLower(key)]; found {
		return GetConfigFieldString(pkr.confValue, field)
	}

	return "", fmt.Errorf("%v is not a valid config key", key)
}

// Set sets the value of it's bound struct's property, tagged with the corresponding key
func (pkr *PropKeyResolver) Set(key string, value string) error {
	return pkr.set(pkr.confValue, key, value)
}

// set sets the value of a target struct tagged with the corresponding key
func (pkr *PropKeyResolver) set(target reflect.Value, key string, value string) error {
	if field, found := pkr.keyFields[strings.ToLower(key)]; found {
		valid, err := SetConfigField(target, field, value)
		if !valid && err == nil {
			return errors.New("invalid value for type")
		}
		return err
	}

	return fmt.Errorf("%v is not a valid config key %v", key, pkr.keys)
}

// UpdateConfigFromParams mutates the provided config, updating the values from it's corresponding params
// If the provided config is nil, the internal config will be updated instead.
// The error returned is the first error that occurred, subsequent errors are just discarded.
func (pkr *PropKeyResolver) UpdateConfigFromParams(config t.ServiceConfig, params *t.Params) (firstError error) {
	confValue := pkr.configValueOrInternal(config)
	if params != nil {
		for key, val := range *params {
			if err := pkr.set(confValue, key, val); err != nil && firstError == nil {
				firstError = err
			}
		}
	}
	return
}

// SetDefaultProps mutates the provided config, setting the tagged fields with their default values
// If the provided config is nil, the internal config will be updated instead.
// The error returned is the first error that occurred, subsequent errors are just discarded.
func (pkr *PropKeyResolver) SetDefaultProps(config t.ServiceConfig) (firstError error) {
	confValue := pkr.configValueOrInternal(config)
	for key, info := range pkr.keyFields {
		if err := pkr.set(confValue, key, info.DefaultValue); err != nil && firstError == nil {
			firstError = err
		}
	}
	return
}

// Bind is called to create a new instance of the PropKeyResolver, with he internal config reference
// set to the provided config. This should only be used for configs of the same type.
func (pkr *PropKeyResolver) Bind(config t.ServiceConfig) PropKeyResolver {
	bound := *pkr
	bound.confValue = configValue(config)
	return bound
}

// GetConfigQueryResolver returns the config itself if it implements ConfigQueryResolver
// otherwise it creates and returns a PropKeyResolver that implements it
func GetConfigQueryResolver(config t.ServiceConfig) t.ConfigQueryResolver {
	var resolver t.ConfigQueryResolver
	var ok bool
	if resolver, ok = config.(t.ConfigQueryResolver); !ok {
		pkr := NewPropKeyResolver(config)
		resolver = &pkr
	}
	return resolver
}

// KeyIsPrimary returns whether the key is the primary (and not an alias)
func (pkr *PropKeyResolver) KeyIsPrimary(key string) bool {
	return pkr.keyFields[key].Keys[0] == key
}

func (pkr *PropKeyResolver) configValueOrInternal(config t.ServiceConfig) reflect.Value {
	if config != nil {
		return configValue(config)
	}
	return pkr.confValue
}

func configValue(config t.ServiceConfig) reflect.Value {
	return reflect.Indirect(reflect.ValueOf(config))
}

// IsDefault returns whether the specified key value is the default value
func (pkr *PropKeyResolver) IsDefault(key string, value string) bool {
	return pkr.keyFields[key].DefaultValue == value
}
