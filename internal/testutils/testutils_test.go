package testutils_test

import (
	"testing"

	. "github.com/containrrr/shoutrrr/internal/testutils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTestUtils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr TestUtils Suite")
}

var _ = Describe("the testutils package", func() {
	When("calling function TestLogger", func() {
		It("should not return nil", func() {
			Expect(TestLogger()).NotTo(Equal(nil))
		})
		It("should have the prefix \"Test\"", func() {
			Expect(TestLogger().Prefix()).To(Equal("Test"))
		})
	})
})
