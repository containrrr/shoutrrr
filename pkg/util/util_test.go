package util_test

import (
	"testing"
	. "github.com/containrrr/shoutrrr/pkg/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
			min := Min(a, b)

			Expect(min).To(Equal(a))
		})
	})

	When("calling function Max", func() {
		It("should return the largest of two integers", func() {
			max := Max(a, b)

			Expect(max).To(Equal(b))
		})
	})
})