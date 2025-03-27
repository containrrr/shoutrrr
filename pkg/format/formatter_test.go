package format

import (
	"reflect"
	"strings"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// logger *log.Logger.
var ts *testStruct

var _ = ginkgo.Describe("SetConfigField", func() {
	var tv reflect.Value
	testConfig := testStruct{}
	tt := reflect.TypeOf(testConfig)
	rootNode := getRootNode(&testConfig)
	nodeMap := make(map[string]Node, len(rootNode.Items))
	for _, item := range rootNode.Items {
		field := item.Field()
		nodeMap[field.Name] = item
	}
	ginkgo.When("updating a struct", func() {
		ginkgo.BeforeEach(func() {
			tsPtr := reflect.New(tt)
			tv = tsPtr.Elem()
			ts = tsPtr.Interface().(*testStruct)
		})
		ginkgo.When("setting an integer value", func() {
			ginkgo.When("the value is valid", func() {
				ginkgo.It("should set it", func() {
					valid, err := SetConfigField(tv, *nodeMap["Signed"].Field(), "3")
					gomega.Expect(valid).To(gomega.BeTrue())
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(ts.Signed).To(gomega.Equal(3))
				})
			})
			ginkgo.When("the value is invalid", func() {
				ginkgo.It("should return an error", func() {
					ts.Signed = 2
					valid, err := SetConfigField(tv, *nodeMap["Signed"].Field(), "z7")
					gomega.Expect(valid).To(gomega.BeFalse())
					gomega.Expect(err).To(gomega.HaveOccurred())
					gomega.Expect(ts.Signed).To(gomega.Equal(2))
				})
			})
		})
		ginkgo.When("setting an unsigned integer value", func() {
			ginkgo.When("the value is valid", func() {
				ginkgo.It("should set it", func() {
					valid, err := SetConfigField(tv, *nodeMap["Unsigned"].Field(), "6")
					gomega.Expect(valid).To(gomega.BeTrue())
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(ts.Unsigned).To(gomega.Equal(uint(6)))
				})
			})
			ginkgo.When("the value is invalid", func() {
				ginkgo.It("should return an error", func() {
					ts.Unsigned = 2
					valid, err := SetConfigField(tv, *nodeMap["Unsigned"].Field(), "-3")
					gomega.Expect(ts.Unsigned).To(gomega.Equal(uint(2)))
					gomega.Expect(valid).To(gomega.BeFalse())
					gomega.Expect(err).To(gomega.HaveOccurred())
				})
			})
		})
		ginkgo.When("setting a string slice value", func() {
			ginkgo.When("the value is valid", func() {
				ginkgo.It("should set it", func() {
					valid, err := SetConfigField(tv, *nodeMap["StrSlice"].Field(), "meawannowalkalitabitalleh,meawannofeelalitabitstrongah")
					gomega.Expect(valid).To(gomega.BeTrue())
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(ts.StrSlice).To(gomega.HaveLen(2))
				})
			})
		})
		ginkgo.When("setting a string array value", func() {
			ginkgo.When("the value is valid", func() {
				ginkgo.It("should set it", func() {
					valid, err := SetConfigField(tv, *nodeMap["StrArray"].Field(), "meawannowalkalitabitalleh,meawannofeelalitabitstrongah,meawannothinkalitabitsmartah")
					gomega.Expect(valid).To(gomega.BeTrue())
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
				})
			})
			ginkgo.When("the value has too many elements", func() {
				ginkgo.It("should return an error", func() {
					valid, err := SetConfigField(tv, *nodeMap["StrArray"].Field(), "one,two,three,four?")
					gomega.Expect(valid).To(gomega.BeFalse())
					gomega.Expect(err).To(gomega.HaveOccurred())
				})
			})
			ginkgo.When("the value has too few elements", func() {
				ginkgo.It("should return an error", func() {
					valid, err := SetConfigField(tv, *nodeMap["StrArray"].Field(), "one,two")
					gomega.Expect(valid).To(gomega.BeFalse())
					gomega.Expect(err).To(gomega.HaveOccurred())
				})
			})
		})
		ginkgo.When("setting a struct value", func() {
			ginkgo.When("it doesn't implement ConfigProp", func() {
				ginkgo.It("should return an error", func() {
					valid, err := SetConfigField(tv, *nodeMap["Sub"].Field(), "@awol")
					gomega.Expect(err).To(gomega.HaveOccurred())
					gomega.Expect(valid).NotTo(gomega.BeTrue())
				})
			})
			ginkgo.When("it implements ConfigProp", func() {
				ginkgo.When("the value is valid", func() {
					ginkgo.It("should set it", func() {
						valid, err := SetConfigField(tv, *nodeMap["SubProp"].Field(), "@awol")
						gomega.Expect(err).NotTo(gomega.HaveOccurred())
						gomega.Expect(valid).To(gomega.BeTrue())
						gomega.Expect(ts.SubProp.Value).To(gomega.Equal("awol"))
					})
				})
				ginkgo.When("the value is invalid", func() {
					ginkgo.It("should return an error", func() {
						valid, err := SetConfigField(tv, *nodeMap["SubProp"].Field(), "missing initial at symbol")
						gomega.Expect(err).To(gomega.HaveOccurred())
						gomega.Expect(valid).NotTo(gomega.BeTrue())
					})
				})
			})
		})
		ginkgo.When("setting a struct slice value", func() {
			ginkgo.When("the value is valid", func() {
				ginkgo.It("should set it", func() {
					valid, err := SetConfigField(tv, *nodeMap["SubPropSlice"].Field(), "@alice,@merton")
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(valid).To(gomega.BeTrue())
					gomega.Expect(ts.SubPropSlice).To(gomega.HaveLen(2))
				})
			})
		})
		ginkgo.When("setting a struct pointer slice value", func() {
			ginkgo.When("the value is valid", func() {
				ginkgo.It("should set it", func() {
					valid, err := SetConfigField(tv, *nodeMap["SubPropPtrSlice"].Field(), "@the,@best")
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(valid).To(gomega.BeTrue())
					gomega.Expect(ts.SubPropPtrSlice).To(gomega.HaveLen(2))
				})
			})
		})
	})
	ginkgo.When("formatting stuct values", func() {
		ginkgo.BeforeEach(func() {
			tsPtr := reflect.New(tt)
			tv = tsPtr.Elem()
		})
		ginkgo.When("setting and formatting", func() {
			ginkgo.It("should format signed integers identical to input", func() {
				testSetAndFormat(tv, nodeMap["Signed"], "-45", "-45")
			})
			ginkgo.It("should format unsigned integers identical to input", func() {
				testSetAndFormat(tv, nodeMap["Unsigned"], "5", "5")
			})
			ginkgo.It("should format structs identical to input", func() {
				testSetAndFormat(tv, nodeMap["SubProp"], "@whoa", "@whoa")
			})
			ginkgo.It("should format enums identical to input", func() {
				testSetAndFormat(tv, nodeMap["TestEnum"], "Foo", "Foo")
			})
			ginkgo.It("should format string slices identical to input", func() {
				testSetAndFormat(tv, nodeMap["StrSlice"], "one,two,three,four", "[ one, two, three, four ]")
			})
			ginkgo.It("should format string arrays identical to input", func() {
				testSetAndFormat(tv, nodeMap["StrArray"], "one,two,three", "[ one, two, three ]")
			})
			ginkgo.It("should format prop struct slices identical to input", func() {
				testSetAndFormat(tv, nodeMap["SubPropSlice"], "@be,@the,@best", "[ @be, @the, @best ]")
			})
			ginkgo.It("should format prop struct pointer slices identical to input", func() {
				testSetAndFormat(tv, nodeMap["SubPropPtrSlice"], "@diet,@glue", "[ @diet, @glue ]")
			})
			ginkgo.It("should format string maps identical to input", func() {
				testSetAndFormat(tv, nodeMap["StrMap"], "a:1,b:2,c:3", "{ a: 1, b: 2, c: 3 }")
			})
			ginkgo.It("should format int maps identical to input", func() {
				testSetAndFormat(tv, nodeMap["IntMap"], "a:1,b:2,c:3", "{ a: 1, b: 2, c: 3 }")
			})
			ginkgo.It("should format int8 maps identical to input", func() {
				testSetAndFormat(tv, nodeMap["Int8Map"], "a:1,b:2,c:3", "{ a: 1, b: 2, c: 3 }")
			})
			ginkgo.It("should format int16 maps identical to input", func() {
				testSetAndFormat(tv, nodeMap["Int16Map"], "a:1,b:2,c:3", "{ a: 1, b: 2, c: 3 }")
			})
			ginkgo.It("should format int32 maps identical to input", func() {
				testSetAndFormat(tv, nodeMap["Int32Map"], "a:1,b:2,c:3", "{ a: 1, b: 2, c: 3 }")
			})
			ginkgo.It("should format int64 maps identical to input", func() {
				testSetAndFormat(tv, nodeMap["Int64Map"], "a:1,b:2,c:3", "{ a: 1, b: 2, c: 3 }")
			})
			ginkgo.It("should format uint maps identical to input", func() {
				testSetAndFormat(tv, nodeMap["UintMap"], "a:1,b:2,c:3", "{ a: 1, b: 2, c: 3 }")
			})
			ginkgo.It("should format uint8 maps identical to input", func() {
				testSetAndFormat(tv, nodeMap["Uint8Map"], "a:1,b:2,c:3", "{ a: 1, b: 2, c: 3 }")
			})
			ginkgo.It("should format uint16 maps identical to input", func() {
				testSetAndFormat(tv, nodeMap["Uint16Map"], "a:1,b:2,c:3", "{ a: 1, b: 2, c: 3 }")
			})
			ginkgo.It("should format uint32 maps identical to input", func() {
				testSetAndFormat(tv, nodeMap["Uint32Map"], "a:1,b:2,c:3", "{ a: 1, b: 2, c: 3 }")
			})
			ginkgo.It("should format uint64 maps identical to input", func() {
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
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	gomega.Expect(formatted).To(gomega.Equal(value))
	node.Update(tv.FieldByName(field.Name))

	// Used for pretty printing output, coloring etc.
	sb := strings.Builder{}
	renderer := ConsoleTreeRenderer{}
	renderer.writeNodeValue(&sb, node)
	gomega.Expect(sb.String()).To(gomega.Equal(prettyFormat))
}
