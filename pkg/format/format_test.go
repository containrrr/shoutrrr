package format

import (
	"errors"
	"net/url"
	"testing"

	"github.com/fatih/color"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
	t "github.com/containrrr/shoutrrr/pkg/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFormat(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Format Suite")
}

var _ = BeforeSuite(func() {
	// Disable color output for tests to have them match the string format rather than the colors
	color.NoColor = true
})

var _ = Describe("the format package", func() {
	Describe("Generic Format Utils", func() {
		When("parsing a bool", func() {
			var testParseValidBool = func(raw string, expected bool) {
				parsed, ok := ParseBool(raw, !expected)
				Expect(parsed).To(Equal(expected))
				Expect(ok).To(BeTrue())
			}
			It("should parse truthy values as true", func() {
				testParseValidBool("true", true)
				testParseValidBool("1", true)
				testParseValidBool("yes", true)
			})
			It("should parse falsy values as false", func() {
				testParseValidBool("false", false)
				testParseValidBool("0", false)
				testParseValidBool("no", false)
			})
			It("should match regardless of case", func() {
				testParseValidBool("trUE", true)
			})
			It("should return the default if no value matches", func() {
				parsed, ok := ParseBool("bad", true)
				Expect(parsed).To(Equal(true))
				Expect(ok).To(BeFalse())
				parsed, ok = ParseBool("values", false)
				Expect(parsed).To(Equal(false))
				Expect(ok).To(BeFalse())
			})
		})
		When("printing a bool", func() {
			It("should return yes or no", func() {
				Expect(PrintBool(true)).To(Equal("Yes"))
				Expect(PrintBool(false)).To(Equal("No"))
			})
		})
		When("checking for number-like strings", func() {
			It("should be true for numbers", func() {
				Expect(IsNumber("1.5")).To(BeTrue())
				Expect(IsNumber("0")).To(BeTrue())
				Expect(IsNumber("NaN")).To(BeTrue())
			})
			It("should be false for non-numbers", func() {
				Expect(IsNumber("baNaNa")).To(BeFalse())
			})
		})
	})
	Describe("Enum Formatter", func() {
		It("should return all enum values on listing", func() {
			Expect(testEnum.Names()).To(ConsistOf("None", "Foo", "Bar"))
		})
	})
})

type testStruct struct {
	Signed          int `key:"signed" default:"0"`
	Unsigned        uint
	Str             string `key:"str" default:"notempty"`
	StrSlice        []string
	StrArray        [3]string
	Sub             subStruct
	TestEnum        int `key:"testenum" default:"None"`
	SubProp         subPropStruct
	SubSlice        []subStruct
	SubPropSlice    []subPropStruct
	SubPropPtrSlice []*subPropStruct
	StrMap          map[string]string
	IntMap          map[string]int
	Int8Map         map[string]int8
	Int16Map        map[string]int16
	Int32Map        map[string]int32
	Int64Map        map[string]int64
	UintMap         map[string]uint
	Uint8Map        map[string]int8
	Uint16Map       map[string]int16
	Uint32Map       map[string]int32
	Uint64Map       map[string]int64
}

func (t *testStruct) GetURL() *url.URL {
	panic("not implemented")
}

func (t *testStruct) SetURL(_ *url.URL) error {
	panic("not implemented")
}

func (t *testStruct) Enums() map[string]t.EnumFormatter {
	return enums
}

type subStruct struct {
	Value string
}

type subPropStruct struct {
	Value string
}

func (s *subPropStruct) SetFromProp(propValue string) error {
	if len(propValue) < 1 || propValue[0] != '@' {
		return errors.New("invalid value")
	}
	s.Value = propValue[1:]
	return nil
}
func (s *subPropStruct) GetPropValue() (string, error) {
	return "@" + s.Value, nil
}

var (
	testEnum = CreateEnumFormatter([]string{"None", "Foo", "Bar"})
	enums    = map[string]t.EnumFormatter{
		"TestEnum": testEnum,
	}
)

type testStructBadDefault struct {
	standard.EnumlessConfig
	Value int `key:"value" default:"NaN"`
}

func (t *testStructBadDefault) GetURL() *url.URL {
	panic("not implemented")
}

func (t *testStructBadDefault) SetURL(_ *url.URL) error {
	panic("not implemented")
}
