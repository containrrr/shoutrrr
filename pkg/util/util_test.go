package util_test

import (
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
})