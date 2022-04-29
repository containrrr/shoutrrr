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

	Describe("Must helpers", func() {
		Describe("URLMust", func() {
			It("should panic when an invalid URL is passed", func() {
				failures := InterceptGomegaFailures(func() { URLMust(":") })
				Expect(failures).To(HaveLen(1))
			})
		})

		Describe("JSONRespondMust", func() {
			It("should panic when an invalid struct is passed", func() {
				notAValidJSONSource := func() {}
				failures := InterceptGomegaFailures(func() { JSONRespondMust(200, notAValidJSONSource) })
				Expect(failures).To(HaveLen(1))
			})
		})
	})
})
