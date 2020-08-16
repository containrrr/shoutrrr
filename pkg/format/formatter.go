package format

import (
	"fmt"
	"github.com/fatih/color"
	"reflect"
	"strings"
	"unsafe"

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

func GetConfigFormat(service types.Service) (reflect.Type, []fieldInfo) {
	configRef := reflect.ValueOf(service).Elem().FieldByName("config")
	configType := configRef.Type().Elem()

	config := reflect.New(configType)

	serviceConfig := config.Interface().(types.ServiceConfig)

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
				value, valueLen = fmtr.getFieldValueString(values.Field(i), nextDepth)
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

type fieldInfo struct {
	Name          string
	Type          reflect.Type
	EnumFormatter types.EnumFormatter
	Description   string
	DefaultValue  string
	Template      string
	Required      bool
}

func (fmtr *formatter) getStructFieldInfo(structType reflect.Type) []fieldInfo {

	numFields := structType.NumField()
	fields := make([]fieldInfo, numFields)
	maxKeyLen := 0

	for i := 0; i < numFields; i++ {
		fieldDef := structType.Field(i)

		if fieldDef.Anonymous {
			// This is an embedded field, which should not be part of the Config output
			continue
		}

		info := fieldInfo{
			Name:     fieldDef.Name,
			Type:     fieldDef.Type,
			Required: true,
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

		if ef, isEnum := fmtr.EnumFormatters[fieldDef.Name]; isEnum {
			info.EnumFormatter = ef
		}

		fields[i] = info
		keyLen := len(fieldDef.Name)
		if keyLen > maxKeyLen {
			maxKeyLen = keyLen
		}
	}

	return fields
}

func (fmtr *formatter) getFieldValueString(field reflect.Value, depth uint8) (string, int) {

	nextDepth := depth + 1
	kind := field.Kind()

	if util.IsUnsignedDecimal(kind) {
		strVal := fmt.Sprintf("%d", field.Uint())
		return ColorizeNumber(fmt.Sprintf("%s", strVal)), len(strVal)
	}
	if util.IsSignedDecimal(kind) {
		strVal := fmt.Sprintf("%d", field.Int())
		return ColorizeNumber(fmt.Sprintf("%s", strVal)), len(strVal)
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
			items[i], itemLen = fmtr.getFieldValueString(field.Index(i), nextDepth)
			totalLen += itemLen
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
			key, keyLen := fmtr.getFieldValueString(iter.Key(), nextDepth)
			value, valueLen := fmtr.getFieldValueString(iter.Value(), nextDepth)
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
