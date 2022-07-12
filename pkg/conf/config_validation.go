package conf

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	ipPortMinimum int64 = 1
	ipPortMaximum int64 = 65535
)

type RangeValidator struct {
	Minimum *int64
	Maximum *int64
	Known   string
}

func (v *RangeValidator) Verify() error {
	switch strings.ToLower(v.Known) {
	case "ipport", "tcpport", "udpport", "port":
		v.Minimum = &ipPortMinimum
		v.Maximum = &ipPortMaximum

	case "":
		//
	default:
		return fmt.Errorf("unknown range %q", v.Known)
	}

	if v.Minimum != nil && v.Maximum != nil && *v.Maximum < *v.Minimum {
		return fmt.Errorf("invalid range %v-%v, minimum cannot to be larger than maximum", v.Minimum, v.Maximum)
	}

	if v.Minimum == nil && v.Maximum == nil {
		return fmt.Errorf("no minimum or maximum values specified for range validation")
	}

	return nil
}
func (v *RangeValidator) TestCall(sp *SpecProp, variable string) string {
	cond := strings.Builder{}
	if v.Minimum != nil {
		fmt.Fprintf(&cond, `%v < %v`, variable, *v.Minimum)
	}
	if v.Maximum != nil {
		if v.Minimum != nil {
			cond.WriteString(` || `)
		}
		fmt.Fprintf(&cond, `%v > %v`, variable, *v.Maximum)
	}

	return cond.String()
}
func (v *RangeValidator) FailMessage(propName string, variable string) string {
	min := "*"
	if v.Minimum != nil {
		min = strconv.FormatInt(*v.Minimum, 10)
	}
	max := "*"
	if v.Maximum != nil {
		max = strconv.FormatInt(*v.Maximum, 10)
	}
	return fmt.Sprintf(`"value %%v for %v is not in the range %v-%v", %v`, propName, min, max, variable)
}

type RegexValidator struct {
	Pattern string `yaml:"matches"`
	Known   string
}

const uuid4Pattern = "[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}"
const hex32Pattern = "[A-Za-z0-9]{32}"

func (v *RegexValidator) Verify() (err error) {
	switch v.Known {
	case "uuid", "uuid4", "guid":
		v.Pattern = uuid4Pattern
	case "hex32":
		v.Pattern = hex32Pattern
	}
	_, err = regexp.Compile(v.Pattern)
	return
}
func (v *RegexValidator) TestCall(sp *SpecProp, value string) string {
	return fmt.Sprintf(`!conf.ValueMatchesPattern(%v, %q)`, value, v.Pattern)
}
func (v *RegexValidator) FailMessage(propName string, variable string) string {
	return fmt.Sprintf(`"value %%v for %v does not match the expected format", %v`, propName, variable)
}

type NotEmptyValidator struct{}

func (v *NotEmptyValidator) Verify() error { return nil }
func (v *NotEmptyValidator) TestCall(sp *SpecProp, value string) string {
	return sp.Type.EmptyCall(sp, value)
}
func (v *NotEmptyValidator) FailMessage(propName string, _ string) string {
	return fmt.Sprintf(`"%v missing from config URL"`, propName)
}

type LengthValidator struct {
	Minimum *int64
	Maximum *int64
	Equals  *int64
}

func (v *LengthValidator) Verify() error {

	if v.Equals != nil {
		if v.Minimum != nil && v.Maximum != nil {
			return fmt.Errorf("minimum/maximum length cannot be specified together with equals")
		}
		v.Minimum = v.Equals
		v.Maximum = v.Equals
	}

	if v.Minimum != nil && v.Maximum != nil && *v.Maximum < *v.Minimum {
		return fmt.Errorf("invalid range %v-%v, minimum cannot to be larger than maximum", v.Minimum, v.Maximum)
	}

	if v.Minimum == nil && v.Maximum == nil {
		return fmt.Errorf("no minimum or maximum values specified for length validation")
	}

	return nil
}
func (v *LengthValidator) TestCall(sp *SpecProp, variable string) string {
	cond := strings.Builder{}
	if v.Minimum != nil {
		fmt.Fprintf(&cond, `len(%v) < %v`, variable, *v.Minimum)
	}
	if v.Maximum != nil {
		if v.Minimum != nil {
			if *v.Minimum == *v.Maximum {
				return fmt.Sprintf(`len(%v) != %v`, variable, *v.Minimum)
			}
			cond.WriteString(` || `)
		}
		fmt.Fprintf(&cond, `len(%v) > %v`, variable, *v.Maximum)
	}

	return cond.String()
}
func (v *LengthValidator) FailMessage(propName string, variable string) string {
	min := "*"
	if v.Minimum != nil {
		min = strconv.FormatInt(*v.Minimum, 10)
	}
	max := "*"
	if v.Maximum != nil {
		max = strconv.FormatInt(*v.Minimum, 10)
	}
	if min == max {
		return fmt.Sprintf(`"value %%q for %v is not the correct length (%v)", %v`, propName, min, variable)
	}
	return fmt.Sprintf(`"value %%q for %v length is not in the range %v-%v", %v`, propName, min, max, variable)
}

type NotEqualToPropValidator string

func (v NotEqualToPropValidator) Verify() (err error) {
	return nil
}
func (v NotEqualToPropValidator) TestCall(sp *SpecProp, value string) string {
	return fmt.Sprintf(`%v == config.%v`, value, string(v))
}
func (v NotEqualToPropValidator) FailMessage(propName string, variable string) string {
	return fmt.Sprintf(`"value %%v for %v is already used for %v", %v`, propName, v, variable)
}
