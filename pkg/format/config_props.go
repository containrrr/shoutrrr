package format

import (
	"errors"
	"reflect"

	"github.com/nicholas-fedor/shoutrrr/pkg/types"
)

// GetConfigPropFromString deserializes a config property from a string representation using the ConfigProp interface.
func GetConfigPropFromString(structType reflect.Type, value string) (reflect.Value, error) {
	valuePtr := reflect.New(structType)

	configProp, ok := valuePtr.Interface().(types.ConfigProp)
	if !ok {
		return reflect.Value{}, errors.New("struct field cannot be used as a prop")
	}

	if err := configProp.SetFromProp(value); err != nil {
		return reflect.Value{}, err
	}

	return valuePtr, nil
}

// GetConfigPropString serializes a config property to a string representation using the ConfigProp interface.
func GetConfigPropString(propPtr reflect.Value) (string, error) {
	if propPtr.Kind() != reflect.Ptr {
		propVal := propPtr
		propPtr = reflect.New(propVal.Type())
		propPtr.Elem().Set(propVal)
	}

	if propPtr.CanInterface() {
		if configProp, ok := propPtr.Interface().(types.ConfigProp); ok {
			return configProp.GetPropValue()
		}
	}

	return "", errors.New("struct field cannot be used as a prop")
}
