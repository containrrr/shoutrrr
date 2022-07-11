package format

import (
	"testing"

	"github.com/fatih/color"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFormat(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Format Suite")
}

var _ = Describe("the format package", func() {
	BeforeSuite(func() {
		// logger = log.New(GinkgoWriter, "Test", log.LstdFlags)

		// Disable color output for tests to have them match the string format rather than the colors
		color.NoColor = true
	})

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
			testEnum := CreateEnumFormatter([]string{"None", "Foo", "Bar"})
			Expect(testEnum.Names()).To(ConsistOf("None", "Foo", "Bar"))
		})
	})
})
