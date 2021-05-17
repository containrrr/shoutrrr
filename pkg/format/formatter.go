package format

import (
	"errors"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util"
	r "reflect"
	"strconv"
	"strings"
	"unsafe"
)

// GetServiceConfig returns the inner config of a service
func GetServiceConfig(service types.Service) types.ServiceConfig {
	serviceValue := r.Indirect(r.ValueOf(service))
	configField, _ := serviceValue.Type().FieldByName("config")
	configRef := serviceValue.FieldByIndex(configField.Index)

	var ourRef r.Value
	if configRef.IsNil() {
		configType := configField.Type
		if configType.Kind() == r.Ptr {
			configType = configType.Elem()
		}
		ourRef = r.New(configType)
	} else {
		ourRef = r.NewAt(configRef.Type(), unsafe.Pointer(configRef.UnsafeAddr())).Elem()
	}

	//
	return ourRef.Interface().(types.ServiceConfig)
}

// ColorFormatTree returns a color highlighted string representation of a node tree
func ColorFormatTree(rootNode *ContainerNode, withValues bool) string {
	return ConsoleTreeRenderer{WithValues: withValues}.RenderTree(rootNode, "")
}

// GetServiceConfigFormat returns type and field information about a ServiceConfig, resolved from it's Service
func GetServiceConfigFormat(service types.Service) *ContainerNode {
	serviceConfig := GetServiceConfig(service)
	return GetConfigFormat(serviceConfig)
}

// GetConfigFormat returns type and field information about a ServiceConfig
func GetConfigFormat(config types.ServiceConfig) *ContainerNode {
	return getRootNode(config)
}

// SetConfigField deserializes the inputValue and sets the field of a config to that value
func SetConfigField(config r.Value, field FieldInfo, inputValue string) (valid bool, err error) {
	configField := config.FieldByName(field.Name)
	fieldKind := field.Type.Kind()

	if fieldKind == r.String {
		configField.SetString(inputValue)
		return true, nil
	} else if field.EnumFormatter != nil {
		value := field.EnumFormatter.Parse(inputValue)
		if value == EnumInvalid {
			enumNames := strings.Join(field.EnumFormatter.Names(), ", ")
			return false, fmt.Errorf("not a one of %v", enumNames)
		}

		configField.SetInt(int64(value))
		return true, nil

	} else if fieldKind >= r.Uint && fieldKind <= r.Uint64 {
		var value uint64
		number, base := util.StripNumberPrefix(inputValue)
		value, err = strconv.ParseUint(number, base, field.Type.Bits())
		if err == nil {
			configField.SetUint(value)
			return true, nil
		}
	} else if fieldKind >= r.Int && fieldKind <= r.Int64 {
		var value int64
		number, base := util.StripNumberPrefix(inputValue)
		value, err = strconv.ParseInt(number, base, field.Type.Bits())
		if err == nil {
			configField.SetInt(value)
			return true, nil
		}
	} else if fieldKind == r.Bool {
		value, ok := ParseBool(inputValue, false)
		if !ok {
			return false, errors.New("accepted values are 1, true, yes or 0, false, no")
		}

		configField.SetBool(value)
		return true, nil
	} else if fieldKind == r.Map {
		keyKind := field.Type.Key().Kind()
		valueType := field.Type.Elem()
		if keyKind != r.String {
			return false, errors.New("field key format is not supported")
		}

		mapValue := r.MakeMap(field.Type)

		pairs := strings.Split(inputValue, ",")
		for _, pair := range pairs {
			elems := strings.Split(pair, ":")
			if len(elems) != 2 {
				return false, errors.New("invalid field value format")
			}
			key := elems[0]
			valueRaw := elems[1]
			value, err := getMapValue(valueType, valueRaw)
			if err != nil {
				return false, err
			}

			mapValue.SetMapIndex(r.ValueOf(key), value)
		}

		configField.Set(mapValue)
		return true, nil
	} else if fieldKind == r.Struct {
		valuePtr, err := GetConfigPropFromString(field.Type, inputValue)
		if err != nil {
			return false, err
		}
		configField.Set(valuePtr.Elem())
		return true, nil
	} else if fieldKind >= r.Slice || fieldKind == r.Array {
		elemType := field.Type.Elem()
		elemValType := elemType
		elemKind := elemType.Kind()

		if elemKind == r.Ptr {
			// When updating a pointer slice, use the value type kind
			elemValType = elemType.Elem()
			elemKind = elemValType.Kind()
		}

		if elemKind != r.Struct && elemKind != r.String {
			return false, errors.New("field format is not supported")
		}

		values := strings.Split(inputValue, ",")

		var value r.Value
		if elemKind == r.Struct {
			propValues := r.MakeSlice(r.SliceOf(elemType), 0, len(values))
			for _, v := range values {
				propPtr, err := GetConfigPropFromString(elemValType, v)
				if err != nil {
					return false, err
				}
				propVal := propPtr

				// If not a pointer slice, dereference the value
				if elemType.Kind() != r.Ptr {
					propVal = propPtr.Elem()
				}
				propValues = r.Append(propValues, propVal)
			}
			value = propValues
		} else {
			// Use the split string parts as the target value
			value = r.ValueOf(values)
		}

		if fieldKind == r.Array {
			arrayLen := field.Type.Len()
			if len(values) != arrayLen {
				return false, fmt.Errorf("field value count needs to be %d", arrayLen)
			}
			arr := r.Indirect(r.New(field.Type))
			r.Copy(arr, value)
			value = arr
		}

		configField.Set(value)
		return true, nil

	} else {
		err = fmt.Errorf("invalid field kind %v", fieldKind)
	}

	return false, err

}

func getMapValue(valueType r.Type, valueRaw string) (r.Value, error) {
	kind := valueType.Kind()
	switch kind {
	case r.Uint, r.Uint8, r.Uint16, r.Uint32, r.Uint64:
		return getMapUintValue(valueRaw, valueType.Bits(), kind)
	case r.Int, r.Int8, r.Int16, r.Int32, r.Int64:
		return getMapIntValue(valueRaw, valueType.Bits(), kind)
	case r.String:
		return r.ValueOf(valueRaw), nil
	default:
	}
	return r.Value{}, errors.New("map value format is not supported")
}

func getMapUintValue(valueRaw string, bits int, kind r.Kind) (r.Value, error) {
	number, base := util.StripNumberPrefix(valueRaw)
	numValue, err := strconv.ParseUint(number, base, bits)

	switch kind {
	case r.Uint:
		return r.ValueOf(uint(numValue)), err
	case r.Uint8:
		return r.ValueOf(uint8(numValue)), err
	case r.Uint16:
		return r.ValueOf(uint16(numValue)), err
	case r.Uint32:
		return r.ValueOf(uint32(numValue)), err
	case r.Uint64:
	default:
	}
	return r.ValueOf(numValue), err
}

func getMapIntValue(valueRaw string, bits int, kind r.Kind) (r.Value, error) {
	number, base := util.StripNumberPrefix(valueRaw)
	numValue, err := strconv.ParseInt(number, base, bits)

	switch kind {
	case r.Int:
		return r.ValueOf(int(numValue)), err
	case r.Int8:
		return r.ValueOf(int8(numValue)), err
	case r.Int16:
		return r.ValueOf(int16(numValue)), err
	case r.Int32:
		return r.ValueOf(int32(numValue)), err
	case r.Int64:
	default:
	}
	return r.ValueOf(numValue), err
}

// GetConfigFieldString serializes the config field value to a string representation
func GetConfigFieldString(config r.Value, field FieldInfo) (value string, err error) {
	configField := config.FieldByName(field.Name)

	strVal, token := getValueNodeValue(configField, &field)
	if token == ErrorToken {
		err = errors.New("invalid field value")
	}
	return strVal, err
}
