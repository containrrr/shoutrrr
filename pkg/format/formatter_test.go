package format

import (
	"reflect"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	// logger *log.Logger
	ts *testStruct
)

var _ = Describe("SetConfigField", func() {
	var (
		tv reflect.Value
	)
	testConfig := testStruct{}
	tt := reflect.TypeOf(testConfig)
	rootNode := getRootNode(&testConfig)

	nodeMap := make(map[string]Node, len(rootNode.Items))
	for _, item := range rootNode.Items {
		field := item.Field()
		nodeMap[field.Name] = item
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
					valid, err := SetConfigField(tv, *nodeMap["Signed"].Field(), "3")
					Expect(valid).To(BeTrue())
					Expect(err).NotTo(HaveOccurred())

					Expect(ts.Signed).To(Equal(3))
				})
			})
			When("the value is invalid", func() {
				It("should return an error", func() {
					ts.Signed = 2
					valid, err := SetConfigField(tv, *nodeMap["Signed"].Field(), "z7")
					Expect(valid).To(BeFalse())
					Expect(err).To(HaveOccurred())

					Expect(ts.Signed).To(Equal(2))
				})
			})
		})

		When("setting an unsigned integer value", func() {
			When("the value is valid", func() {
				It("should set it", func() {
					valid, err := SetConfigField(tv, *nodeMap["Unsigned"].Field(), "6")
					Expect(valid).To(BeTrue())
					Expect(err).NotTo(HaveOccurred())

					Expect(ts.Unsigned).To(Equal(uint(6)))
				})
			})
			When("the value is invalid", func() {
				It("should return an error", func() {
					ts.Unsigned = 2
					valid, err := SetConfigField(tv, *nodeMap["Unsigned"].Field(), "-3")

					Expect(ts.Unsigned).To(Equal(uint(2)))
					Expect(valid).To(BeFalse())
					Expect(err).To(HaveOccurred())
				})
			})
		})

		When("setting a string slice value", func() {
			When("the value is valid", func() {
				It("should set it", func() {
					valid, err := SetConfigField(tv, *nodeMap["StrSlice"].Field(), "meawannowalkalitabitalleh,meawannofeelalitabitstrongah")
					Expect(valid).To(BeTrue())
					Expect(err).NotTo(HaveOccurred())

					Expect(ts.StrSlice).To(HaveLen(2))
				})
			})
			When("the value is empty", func() {
				It("should be set to an empty slice ", func() {
					valid, err := SetConfigField(tv, *nodeMap["StrSlice"].Field(), "")
					Expect(valid).To(BeTrue())
					Expect(err).NotTo(HaveOccurred())
					Expect(ts.StrSlice).To(HaveLen(0))
				})
			})
		})

		When("setting a string array value", func() {
			When("the value is valid", func() {
				It("should set it", func() {
					valid, err := SetConfigField(tv, *nodeMap["StrArray"].Field(), "meawannowalkalitabitalleh,meawannofeelalitabitstrongah,meawannothinkalitabitsmartah")
					Expect(valid).To(BeTrue())
					Expect(err).NotTo(HaveOccurred())
				})
			})
			When("the value has too many elements", func() {
				It("should return an error", func() {
					valid, err := SetConfigField(tv, *nodeMap["StrArray"].Field(), "one,two,three,four?")
					Expect(valid).To(BeFalse())
					Expect(err).To(HaveOccurred())
				})
			})
			When("the value has too few elements", func() {
				It("should return an error", func() {
					valid, err := SetConfigField(tv, *nodeMap["StrArray"].Field(), "one,two")
					Expect(valid).To(BeFalse())
					Expect(err).To(HaveOccurred())
				})
			})
		})

		When("setting a struct value", func() {
			When("it doesn't implement ConfigProp", func() {
				It("should return an error", func() {
					valid, err := SetConfigField(tv, *nodeMap["Sub"].Field(), "@awol")
					Expect(err).To(HaveOccurred())
					Expect(valid).NotTo(BeTrue())
				})
			})
			When("it implements ConfigProp", func() {
				When("the value is valid", func() {
					It("should set it", func() {
						valid, err := SetConfigField(tv, *nodeMap["SubProp"].Field(), "@awol")
						Expect(err).NotTo(HaveOccurred())
						Expect(valid).To(BeTrue())

						Expect(ts.SubProp.Value).To(Equal("awol"))
					})
				})
				When("the value is invalid", func() {
					It("should return an error", func() {
						valid, err := SetConfigField(tv, *nodeMap["SubProp"].Field(), "missing initial at symbol")
						Expect(err).To(HaveOccurred())
						Expect(valid).NotTo(BeTrue())
					})
				})
			})
		})

		When("setting a struct slice value", func() {
			When("the value is valid", func() {
				It("should set it", func() {
					valid, err := SetConfigField(tv, *nodeMap["SubPropSlice"].Field(), "@alice,@merton")
					Expect(err).NotTo(HaveOccurred())
					Expect(valid).To(BeTrue())

					Expect(ts.SubPropSlice).To(HaveLen(2))
				})
			})
		})

		When("setting a struct pointer slice value", func() {
			When("the value is valid", func() {
				It("should set it", func() {
					valid, err := SetConfigField(tv, *nodeMap["SubPropPtrSlice"].Field(), "@the,@best")
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
				testSetAndFormat(tv, nodeMap["Signed"], "-45", "-45")
			})
			It("should format unsigned integers identical to input", func() {
				testSetAndFormat(tv, nodeMap["Unsigned"], "5", "5")
			})
			It("should format structs identical to input", func() {
				testSetAndFormat(tv, nodeMap["SubProp"], "@whoa", "@whoa")
			})
			It("should format enums identical to input", func() {
				testSetAndFormat(tv, nodeMap["TestEnum"], "Foo", "Foo")
			})
			It("should format string slices identical to input", func() {
				testSetAndFormat(tv, nodeMap["StrSlice"], "one,two,three,four", "[ one, two, three, four ]")
			})
			It("should format string arrays identical to input", func() {
				testSetAndFormat(tv, nodeMap["StrArray"], "one,two,three", "[ one, two, three ]")
			})
			It("should format prop struct slices identical to input", func() {
				testSetAndFormat(tv, nodeMap["SubPropSlice"], "@be,@the,@best", "[ @be, @the, @best ]")
			})
			It("should format prop struct pointer slices identical to input", func() {
				testSetAndFormat(tv, nodeMap["SubPropPtrSlice"], "@diet,@glue", "[ @diet, @glue ]")
			})
			It("should format string maps identical to input", func() {
				testSetAndFormat(tv, nodeMap["StrMap"], "a:1,b:2,c:3", "{ a: 1, b: 2, c: 3 }")
			})

			It("should format int maps identical to input", func() {
				testSetAndFormat(tv, nodeMap["IntMap"], "a:1,b:2,c:3", "{ a: 1, b: 2, c: 3 }")
			})
			It("should format int8 maps identical to input", func() {
				testSetAndFormat(tv, nodeMap["Int8Map"], "a:1,b:2,c:3", "{ a: 1, b: 2, c: 3 }")
			})
			It("should format int16 maps identical to input", func() {
				testSetAndFormat(tv, nodeMap["Int16Map"], "a:1,b:2,c:3", "{ a: 1, b: 2, c: 3 }")
			})
			It("should format int32 maps identical to input", func() {
				testSetAndFormat(tv, nodeMap["Int32Map"], "a:1,b:2,c:3", "{ a: 1, b: 2, c: 3 }")
			})
			It("should format int64 maps identical to input", func() {
				testSetAndFormat(tv, nodeMap["Int64Map"], "a:1,b:2,c:3", "{ a: 1, b: 2, c: 3 }")
			})

			It("should format uint maps identical to input", func() {
				testSetAndFormat(tv, nodeMap["UintMap"], "a:1,b:2,c:3", "{ a: 1, b: 2, c: 3 }")
			})
			It("should format uint8 maps identical to input", func() {
				testSetAndFormat(tv, nodeMap["Uint8Map"], "a:1,b:2,c:3", "{ a: 1, b: 2, c: 3 }")
			})
			It("should format uint16 maps identical to input", func() {
				testSetAndFormat(tv, nodeMap["Uint16Map"], "a:1,b:2,c:3", "{ a: 1, b: 2, c: 3 }")
			})
			It("should format uint32 maps identical to input", func() {
				testSetAndFormat(tv, nodeMap["Uint32Map"], "a:1,b:2,c:3", "{ a: 1, b: 2, c: 3 }")
			})
			It("should format uint64 maps identical to input", func() {
				testSetAndFormat(tv, nodeMap["Uint64Map"], "a:1,b:2,c:3", "{ a: 1, b: 2, c: 3 }")
			})
		})
	})
})

func testSetAndFormat(tv reflect.Value, node Node, value string, prettyFormat string) {
	field := node.Field()
	_, _ = SetConfigField(tv, *field, value)

	// Used for de-/serializing configuration
	formatted, err := GetConfigFieldString(tv, *field)
	Expect(err).NotTo(HaveOccurred())
	Expect(formatted).To(Equal(value))

	node.Update(tv.FieldByName(field.Name))

	// Used for pretty printing output, coloring etc.
	sb := strings.Builder{}
	renderer := ConsoleTreeRenderer{}
	renderer.writeNodeValue(&sb, node)
	Expect(sb.String()).To(Equal(prettyFormat))
}
