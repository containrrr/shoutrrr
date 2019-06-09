package format

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/types"
    "github.com/fatih/color"
	"reflect"
	"strings"
)

// GetConfigMap returns a string map of a given Config struct
func GetConfigMap(config types.ServiceConfig) map[string]string {
	formatter := formatter{
		EnumFormatters: config.Enums(),
		MaxDepth: 10,
	}
	return formatter.getStructMap(config, 0)
}

var colorizeDesc = color.New(color.FgHiBlack).SprintFunc()
var colorizeTrue = color.New(color.FgHiGreen).SprintFunc()
var colorizeFalse = color.New(color.FgHiRed).SprintFunc()
var colorizeNumber = color.New(color.FgHiCyan).SprintFunc()
var colorizeString = color.New(color.FgHiYellow).SprintFunc()

type formatter struct {
	EnumFormatters map[string]types.EnumFormatter
	MaxDepth uint8
	Errors []error
}

func (fmtr *formatter) getStructMap(structItem interface{}, depth uint8) map[string]string {
	values := reflect.ValueOf(structItem).Elem()
	defs := reflect.TypeOf(structItem).Elem()
	numFields := values.NumField()
	valueMap := make(map[string]string, numFields)
	nextDepth := depth + 1

	for i := 0; i < numFields; i++ {
		fieldDef := defs.Field(i)
		value := fmt.Sprintf("(%s)", fieldDef.Type.Name())
		if ef, isEnum := fmtr.EnumFormatters[fieldDef.Name]; isEnum {
			fieldVal := values.Field(i);
			kind := fieldVal.Kind()
			if kind == reflect.Int {
				value = ef.Print(int(fieldVal.Int()))
			} else {
				err := fmt.Errorf("incorrect enum type '%s' for field '%s'", kind, fieldDef.Name)
				fmtr.Errors = append(fmtr.Errors, err)
			}
		} else if nextDepth < fmtr.MaxDepth {
			value = fmtr.getFieldValueString(values.Field(i), nextDepth)
		}

		if desc, ok := fieldDef.Tag.Lookup("desc"); ok {
			value += " " + colorizeDesc(desc)
		}
		valueMap[fieldDef.Name] = value
	}

	return valueMap
}

func (fmtr *formatter) getFieldValueString(field reflect.Value, depth uint8) string {

	nextDepth := depth + 1
	kind := field.Kind()

	switch kind {
	case reflect.String:
		return colorizeString(field.String())

	case reflect.Int8: fallthrough
	case reflect.Int16: fallthrough
	case reflect.Int32: fallthrough
	case reflect.Int64: fallthrough
	case reflect.Int:
		return colorizeNumber(fmt.Sprintf("%d", field.Int()))

	case reflect.Uint8: fallthrough
	case reflect.Uint16: fallthrough
	case reflect.Uint32: fallthrough
	case reflect.Uint64: fallthrough
	case reflect.Uint:
		return colorizeNumber(fmt.Sprintf("%d", field.Uint()))

	case reflect.Bool:
		val := field.Bool()
		if val {
			return colorizeTrue(PrintBool(val))
		}
		return colorizeFalse(PrintBool(val))


	case reflect.Slice:
		fallthrough
	case reflect.Array:
		len := field.Len()
		items := make([]string, len)
		for i := 0; i < field.Len(); i++ {
			items[i] = fmtr.getFieldValueString(field.Index(i), nextDepth)
		}
		return fmt.Sprintf("[ %s ]", strings.Join(items, ", "))

	case reflect.Map:
		items := make([]string, field.Len())
		iter := field.MapRange()
		index := 0
		for iter.Next() {
			key := fmtr.getFieldValueString(iter.Key(), nextDepth)
			value := fmtr.getFieldValueString(iter.Value(), nextDepth)
			items[index] = fmt.Sprintf("%s: %s", key, value)
		}
		return fmt.Sprintf("{ %s }", strings.Join(items, ", "))

	case reflect.Struct:
		structMap := fmtr.getStructMap(field, depth + 1)
		items := make([]string, len(structMap))
		index := 0
		for key, value := range structMap {
			items[index] = fmt.Sprintf("%s: %s", key, value)
			index++
		}
		return fmt.Sprintf("< %s >", strings.Join(items, ", "))

	default:
		return fmt.Sprintf("<?%s>", kind.String())
	}
}
