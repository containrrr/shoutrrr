package standard

import (
	"errors"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
	"log"
	"reflect"
	"sort"
	"strings"
)

// KeyPropConfig implements the ServiceConfig interface for services that uses key tags for query props
type KeyPropConfig struct {
	confValue reflect.Value
	keyFields map[string]format.FieldInfo
	keys      []string
}

// BindKeys is called to map config fields to it's tagged query keys
func (c *KeyPropConfig) BindKeys(config types.ServiceConfig) {
	_, fields := format.GetConfigFormat(config)
	keyFields := make(map[string]format.FieldInfo, len(fields))
	keys := make([]string, 0, len(fields))
	for _, field := range fields {
		key := strings.ToLower(field.Key)
		if key != "" {
			keys = append(keys, key)
			keyFields[key] = field
		}
	}
	c.keyFields = keyFields
	c.confValue = reflect.ValueOf(config)
	if c.confValue.Kind() == reflect.Ptr {
		c.confValue = c.confValue.Elem()
	}
	sort.Strings(keys)
	c.keys = keys
}

// QueryFields returns a list of tagged keys
func (c *KeyPropConfig) QueryFields() []string {
	if c.keys == nil {
		log.Panic("KeyPropConfig.QueryFields called before BindKeys")
	}

	return c.keys
}

// Get returns the value of a config property tagged with the corresponding key
func (c *KeyPropConfig) Get(key string) (string, error) {
	if c.keyFields == nil {
		return "", errors.New("KeyPropConfig.Get called before BindKeys")
	}

	if field, found := c.keyFields[strings.ToLower(key)]; found {
		return format.GetConfigFieldString(c.confValue, field)
	}

	return "", fmt.Errorf("%v is not a valid config key", key)
}

// Set sets the value of a config property tagged with the corresponding key
func (c *KeyPropConfig) Set(key string, value string) error {
	if c.keyFields == nil {
		return errors.New("KeyPropConfig.Set called before BindKeys")
	}

	if field, found := c.keyFields[strings.ToLower(key)]; found {
		valid, err := format.SetConfigField(c.confValue, field, value)
		if !valid && err == nil {
			return errors.New("invalid value for type")
		}
		return err
	}

	return fmt.Errorf("%v is not a valid config key", key)
}

// UpdateConfigFromParams mutates the provided config, updating the values from it's corresponding params
func (c *KeyPropConfig) UpdateConfigFromParams(config types.ServiceConfig, params *types.Params) error {
	for key, val := range *params {
		if err := config.Set(key, val); err != nil {
			return err
		}
	}
	return nil
}
