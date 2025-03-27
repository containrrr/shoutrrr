package format

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/nicholas-fedor/shoutrrr/pkg/types"
	"github.com/nicholas-fedor/shoutrrr/pkg/util"
)

// Constants for map parsing.
const (
	KeyValuePairSize = 2 // Number of elements in a key:value pair
)

// GetServiceConfig returns the inner config of a service.
func GetServiceConfig(service types.Service) types.ServiceConfig {
	serviceValue := reflect.Indirect(reflect.ValueOf(service))

	configField, ok := serviceValue.Type().FieldByName("Config")
	if !ok {
		panic("service does not have a Config field") // Or handle gracefully
	}

	configRef := serviceValue.FieldByIndex(configField.Index)

	if configRef.IsNil() {
		configType := configField.Type
		if configType.Kind() == reflect.Ptr {
			configType = configType.Elem()
		}

		return reflect.New(configType).Interface().(types.ServiceConfig)
	}

	return configRef.Interface().(types.ServiceConfig)
}

// ColorFormatTree returns a color highlighted string representation of a node tree.
func ColorFormatTree(rootNode *ContainerNode, withValues bool) string {
	return ConsoleTreeRenderer{WithValues: withValues}.RenderTree(rootNode, "")
}

// GetServiceConfigFormat returns type and field information about a ServiceConfig, resolved from it's Service.
func GetServiceConfigFormat(service types.Service) *ContainerNode {
	serviceConfig := GetServiceConfig(service)

	return GetConfigFormat(serviceConfig)
}

// GetConfigFormat returns type and field information about a ServiceConfig.
func GetConfigFormat(config types.ServiceConfig) *ContainerNode {
	return getRootNode(config)
}

// SetConfigField deserializes the inputValue and sets the field of a config to that value.
func SetConfigField(config reflect.Value, field FieldInfo, inputValue string) (valid bool, err error) {
	configField := config.FieldByName(field.Name)
	fieldKind := field.Type.Kind()

	if fieldKind == reflect.String {
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
	} else if fieldKind >= reflect.Uint && fieldKind <= reflect.Uint64 {
		var value uint64

		number, base := util.StripNumberPrefix(inputValue)

		value, err = strconv.ParseUint(number, base, field.Type.Bits())
		if err == nil {
			configField.SetUint(value)

			return true, nil
		}
	} else if fieldKind >= reflect.Int && fieldKind <= reflect.Int64 {
		var value int64

		number, base := util.StripNumberPrefix(inputValue)

		value, err = strconv.ParseInt(number, base, field.Type.Bits())
		if err == nil {
			configField.SetInt(value)

			return true, nil
		}
	} else if fieldKind == reflect.Bool {
		value, ok := ParseBool(inputValue, false)
		if !ok {
			return false, errors.New("accepted values are 1, true, yes or 0, false, no")
		}

		configField.SetBool(value)

		return true, nil
	} else if fieldKind == reflect.Map {
		keyKind := field.Type.Key().Kind()
		valueType := field.Type.Elem()

		if keyKind != reflect.String {
			return false, errors.New("field key format is not supported")
		}

		mapValue := reflect.MakeMap(field.Type)
		pairs := strings.Split(inputValue, ",")

		for _, pair := range pairs {
			elems := strings.Split(pair, ":")

			if len(elems) != KeyValuePairSize {
				return false, errors.New("invalid field value format")
			}

			key := elems[0]
			valueRaw := elems[1]

			value, err := getMapValue(valueType, valueRaw)
			if err != nil {
				return false, err
			}

			mapValue.SetMapIndex(reflect.ValueOf(key), value)
		}

		configField.Set(mapValue)

		return true, nil
	} else if fieldKind == reflect.Struct {
		valuePtr, err := GetConfigPropFromString(field.Type, inputValue)
		if err != nil {
			return false, err
		}

		configField.Set(valuePtr.Elem())

		return true, nil
	} else if fieldKind >= reflect.Slice || fieldKind == reflect.Array {
		elemType := field.Type.Elem()
		elemValType := elemType
		elemKind := elemType.Kind()

		if elemKind == reflect.Ptr {
			// When updating a pointer slice, use the value type kind
			elemValType = elemType.Elem()
			elemKind = elemValType.Kind()
		}

		if elemKind != reflect.Struct && elemKind != reflect.String {
			return false, errors.New("field format is not supported")
		}

		values := strings.Split(inputValue, string(field.ItemSeparator))

		var value reflect.Value

		if elemKind == reflect.Struct {
			propValues := reflect.MakeSlice(reflect.SliceOf(elemType), 0, len(values))

			for _, v := range values {
				propPtr, err := GetConfigPropFromString(elemValType, v)
				if err != nil {
					return false, err
				}

				propVal := propPtr

				// If not a pointer slice, dereference the value
				if elemType.Kind() != reflect.Ptr {
					propVal = propPtr.Elem()
				}

				propValues = reflect.Append(propValues, propVal)
			}

			value = propValues
		} else {
			// Use the split string parts as the target value
			value = reflect.ValueOf(values)
		}

		if fieldKind == reflect.Array {
			arrayLen := field.Type.Len()
			if len(values) != arrayLen {
				return false, fmt.Errorf("field value count needs to be %d", arrayLen)
			}

			arr := reflect.Indirect(reflect.New(field.Type))
			reflect.Copy(arr, value)
			value = arr
		}

		configField.Set(value)

		return true, nil
	} else {
		err = fmt.Errorf("invalid field kind %v", fieldKind)
	}

	return false, err
}

func getMapValue(valueType reflect.Type, valueRaw string) (reflect.Value, error) {
	kind := valueType.Kind()
	switch kind {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return getMapUintValue(valueRaw, valueType.Bits(), kind)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return getMapIntValue(valueRaw, valueType.Bits(), kind)
	case reflect.String:
		return reflect.ValueOf(valueRaw), nil
	default:
	}

	return reflect.Value{}, errors.New("map value format is not supported")
}

func getMapUintValue(valueRaw string, bits int, kind reflect.Kind) (reflect.Value, error) {
	number, base := util.StripNumberPrefix(valueRaw)
	numValue, err := strconv.ParseUint(number, base, bits)

	switch kind {
	case reflect.Uint:
		return reflect.ValueOf(uint(numValue)), err
	case reflect.Uint8:
		return reflect.ValueOf(uint8(numValue)), err
	case reflect.Uint16:
		return reflect.ValueOf(uint16(numValue)), err
	case reflect.Uint32:
		return reflect.ValueOf(uint32(numValue)), err
	case reflect.Uint64:
	default:
	}

	return reflect.ValueOf(numValue), err
}

func getMapIntValue(valueRaw string, bits int, kind reflect.Kind) (reflect.Value, error) {
	number, base := util.StripNumberPrefix(valueRaw)
	numValue, err := strconv.ParseInt(number, base, bits)

	switch kind {
	case reflect.Int:
		return reflect.ValueOf(int(numValue)), err
	case reflect.Int8:
		return reflect.ValueOf(int8(numValue)), err
	case reflect.Int16:
		return reflect.ValueOf(int16(numValue)), err
	case reflect.Int32:
		return reflect.ValueOf(int32(numValue)), err
	case reflect.Int64:
	default:
	}

	return reflect.ValueOf(numValue), err
}

// GetConfigFieldString serializes the config field value to a string representation.
func GetConfigFieldString(config reflect.Value, field FieldInfo) (value string, err error) {
	configField := config.FieldByName(field.Name)

	strVal, token := getValueNodeValue(configField, &field)
	if token == ErrorToken {
		err = errors.New("invalid field value")
	}

	return strVal, err
}
