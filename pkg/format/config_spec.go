package format

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util"
)

type ConfigSpec struct {
	Version uint32
	Props   map[string]*ConfigSpecProp
}

type ConfigSpecProp struct {
	Type         ConfigPropType
	Description  string
	Name         string    `yaml:"-"`
	DefaultValue string    `yaml:"default,omitempty"`
	Template     string    `yaml:",omitempty"`
	Required     bool      `yaml:",omitempty"`
	URLParts     []URLPart `yaml:",omitempty,flow"`
	Title        bool      `yaml:",omitempty"`
	Base         int       `yaml:",omitempty"`
	Keys         []string  `yaml:",omitempty,flow"`
	Values       []string  `yaml:",omitempty,flow"`
	CustomType   string    `yaml:"customType,omitempty"`
}

type ConfigPropType string

const (
	NumberPropType   ConfigPropType = "number"
	TextPropType     ConfigPropType = "text"
	ListPropType     ConfigPropType = "list"
	ColorPropType    ConfigPropType = "color"
	DurationPropType ConfigPropType = "duration"
	OptionPropType   ConfigPropType = "option"
	TogglePropType   ConfigPropType = "toggle"
	MapPropType      ConfigPropType = "map"
	CustomPropType   ConfigPropType = "custom"
)

func (pt ConfigPropType) ParserCall(sp *ConfigSpecProp, valueVar string) string {
	nm := string(pt)
	args := []string{valueVar}
	switch pt {
	case NumberPropType:
		args = append(args, fmt.Sprintf("%d", sp.Base))
	case OptionPropType:
		args = append(args, fmt.Sprintf("%vFormatter", sp.Name))
	}

	return fmt.Sprintf("Parse%v%vValue(%v)", strings.ToUpper(nm[:1]), nm[1:], strings.Join(args, ", "))
}

func ConfigPropTypeFromType(rtype reflect.Type, ttype NodeTokenType) ConfigPropType {
	kind := rtype.Kind()
	if util.IsNumeric(kind) {
		return NumberPropType
	}
	switch kind {
	case reflect.Array:
	case reflect.Slice:
		return ListPropType
	case reflect.String:
		return TextPropType
	case reflect.Bool:
		return TogglePropType
	case reflect.Map:
		return MapPropType
	case reflect.Struct:
		if ttype == PropToken {
			return CustomPropType
		}
	}
	panic(fmt.Sprintf("Invalid config prop type %q (%q)", rtype.String(), kind.String()))
}

type Option int

func ParseNumberValue(v string, base int) (int64, error) {
	return strconv.ParseInt(v, base, 64)
}

func ParseTextValue(v string) (string, error) {
	return v, nil
}

func ParseListValue(v string) ([]string, error) {
	return strings.Split(v, ","), nil
}

func ParseColorValue(v string) (int, error) {
	return 0, fmt.Errorf("color value parser is not implemented")
}

func ParseDurationValue(v string) (time.Duration, error) {
	return time.ParseDuration(v)
}

func ParseOptionValue(v string, ef types.EnumFormatter) (Option, error) {
	val := Option(ef.Parse(v))
	if val == Option(EnumInvalid) {
		return val, fmt.Errorf("invalid option value %q", v)
	}
	return val, nil
}

func ParseToggleValue(v string) (bool, error) {
	if val, ok := ParseBool(v, false); ok {
		return val, nil
	} else {
		return false, fmt.Errorf("invalid toggle value %q", v)
	}
}

func ParseMapValue(v string) (map[string]string, error) {
	return nil, fmt.Errorf("map value parser is not implemented")
}
