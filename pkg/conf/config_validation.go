package conf

import (
	"fmt"
	"regexp"
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
	case "ipport":
	case "tcpport":
	case "udpport":
	case "port":
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
	return fmt.Sprintf(`"value %%v for %v is not in the range %v-%v", %v`, propName, v.Minimum, v.Maximum, variable)
}

type RegexValidator struct {
	Pattern string `yaml:"matches"`
}

func (v *RegexValidator) Verify() (err error) {
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
	// return fmt.Sprintf(`!(%v)`, sp.Type.EmptyCall(sp, value))
	return sp.Type.EmptyCall(sp, value)
}
func (v *NotEmptyValidator) FailMessage(propName string, _ string) string {
	return fmt.Sprintf(`"%v missing from config URL"`, propName)
}

//
