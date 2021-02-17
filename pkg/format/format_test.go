package format

import (
	"errors"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/fatih/color"
	"reflect"

	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	// logger *log.Logger
	f = formatter{
		EnumFormatters: map[string]types.EnumFormatter{
			"TestEnum": testEnum,
		},
		MaxDepth: 2,
	}
	ts *testStruct
)

func TestFormat(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Discord Suite")
}

var _ = Describe("the format package", func() {
	BeforeSuite(func() {
		// logger = log.New(GinkgoWriter, "Test", log.LstdFlags)

		// Disable color output for tests to have them match the string format rather than the colors
		color.NoColor = true
	})

	Describe("SetConfigField", func() {
		var (
			tv reflect.Value
		)
		tt := reflect.TypeOf(testStruct{})
		fields := f.getStructFieldInfo(tt)

		fieldMap := make(map[string]FieldInfo, len(fields))
		for _, fi := range fields {
			fieldMap[fi.Name] = fi
		}
		When("updating a struct", func() {

			BeforeEach(func() {
				tsPtr := reflect.New(tt)
				tv = tsPtr.Elem()
				ts = tsPtr.Interface().(*testStruct)
			})

			When("setting an integer value", func() {
				When("the value is valid", func() {
					It("should set it", func() {
						valid, err := SetConfigField(tv, fieldMap["Signed"], "3")
						Expect(valid).To(BeTrue())
						Expect(err).NotTo(HaveOccurred())

						Expect(ts.Signed).To(Equal(3))
					})
				})
				When("the value is invalid", func() {
					It("should return an error", func() {
						ts.Signed = 2
						valid, err := SetConfigField(tv, fieldMap["Signed"], "z7")
						Expect(valid).To(BeFalse())
						Expect(err).To(HaveOccurred())

						Expect(ts.Signed).To(Equal(2))
					})
				})
			})

			When("setting an unsigned integer value", func() {
				When("the value is valid", func() {
					It("should set it", func() {
						valid, err := SetConfigField(tv, fieldMap["Unsigned"], "6")
						Expect(valid).To(BeTrue())
						Expect(err).NotTo(HaveOccurred())

						Expect(ts.Unsigned).To(Equal(uint(6)))
					})
				})
				When("the value is invalid", func() {
					It("should return an error", func() {
						ts.Unsigned = 2
						valid, err := SetConfigField(tv, fieldMap["Unsigned"], "-3")

						Expect(ts.Unsigned).To(Equal(uint(2)))
						Expect(valid).To(BeFalse())
						Expect(err).To(HaveOccurred())
					})
				})
			})

			When("setting a string slice value", func() {
				When("the value is valid", func() {
					It("should set it", func() {
						valid, err := SetConfigField(tv, fieldMap["StrSlice"], "meawannowalkalitabitalleh,meawannofeelalitabitstrongah")
						Expect(valid).To(BeTrue())
						Expect(err).NotTo(HaveOccurred())

						Expect(ts.StrSlice).To(HaveLen(2))
					})
				})
			})

			When("setting a string array value", func() {
				When("the value is valid", func() {
					It("should set it", func() {
						valid, err := SetConfigField(tv, fieldMap["StrArray"], "meawannowalkalitabitalleh,meawannofeelalitabitstrongah,meawannothinkalitabitsmartah")
						Expect(valid).To(BeTrue())
						Expect(err).NotTo(HaveOccurred())
					})
				})
				When("the value has too many elements", func() {
					It("should return an error", func() {
						valid, err := SetConfigField(tv, fieldMap["StrArray"], "one,two,three,four?")
						Expect(valid).To(BeFalse())
						Expect(err).To(HaveOccurred())
					})
				})
				When("the value has too few elements", func() {
					It("should return an error", func() {
						valid, err := SetConfigField(tv, fieldMap["StrArray"], "one,two")
						Expect(valid).To(BeFalse())
						Expect(err).To(HaveOccurred())
					})
				})
			})

			When("setting a struct value", func() {
				When("it implements ConfigProp", func() {
					It("should return an error", func() {
						valid, err := SetConfigField(tv, fieldMap["Sub"], "@awol")
						Expect(err).To(HaveOccurred())
						Expect(valid).NotTo(BeTrue())
					})
				})
				When("it implements ConfigProp", func() {
					When("the value is valid", func() {
						It("should set it", func() {
							valid, err := SetConfigField(tv, fieldMap["SubProp"], "@awol")
							Expect(err).NotTo(HaveOccurred())
							Expect(valid).To(BeTrue())

							Expect(ts.SubProp.Value).To(Equal("awol"))
						})
					})
					When("the value is invalid", func() {
						It("should return an error", func() {
							valid, err := SetConfigField(tv, fieldMap["SubProp"], "missing initial at symbol")
							Expect(err).To(HaveOccurred())
							Expect(valid).NotTo(BeTrue())
						})
					})
				})
			})

			When("setting a struct slice value", func() {
				When("the value is valid", func() {
					It("should set it", func() {
						valid, err := SetConfigField(tv, fieldMap["SubPropSlice"], "@alice,@merton")
						Expect(err).NotTo(HaveOccurred())
						Expect(valid).To(BeTrue())

						Expect(ts.SubPropSlice).To(HaveLen(2))
					})
				})
			})

			When("setting a struct pointer slice value", func() {
				When("the value is valid", func() {
					It("should set it", func() {
						valid, err := SetConfigField(tv, fieldMap["SubPropPtrSlice"], "@the,@best")
						Expect(err).NotTo(HaveOccurred())
						Expect(valid).To(BeTrue())

						Expect(ts.SubPropPtrSlice).To(HaveLen(2))
					})
				})
			})
		})

		When("formatting stuct values", func() {
			BeforeEach(func() {
				tsPtr := reflect.New(tt)
				tv = tsPtr.Elem()
			})
			When("setting and formatting", func() {
				It("should format signed integers identical to input", func() {
					testSetAndFormat(tv, fieldMap["Signed"], "-45", "-45")
				})
				It("should format unsigned integers identical to input", func() {
					testSetAndFormat(tv, fieldMap["Unsigned"], "5", "5")
				})
				It("should format structs identical to input", func() {
					testSetAndFormat(tv, fieldMap["SubProp"], "@whoa", "@whoa")
				})
				It("should format enums identical to input", func() {
					testSetAndFormat(tv, fieldMap["TestEnum"], "Foo", "Foo")
				})
				It("should format string slices identical to input", func() {
					testSetAndFormat(tv, fieldMap["StrSlice"], "one,two,three,four", "[ one, two, three, four ]")
				})
				It("should format string arrays identical to input", func() {
					testSetAndFormat(tv, fieldMap["StrArray"], "one,two,three", "[ one, two, three ]")
				})
				It("should format prop struct slices identical to input", func() {
					testSetAndFormat(tv, fieldMap["SubPropSlice"], "@be,@the,@best", "[ @be, @the, @best ]")
				})
				It("should format prop struct slices identical to input", func() {
					testSetAndFormat(tv, fieldMap["SubPropPtrSlice"], "@diet,@glue", "[ @diet, @glue ]")
				})
				It("should format prop struct slices identical to input", func() {
					testSetAndFormat(tv, fieldMap["StrMap"], "a:1,c:3,b:2", "{ a: 1, b: 2, c: 3 }")
				})
			})
		})
	})
})

func testSetAndFormat(tv reflect.Value, field FieldInfo, value string, prettyFormat string) {
	_, _ = SetConfigField(tv, field, value)
	fieldValue := tv.FieldByName(field.Name)

	// Used for de-/serializing configuration
	formatted, err := GetConfigFieldString(tv, field)
	Expect(err).NotTo(HaveOccurred())
	Expect(formatted).To(Equal(value))

	// Used for pretty printing output, coloring etc.
	formatted, _ = f.getStructFieldValueString(fieldValue, field, 0)
	Expect(formatted).To(Equal(prettyFormat))
}

type testStruct struct {
	Signed          int
	Unsigned        uint
	Str             string
	StrSlice        []string
	StrArray        [3]string
	Sub             subStruct
	TestEnum        int
	SubProp         subPropStruct
	SubSlice        []subStruct
	SubPropSlice    []subPropStruct
	SubPropPtrSlice []*subPropStruct
	StrMap          map[string]string
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

var testEnum = CreateEnumFormatter([]string{"None", "Foo", "Bar"})
