package format

import (
	"errors"
	"github.com/containrrr/shoutrrr/pkg/types"
	r "reflect"
)

// GetConfigPropFromString deserializes a config property from a string representation using the ConfigProp interface
func GetConfigPropFromString(structType r.Type, value string) (r.Value, error) {
	valuePtr := r.New(structType)
	configProp, ok := valuePtr.Interface().(types.ConfigProp)
	if !ok {
		return r.Value{}, errors.New("struct field cannot be used as a prop")
	}

	if err := configProp.SetFromProp(value); err != nil {
		return r.Value{}, err
	}

	return valuePtr, nil
}

// GetConfigPropString serializes a config property to a string representation using the ConfigProp interface
func GetConfigPropString(propPtr r.Value) (string, error) {

	if propPtr.Kind() != r.Ptr {
		propVal := propPtr
		propPtr = r.New(propVal.Type())
		propPtr.Elem().Set(propVal)
	}

	if propPtr.CanInterface() {
		if configProp, ok := propPtr.Interface().(types.ConfigProp); ok {
	return configProp.GetPropValue()
}
	}

	return "", errors.New("struct field cannot be used as a prop")
}
