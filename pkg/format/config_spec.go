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
	Scheme  string
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
	Credential   bool      `yaml:"credential"`
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

func (pt ConfigPropType) FormatCall(sp *ConfigSpecProp, valueVar string) string {
	nm := string(pt)
	args := []string{valueVar}
	switch pt {
	case NumberPropType:
		args = append(args, fmt.Sprintf("%d", sp.Base))
	case CustomPropType:
		return fmt.Sprintf("config.get%v()", sp.Name)
	case OptionPropType:
		return fmt.Sprintf(`%vOptions.Formatter.Print(int(%v))`, sp.Name, args[0])
		// args = append(args, fmt.Sprintf("%vOptions.Formatter", sp.Name))
	}

	return fmt.Sprintf("format.Format%v%vValue(%v)", strings.ToUpper(nm[:1]), nm[1:], strings.Join(args, ", "))
}

func (pt ConfigPropType) ParserCall(sp *ConfigSpecProp, valueVar string) string {
	nm := string(pt)
	args := []string{valueVar}
	switch pt {
	case NumberPropType:
		args = append(args, fmt.Sprintf("%d", sp.Base))
	case CustomPropType:
		return fmt.Sprintf("config.set%v(%v)", sp.Name, valueVar)
	case OptionPropType:
		return fmt.Sprintf(`%vOptions.Parse(%v)`, sp.Name, valueVar)
		// args = append(args, fmt.Sprintf("%vOptions.Formatter", sp.Name))
	}

	return fmt.Sprintf("format.Parse%v%vValue(%v)", strings.ToUpper(nm[:1]), nm[1:], strings.Join(args, ", "))
}

func (pt ConfigPropType) EmptyCall(sp *ConfigSpecProp, valueVar string) string {
	switch pt {
	case NumberPropType:
		return fmt.Sprintf(`%v == 0`, valueVar)
	case TextPropType:
		return fmt.Sprintf(`%v == ""`, valueVar)
	case ListPropType:
		return fmt.Sprintf(`len(%v) == 0`, valueVar)
	case OptionPropType:
		return fmt.Sprintf(`int(%v) == 0`, valueVar)
	// args = append(args, fmt.Sprintf("%vOptions.Formatter", sp.Name))
	default:
		panic(fmt.Sprintf("EmptyCall not implemented for %v", pt))
	}
}

func ConfigPropTypeFromType(rtype reflect.Type, ttype NodeTokenType) ConfigPropType {
	if ttype == EnumToken {
		return OptionPropType
	}
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
	if v == "" {
		// If the source string is empty, the list contains no items
		return []string{}, nil
	}
	return strings.Split(v, ","), nil
}

func ParseColorValue(v string) (uint32, error) {
	if len(v) > 0 && v[0] == '#' {
		v = v[1:]
	}
	if len(v) > 1 && v[:2] == "0x" {
		v = v[2:]
	}
	color, err := strconv.ParseUint(v, 16, 32)
	return uint32(color), err
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

func FormatNumberValue(v int64, base int) string {
	return strconv.FormatInt(v, base)
}

func FormatTextValue(v string) string {
	return v
}

func FormatListValue(v []string) string {
	return strings.Join(v, ",")
}

func FormatColorValue(v uint32) string {
	return fmt.Sprintf("0x%06x", v)
}

func FormatDurationValue(v time.Duration) string {
	return v.String()
}

func FormatOptionValue(v Option, ef types.EnumFormatter) string {
	return ef.Print(int(v))
}

func FormatToggleValue(v bool) string {
	return PrintBool(v)
}

func FormatMapValue(v map[string]string) string {
	panic("map value parser is not implemented")
}
