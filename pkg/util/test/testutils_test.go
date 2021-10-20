package test_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/containrrr/shoutrrr/pkg/util/test"
)

func TestUtil(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Util Test Suite")
}

var _ = Describe("the util package", func() {
	When("calling function TestLogger", func() {
		It("should not return nil", func() {
			Expect(TestLogger()).NotTo(Equal(nil))
		})
		It("should have the prefix \"Test\"", func() {
			Expect(TestLogger().Prefix()).To(Equal("Test"))
		})
	})
})
