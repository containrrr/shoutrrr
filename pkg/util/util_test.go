package util_test

import (
	"reflect"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/containrrr/shoutrrr/pkg/util"
)

func TestUtil(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Util Suite")
}

const a = 10
const b = 20

var _ = Describe("the util package", func() {
	When("calling function Min", func() {
		It("should return the smallest of two integers", func() {
			Expect(Min(a, b)).To(Equal(a))
			Expect(Min(b, a)).To(Equal(a))
		})
	})

	When("calling function Max", func() {
		It("should return the largest of two integers", func() {
			Expect(Max(a, b)).To(Equal(b))
			Expect(Max(b, a)).To(Equal(b))
		})
	})

	When("calling function TestLogger", func() {
		It("should not return nil", func() {
			Expect(TestLogger()).NotTo(Equal(nil))
		})
		It("should have the prefix \"Test\"", func() {
			Expect(TestLogger().Prefix()).To(Equal("Test"))
		})
	})
	When("checking if a supplied kind is of the signed integer kind", func() {
		It("should be true if the kind is Int", func() {
			Expect(IsSignedInt(reflect.Int)).To(BeTrue())
		})
		It("should be false if the kind is String", func() {
			Expect(IsSignedInt(reflect.String)).To(BeFalse())
		})
	})

	When("checking if a supplied kind is of the unsigned integer kind", func() {
		It("should be true if the kind is Uint", func() {
			Expect(IsUnsignedInt(reflect.Uint)).To(BeTrue())
		})
		It("should be false if the kind is Int", func() {
			Expect(IsUnsignedInt(reflect.Int)).To(BeFalse())
		})
	})

	When("checking if a supplied kind is of the collection kind", func() {
		It("should be true if the kind is slice", func() {
			Expect(IsCollection(reflect.Slice)).To(BeTrue())
		})
		It("should be false if the kind is map", func() {
			Expect(IsCollection(reflect.Map)).To(BeFalse())
		})
	})

	When("checking if a supplied kind is numeric", func() {
		It("should be true if supplied a constant integer", func() {
			Expect(IsNumeric(reflect.TypeOf(5).Kind())).To(BeTrue())
		})
		It("should be true if supplied a constant float", func() {
			Expect(IsNumeric(reflect.TypeOf(2.5).Kind())).To(BeTrue())
		})
		It("should be false if supplied a constant string", func() {
			Expect(IsNumeric(reflect.TypeOf("3").Kind())).To(BeFalse())
		})
	})
})
