package format

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/fatih/color"

	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util"
)

// GetConfigMap returns a string map of a given Config struct
func GetConfigMap(service types.Service) (map[string]string, int) {

	configRef := reflect.ValueOf(service).Elem().FieldByName("config")
	configType := configRef.Type().Elem()
	cr := reflect.NewAt(configRef.Type(), unsafe.Pointer(configRef.UnsafeAddr())).Elem()
	config := cr.Interface().(types.ServiceConfig)

	formatter := formatter{
		EnumFormatters: config.Enums(),
		MaxDepth:       10,
	}
	return formatter.formatStructMap(configType, config, 0)
}

// GetServiceConfigFormat returns type and field information about a ServiceConfig, resolved from it's Service
func GetServiceConfigFormat(service types.Service) (reflect.Type, []FieldInfo) {
	configRef := reflect.ValueOf(service).Elem().FieldByName("config")
	configType := configRef.Type().Elem()

	config := reflect.New(configType)
	serviceConfig := config.Interface().(types.ServiceConfig)

	return GetConfigFormat(serviceConfig)
}

// GetConfigFormat returns type and field information about a ServiceConfig
func GetConfigFormat(serviceConfig types.ServiceConfig) (reflect.Type, []FieldInfo) {
	configType := reflect.TypeOf(serviceConfig)
	if configType.Kind() == reflect.Ptr {
		configType = configType.Elem()
	}

	formatter := formatter{
		EnumFormatters: serviceConfig.Enums(),
		MaxDepth:       10,
	}
	return configType, formatter.getStructFieldInfo(configType)
}

type formatter struct {
	EnumFormatters map[string]types.EnumFormatter
	MaxDepth       uint8
	Errors         []error
}

func (fmtr *formatter) formatStructMap(structType reflect.Type, structItem interface{}, depth uint8) (map[string]string, int) {
	values := reflect.ValueOf(structItem)

	if values.Kind() == reflect.Ptr {
		values = values.Elem()
	}

	infoFields := fmtr.getStructFieldInfo(structType)

	numFields := len(infoFields)
	valueMap := make(map[string]string, numFields)
	nextDepth := depth + 1
	maxKeyLen := 0

	for i := 0; i < numFields; i++ {
		field := infoFields[i]

		var value string
		valueLen := 0
		preLen := 16

		isEnum := field.EnumFormatter != nil

		if values.IsValid() {
			// Add some space to print the value
			preLen = 40
			if isEnum {
				fieldVal := values.Field(i)
				kind := fieldVal.Kind()
				if kind == reflect.Int {
					valueStr := field.EnumFormatter.Print(int(fieldVal.Int()))
					value = ColorizeEnum(valueStr)
					valueLen = len(valueStr)
				} else {
					err := fmt.Errorf("incorrect enum type '%s' for field '%s'", kind, field.Name)
					fmtr.Errors = append(fmtr.Errors, err)
				}
			} else if nextDepth < fmtr.MaxDepth {
				value, valueLen = fmtr.getFieldValueString(values.FieldByName(field.Name), field.Base, nextDepth)
			}
		} else {
			// Since no values was supplied, let's substitute the value with the type
			typeName := field.Type.String()
			valueLen = len(typeName)
			value = color.CyanString(typeName)
		}

		if len(field.Description) > 0 {

			prePad := strings.Repeat(" ", util.Max(preLen-valueLen, 1))
			postPad := strings.Repeat(" ", util.Max(60-len(field.Description), 1))

			value += " " + prePad + ColorizeDesc(field.Description) + postPad
		}

		if len(field.Template) > 0 {
			value += fmt.Sprintf(" <Template: %s>", ColorizeString(field.Template))
		}

		if len(field.DefaultValue) > 0 {
			value += fmt.Sprintf(" <Default: %s>", ColorizeValue(field.DefaultValue, isEnum))
		}

		if isEnum {
			value += fmt.Sprintf(" [%s]", strings.Join(field.EnumFormatter.Names(), ", "))
		}

		valueMap[field.Name] = value
		keyLen := len(field.Name)
		if keyLen > maxKeyLen {
			maxKeyLen = keyLen
		}
	}

	return valueMap, maxKeyLen
}

// FieldInfo is the meta data about a config field
type FieldInfo struct {
	Name          string
	Type          reflect.Type
	EnumFormatter types.EnumFormatter
	Description   string
	DefaultValue  string
	Template      string
	Required      bool
	Title         bool
	Base          int
	Keys          []string
}

func (fmtr *formatter) getStructFieldInfo(structType reflect.Type) []FieldInfo {

	numFields := structType.NumField()
	fields := make([]FieldInfo, 0, numFields)
	maxKeyLen := 0

	for i := 0; i < numFields; i++ {
		fieldDef := structType.Field(i)

		if fieldDef.Anonymous {
			// This is an embedded field, which should not be part of the Config output
			continue
		}

		info := FieldInfo{
			Name:     fieldDef.Name,
			Type:     fieldDef.Type,
			Required: true,
			Title:    false,
		}

		if util.IsNumeric(fieldDef.Type.Kind()) {
			info.Base = getFieldBase(fieldDef)
		}

		if tag, ok := fieldDef.Tag.Lookup("desc"); ok {
			info.Description = tag
		}

		if tag, ok := fieldDef.Tag.Lookup("tpl"); ok {
			info.Template = tag
		}

		if tag, ok := fieldDef.Tag.Lookup("default"); ok {
			info.Required = false
			info.DefaultValue = tag
		}

		if _, ok := fieldDef.Tag.Lookup("optional"); ok {
			info.Required = false
		}

		if _, ok := fieldDef.Tag.Lookup("title"); ok {
			info.Title = true
		}

		if tag, ok := fieldDef.Tag.Lookup("key"); ok {
			tag := strings.ToLower(tag)
			info.Keys = strings.Split(tag, ",")
		}

		if ef, isEnum := fmtr.EnumFormatters[fieldDef.Name]; isEnum {
			info.EnumFormatter = ef
		}

		fields = append(fields, info)
		keyLen := len(fieldDef.Name)
		if keyLen > maxKeyLen {
			maxKeyLen = keyLen
		}
	}

	return fields
}

func getFieldBase(field reflect.StructField) int {
	if tag, ok := field.Tag.Lookup("base"); ok {
		if base, err := strconv.ParseUint(tag, 10, 8); err == nil {
			return int(base)
		}
	}

	// Default to base 10 if not tagged
	return 10
}

func (fmtr *formatter) getFieldValueString(field reflect.Value, base int, depth uint8) (string, int) {

	nextDepth := depth + 1
	kind := field.Kind()

	if util.IsUnsignedInt(kind) {
		strVal := strconv.FormatUint(field.Uint(), base)
		if base == 16 {
			strVal = "0x" + strVal
		}
		return ColorizeNumber(strVal), len(strVal)
	}
	if util.IsSignedInt(kind) {
		strVal := strconv.FormatInt(field.Int(), base)
		return ColorizeNumber(strVal), len(strVal)
	}
	if kind == reflect.String {
		strVal := field.String()
		return ColorizeString(strVal), len(strVal)
	}
	if kind == reflect.Bool {
		val := field.Bool()
		if val {
			return ColorizeTrue(PrintBool(val)), 3
		}
		return ColorizeFalse(PrintBool(val)), 2

	}

	if util.IsCollection(kind) {
		fieldLen := field.Len()
		items := make([]string, fieldLen)
		totalLen := 4
		var itemLen int
		for i := 0; i < fieldLen; i++ {
			items[i], itemLen = fmtr.getFieldValueString(field.Index(i), base, nextDepth)
			totalLen += itemLen
		}
		if fieldLen > 1 {
			// Add space for separators
			totalLen += (fieldLen - 1) * 2
		}
		return fmt.Sprintf("[ %s ]", strings.Join(items, ", ")), totalLen
	}

	if kind == reflect.Map {
		items := make([]string, field.Len())
		iter := field.MapRange()
		index := 0
		// initial value for totalLen is surrounding curlies and spaces, and separating commas
		totalLen := 4 + (field.Len() - 1)
		for iter.Next() {
			key, keyLen := fmtr.getFieldValueString(iter.Key(), base, nextDepth)
			value, valueLen := fmtr.getFieldValueString(iter.Value(), base, nextDepth)
			items[index] = fmt.Sprintf("%s: %s", key, value)
			totalLen += keyLen + valueLen + 2
		}

		return fmt.Sprintf("{ %s }", strings.Join(items, ", ")), totalLen
	}
	if kind == reflect.Struct {
		structMap, _ := fmtr.formatStructMap(field.Type(), field, depth+1)
		structFieldCount := len(structMap)
		items := make([]string, structFieldCount)
		index := 0
		totalLen := 4 + (structFieldCount - 1)
		for key, value := range structMap {
			items[index] = fmt.Sprintf("%s: %s", key, value)
			index++
			totalLen += len(key) + 2 + len(value)
		}
		return fmt.Sprintf("< %s >", strings.Join(items, ", ")), totalLen
	}
	strVal := kind.String()
	return fmt.Sprintf("<?%s>", strVal), len(strVal) + 5
}

// SetConfigField deserializes the inputValue and sets the field of a config to that value
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
	} else if fieldKind >= reflect.Slice || fieldKind == reflect.Array {
		elemType := field.Type.Elem()
		elemKind := elemType.Kind()
		if elemKind != reflect.String {
			return false, errors.New("field format is not supported")
		}

		values := strings.Split(inputValue, ",")
		value := reflect.ValueOf(values)
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

// GetConfigFieldString serializes the config field value to a string representation
func GetConfigFieldString(config reflect.Value, field FieldInfo) (value string, err error) {
	configField := config.FieldByName(field.Name)
	fieldKind := field.Type.Kind()

	if fieldKind == reflect.String {
		return configField.String(), nil
	} else if field.EnumFormatter != nil {
		return field.EnumFormatter.Print(int(configField.Int())), nil
	} else if fieldKind >= reflect.Uint && fieldKind <= reflect.Uint64 {
		number := strconv.FormatUint(configField.Uint(), field.Base)
		if field.Base == 16 {
			number = "0x" + number
		}
		return number, nil
	} else if fieldKind >= reflect.Int && fieldKind <= reflect.Int64 {
		return strconv.FormatInt(configField.Int(), field.Base), nil
	} else if fieldKind == reflect.Bool {
		return PrintBool(configField.Bool()), nil
	} else if fieldKind >= reflect.Slice {
		sliceLen := configField.Len()
		sliceValue := configField.Slice(0, sliceLen)
		if field.Type.Elem().Kind() != reflect.String {
			return "", errors.New("field format is not supported")
		}
		slice := sliceValue.Interface().([]string)
		return strings.Join(slice, ","), nil
	}
	return "", fmt.Errorf("field kind %x is not supported", fieldKind)

}
