package conf

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/ref"
	up "github.com/containrrr/shoutrrr/pkg/urlpart"
	"github.com/containrrr/shoutrrr/pkg/util"
)

type Spec struct {
	Version uint32
	Scheme  string
	Options struct {
		ReversePathPrio bool `yaml:"reversePathPrio,omitempty"`
		CustomQueryVars bool `yaml:"customQueryVars,omitempty"`
	} `yaml:",omitempty"`
	Props map[string]*SpecProp
}

type SpecProp struct {
	Type           PropType
	Description    string
	Name           string       `yaml:"-"`
	DefaultValue   string       `yaml:"default,omitempty"`
	Template       string       `yaml:",omitempty"`
	Required       bool         `yaml:",omitempty"`
	URLParts       []up.URLPart `yaml:",omitempty,flow"`
	Title          bool         `yaml:",omitempty"`
	Base           int          `yaml:",omitempty"`
	Keys           []string     `yaml:",omitempty,flow"`
	Values         []string     `yaml:",omitempty,flow"`
	CustomType     string       `yaml:"customType,omitempty"`
	Credential     bool         `yaml:"credential"`
	ItemSeparator  string       `yaml:"itemSeparator,omitempty"`
	ValueSeparator string       `yaml:"valueSeparator,omitempty"`
	Validation     struct {
		ValidRange   *RangeValidator          `yaml:"range,omitempty"`
		ValidLength  *LengthValidator         `yaml:"length,omitempty"`
		MatchesRegex *RegexValidator          `yaml:"regex,omitempty"`
		NotEqual     *NotEqualToPropValidator `yaml:"notEqual,omitempty"`
	} `yaml:"validation"`
}

type PropValidator interface {
	TestCall(sp *SpecProp, variable string) string
	FailMessage(propName, variable string) string
	Verify() error
}

func (sp *SpecProp) Validators() []PropValidator {
	validators := []PropValidator{}
	if sp.Validation.MatchesRegex != nil {
		validators = append(validators, sp.Validation.MatchesRegex)
	}
	if sp.Validation.ValidRange != nil {
		validators = append(validators, sp.Validation.ValidRange)
	}
	if sp.Validation.ValidLength != nil {
		validators = append(validators, sp.Validation.ValidLength)
	}
	if sp.Validation.NotEqual != nil {
		validators = append(validators, sp.Validation.NotEqual)
	}
	if sp.Required {
		validators = append(validators, &NotEmptyValidator{})
	}
	return validators
}

type PropType string

const (
	NumberPropType   PropType = "number"
	TextPropType     PropType = "text"
	ListPropType     PropType = "list"
	ColorPropType    PropType = "color"
	DurationPropType PropType = "duration"
	OptionPropType   PropType = "option"
	TogglePropType   PropType = "toggle"
	MapPropType      PropType = "map"
	CustomPropType   PropType = "custom"
)

func (pt PropType) FormatCall(sp *SpecProp, valueVar string) string {
	nm := string(pt)
	args := []string{valueVar}
	switch pt {
	case NumberPropType:
		args = append(args, fmt.Sprintf("%d", sp.Base))
	case ListPropType:
		if sp.CustomType != "" {
			args[0] = fmt.Sprintf("format%vItems(%v)", sp.CustomType, valueVar)
		}
		args = addSeparatorArgs(args, sp)
	case MapPropType:
		if len(sp.URLParts) > 0 && sp.URLParts[0] == up.Query {
			return fmt.Sprintf(`%v.Encode()`, valueVar)
		}
		args = addSeparatorArgs(args, sp)
	case CustomPropType:
		return fmt.Sprintf("config.get%v()", sp.Name)
	case OptionPropType:
		return fmt.Sprintf(`%vOptions.Formatter.Print(int(%v))`, sp.Name, args[0])
		// args = append(args, fmt.Sprintf("%vOptions.Formatter", sp.Name))
	}

	return fmt.Sprintf("conf.Format%v%vValue(%v)", strings.ToUpper(nm[:1]), nm[1:], strings.Join(args, ", "))
}

func (pt PropType) ParserCall(sp *SpecProp, valueVar string) string {
	nm := string(pt)
	typeParserName := strings.ToUpper(nm[:1]) + nm[1:]
	args := []string{valueVar}
	switch pt {
	case NumberPropType:
		args = append(args, fmt.Sprintf("%d", sp.Base))
	case CustomPropType:
		return fmt.Sprintf("config.set%v(%v)", sp.Name, valueVar)
	case ListPropType:
		if sp.CustomType != "" {
			return fmt.Sprintf("parse%vItems(%v)", sp.CustomType, valueVar)
		}
		if len(sp.URLParts) == 1 && sp.URLParts[0] == up.Path {
			// If the prop is used to store all path elements, use special parser
			typeParserName = "Path" + typeParserName
		} else {
			args = addSeparatorArgs(args, sp)
		}
	case MapPropType:
		if len(sp.URLParts) > 0 && sp.URLParts[0] == up.Query {
			return fmt.Sprintf(`url.ParseQuery(%v)`, valueVar)
		}
		args = addSeparatorArgs(args, sp)
	case OptionPropType:
		return fmt.Sprintf(`%vOptions.Parse(%v)`, sp.Name, valueVar)
	}

	return fmt.Sprintf("conf.Parse%vValue(%v)", typeParserName, strings.Join(args, ", "))
}

func (pt PropType) EmptyCall(sp *SpecProp, valueVar string) string {
	switch pt {
	case NumberPropType:
		return fmt.Sprintf(`%v == 0`, valueVar)
	case TextPropType:
		return fmt.Sprintf(`%v == ""`, valueVar)
	case ListPropType:
		return fmt.Sprintf(`len(%v) == 0`, valueVar)
	case OptionPropType:
		return fmt.Sprintf(`int(%v) == 0`, valueVar)
	case CustomPropType:
		return fmt.Sprintf("config.empty%v(%v)", sp.Name, valueVar)
	default:
		panic(fmt.Sprintf("EmptyCall not implemented for %v", pt))
	}
}

func (pt PropType) GoType(sp *SpecProp) string {
	switch pt {
	case TextPropType:
		return "string"
	case TogglePropType:
		return "bool"
	case NumberPropType:
		return "int64"
	case OptionPropType:
		return OptionTypeName(sp.Name)
	case ListPropType:
		itemType := "string"
		if sp.CustomType != "" {
			itemType = sp.CustomType
		}
		return fmt.Sprintf("[]%v", itemType)
	case MapPropType:
		if len(sp.URLParts) > 0 && sp.URLParts[0] == up.Query {
			return `url.Values`
		}
		return "map[string]string"
	case ColorPropType:
		return "uint32"
	case CustomPropType:
		return CustomTypeName(sp.Name, sp.CustomType)
	default:
		return "interface{}"
	}
}

func addSeparatorArgs(args []string, sp *SpecProp) []string {
	itemSep := sp.ItemSeparator
	if itemSep == "" {
		itemSep = ","
	}
	args = append(args, fmt.Sprintf(`%q`, itemSep))

	if sp.Type == MapPropType {
		valueSep := sp.ValueSeparator
		if valueSep == "" {
			valueSep = ":"
		}
		args = append(args, fmt.Sprintf(`%q`, valueSep))
	}

	return args
}

func ConfigPropTypeFromType(refType reflect.Type, tokenType ref.NodeTokenType) PropType {
	if tokenType == ref.EnumToken {
		return OptionPropType
	}
	kind := refType.Kind()
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
		if tokenType == ref.PropToken {
			return CustomPropType
		}
	}
	panic(fmt.Sprintf("Invalid config prop type %q (%q)", refType.String(), kind.String()))
}

func OptionTypeName(name string) string {
	return strings.ToLower(name[:1]) + name[1:] + `Option`
}

func CustomTypeName(name string, customType string) string {
	if customType != "" {
		return "*" + customType
	}
	return strings.ToLower(name[:1]) + name[1:] + `Type`
}
