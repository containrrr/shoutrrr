package format

import (
	"errors"
	"net/url"
	"testing"

	"github.com/fatih/color"

	"github.com/nicholas-fedor/shoutrrr/pkg/services/standard"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestFormat(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Shoutrrr Format Suite")
}

var _ = ginkgo.BeforeSuite(func() {
	// Disable color output for tests to have them match the string format rather than the colors
	color.NoColor = true
})

var _ = ginkgo.Describe("the format package", func() {
	ginkgo.Describe("Generic Format Utils", func() {
		ginkgo.When("parsing a bool", func() {
			testParseValidBool := func(raw string, expected bool) {
				parsed, ok := ParseBool(raw, !expected)
				gomega.Expect(parsed).To(gomega.Equal(expected))
				gomega.Expect(ok).To(gomega.BeTrue())
			}
			ginkgo.It("should parse truthy values as true", func() {
				testParseValidBool("true", true)
				testParseValidBool("1", true)
				testParseValidBool("yes", true)
			})
			ginkgo.It("should parse falsy values as false", func() {
				testParseValidBool("false", false)
				testParseValidBool("0", false)
				testParseValidBool("no", false)
			})
			ginkgo.It("should match regardless of case", func() {
				testParseValidBool("trUE", true)
			})
			ginkgo.It("should return the default if no value matches", func() {
				parsed, ok := ParseBool("bad", true)
				gomega.Expect(parsed).To(gomega.Equal(true))
				gomega.Expect(ok).To(gomega.BeFalse())
				parsed, ok = ParseBool("values", false)
				gomega.Expect(parsed).To(gomega.Equal(false))
				gomega.Expect(ok).To(gomega.BeFalse())
			})
		})
		ginkgo.When("printing a bool", func() {
			ginkgo.It("should return yes or no", func() {
				gomega.Expect(PrintBool(true)).To(gomega.Equal("Yes"))
				gomega.Expect(PrintBool(false)).To(gomega.Equal("No"))
			})
		})
		ginkgo.When("checking for number-like strings", func() {
			ginkgo.It("should be true for numbers", func() {
				gomega.Expect(IsNumber("1.5")).To(gomega.BeTrue())
				gomega.Expect(IsNumber("0")).To(gomega.BeTrue())
				gomega.Expect(IsNumber("NaN")).To(gomega.BeTrue())
			})
			ginkgo.It("should be false for non-numbers", func() {
				gomega.Expect(IsNumber("baNaNa")).To(gomega.BeFalse())
			})
		})
	})
	ginkgo.Describe("Enum Formatter", func() {
		ginkgo.It("should return all enum values on listing", func() {
			gomega.Expect(testEnum.Names()).To(gomega.ConsistOf("None", "Foo", "Bar"))
		})
	})
})

type testStruct struct {
	Signed          int `default:"0" key:"signed"`
	Unsigned        uint
	Str             string `default:"notempty" key:"str"`
	StrSlice        []string
	StrArray        [3]string
	Sub             subStruct
	TestEnum        int `default:"None" key:"testenum"`
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

func (t *testStruct) Enums() map[string]types.EnumFormatter {
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
	enums    = map[string]types.EnumFormatter{
		"TestEnum": testEnum,
	}
)

type testStructBadDefault struct {
	standard.EnumlessConfig
	Value int `default:"NaN" key:"value"`
}

func (t *testStructBadDefault) GetURL() *url.URL {
	panic("not implemented")
}

func (t *testStructBadDefault) SetURL(_ *url.URL) error {
	panic("not implemented")
}
