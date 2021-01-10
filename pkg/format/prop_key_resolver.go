package format

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/types"
)

// PropKeyResolver implements the ConfigQueryResolver interface for services that uses key tags for query props
type PropKeyResolver struct {
	confValue reflect.Value
	keyFields map[string]FieldInfo
	keys      []string
}

// NewPropKeyResolver creates a new PropKeyResolver and initializes it using the provided config
func NewPropKeyResolver(config types.ServiceConfig) PropKeyResolver {

	_, fields := GetConfigFormat(config)
	keyFields := make(map[string]FieldInfo, len(fields))
	keys := make([]string, 0, len(fields))
	for _, field := range fields {
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
func (pkr *PropKeyResolver) UpdateConfigFromParams(config types.ServiceConfig, params *types.Params) error {
	if params != nil {
		for key, val := range *params {
			if err := pkr.set(reflect.Indirect(reflect.ValueOf(config)), key, val); err != nil {
				return err
			}
		}
	}
	return nil
}

// SetDefaultProps mutates the provided config, setting the tagged fields with their default values
func (pkr *PropKeyResolver) SetDefaultProps(config types.ServiceConfig) error {
	for key, info := range pkr.keyFields {
		if err := pkr.set(reflect.Indirect(reflect.ValueOf(config)), key, info.DefaultValue); err != nil {
			return err
		}
	}
	return nil
}

// Bind is called to set the internal config reference for the PropKeyResolver
func (pkr *PropKeyResolver) Bind(config types.ServiceConfig) PropKeyResolver {
	bound := *pkr
	bound.confValue = reflect.ValueOf(config)
	return bound
}

// GetConfigQueryResolver returns the config itself if it implements ConfigQueryResolver
// otherwise it creates and returns a PropKeyResolver that implements it
func GetConfigQueryResolver(config types.ServiceConfig) types.ConfigQueryResolver {
	var resolver types.ConfigQueryResolver
	var ok bool
	if resolver, ok = config.(types.ConfigQueryResolver); !ok {
		pkr := NewPropKeyResolver(config)
		resolver = &pkr
	}
	return resolver
}

// KeyIsPrimary returns whether the key is the primary (and not an alias)
func (pkr *PropKeyResolver) KeyIsPrimary(key string) bool {
	return pkr.keyFields[key].Keys[0] == key
}
