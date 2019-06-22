package format

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util"
	"reflect"
	"strings"
)

// GetConfigMap returns a string map of a given Config struct
func GetConfigMap(config types.ServiceConfig) (map[string]string, int) {
	formatter := formatter{
		EnumFormatters: config.Enums(),
		MaxDepth: 10,
	}
	return formatter.getStructMap(config, 0)
}



type formatter struct {
	EnumFormatters map[string]types.EnumFormatter
	MaxDepth uint8
	Errors []error
}

func (fmtr *formatter) getStructMap(structItem interface{}, depth uint8) (map[string]string, int) {
	values := reflect.ValueOf(structItem).Elem()
	defs := reflect.TypeOf(structItem).Elem()
	numFields := values.NumField()
	valueMap := make(map[string]string, numFields)
	nextDepth := depth + 1
	maxKeyLen := 0

	for i := 0; i < numFields; i++ {
		fieldDef := defs.Field(i)
		value := fmt.Sprintf("(%s)", fieldDef.Type.Name())
		valueLen := len(value)


		ef, isEnum := fmtr.EnumFormatters[fieldDef.Name];
		if isEnum {
			fieldVal := values.Field(i);
			kind := fieldVal.Kind()
			if kind == reflect.Int {
				valueStr := ef.Print(int(fieldVal.Int()))
				value = ColorizeEnum(valueStr)
				valueLen = len(valueStr)
			} else {
				err := fmt.Errorf("incorrect enum type '%s' for field '%s'", kind, fieldDef.Name)
				fmtr.Errors = append(fmtr.Errors, err)
			}
		} else if nextDepth < fmtr.MaxDepth {
			value, valueLen = fmtr.getFieldValueString(values.Field(i), nextDepth)
		}

		if tag, ok := fieldDef.Tag.Lookup("desc"); ok {

			prePad := strings.Repeat(" ", util.Max(40 - valueLen, 1))
			postPad := strings.Repeat(" ", util.Max(60 - len(tag), 1))

			value += " " + prePad + ColorizeDesc(tag) + postPad
		}

		if tag, ok := fieldDef.Tag.Lookup("tpl"); ok {
			value += fmt.Sprintf(" <Template: %s>", ColorizeString(tag))
		}

		if tag, ok := fieldDef.Tag.Lookup("default"); ok {
			value += fmt.Sprintf(" <Default: %s>", ColorizeValue(tag, isEnum))
		}

		if isEnum {
			value += fmt.Sprintf(" [%s]", strings.Join(ef.Names(), ", "))
		}


		valueMap[fieldDef.Name] = value
		keyLen := len(fieldDef.Name)
		if keyLen > maxKeyLen {
			maxKeyLen = keyLen
		}
	}

	return valueMap, maxKeyLen
}

func (fmtr *formatter) getFieldValueString(field reflect.Value, depth uint8) (string, int) {

	nextDepth := depth + 1
	kind := field.Kind()

	switch kind {
	case reflect.String:
		strVal := field.String()
		return ColorizeString(strVal), len(strVal)

	case reflect.Int8: fallthrough
	case reflect.Int16: fallthrough
	case reflect.Int32: fallthrough
	case reflect.Int64: fallthrough
	case reflect.Int:
		strVal := fmt.Sprintf("%d", field.Int())
		return ColorizeNumber(fmt.Sprintf("%s", strVal)), len(strVal)

	case reflect.Uint8: fallthrough
	case reflect.Uint16: fallthrough
	case reflect.Uint32: fallthrough
	case reflect.Uint64: fallthrough
	case reflect.Uint:
		strVal := fmt.Sprintf("%d", field.Uint())
		return ColorizeNumber(fmt.Sprintf("%s", strVal)), len(strVal)

	case reflect.Bool:
		val := field.Bool()
		if val {
			return ColorizeTrue(PrintBool(val)), 3
		}
		return ColorizeFalse(PrintBool(val)), 2


	case reflect.Slice:
		fallthrough
	case reflect.Array:
		len := field.Len()
		items := make([]string, len)
		totalLen := 4
		var itemLen int
		for i := 0; i < field.Len(); i++ {
			items[i], itemLen = fmtr.getFieldValueString(field.Index(i), nextDepth)
			totalLen += itemLen
		}
		return fmt.Sprintf("[ %s ]", strings.Join(items, ", ")), totalLen

	case reflect.Map:
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

	case reflect.Struct:
		structMap, _ := fmtr.getStructMap(field, depth + 1)
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

	default:
		strVal := kind.String()
		return fmt.Sprintf("<?%s>", strVal), len(strVal) + 5
	}
}
