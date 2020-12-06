package format

import (
	"errors"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/types"
	"reflect"
	"sort"
	"strings"
)

// KeyPropConfig implements the ServiceConfig interface for services that uses key tags for query props
type PropKeyResolver struct {
	confValue reflect.Value
	keyFields map[string]FieldInfo
	keys      []string
}

// BindKeys is called to map config fields to it's tagged query keys
func NewPropKeyResolver(config types.ServiceConfig) PropKeyResolver {

	_, fields := GetConfigFormat(config)
	keyFields := make(map[string]FieldInfo, len(fields))
	keys := make([]string, 0, len(fields))
	for _, field := range fields {
		key := strings.ToLower(field.Key)
		if key != "" {
			keys = append(keys, key)
			keyFields[key] = field
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
func (c *PropKeyResolver) set(target reflect.Value, key string, value string) error {
	if field, found := c.keyFields[strings.ToLower(key)]; found {
		valid, err := SetConfigField(target, field, value)
		if !valid && err == nil {
			return errors.New("invalid value for type")
		}
		return err
	}

	return fmt.Errorf("%v is not a valid config key %v", key, c.keys)
}

// UpdateConfigFromParams mutates the provided config, updating the values from it's corresponding params
func (pkr *PropKeyResolver) UpdateConfigFromParams(config types.ServiceConfig, params *types.Params) error {
	if params != nil {
		for key, val := range *params {
			if err := pkr.set(reflect.ValueOf(config), key, val); err != nil {
				return err
			}
		}
	}
	return nil
}

func (pkr *PropKeyResolver) Bind(config types.ServiceConfig) PropKeyResolver {
	bound := *pkr
	bound.confValue = reflect.ValueOf(config)
	return bound
}

func GetConfigQueryResolver(config types.ServiceConfig) types.ConfigQueryResolver {
	var resolver types.ConfigQueryResolver
	var ok bool
	if resolver, ok = config.(types.ConfigQueryResolver); !ok {
		pkr := NewPropKeyResolver(config)
		resolver = &pkr
	}
	return resolver
}
