package plugins_test

import (
	. "github.com/containrrr/shoutrrr/pkg/plugins"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestPlugins(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Plugins Suite")
}

var _ = Describe("the plugin package", func() {
	When("extract arguments is given a url", func() {
		It("should return the arguments", func() {
			url := "slack://aaaa"
			arguments, err := ExtractArguments(url)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(arguments)).To(Equal(1))
		})
		It("should return an error if no arguments could be found", func() {
			url := "slack://"
			arguments, err := ExtractArguments(url)
			Expect(err).To(HaveOccurred())
			Expect(len(arguments)).To(Equal(0))
		})
		It("should split the arguments by /", func() {
			url := "slack://aaaa/bbb/ccc"
			arguments, err := ExtractArguments(url)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(arguments)).To(Equal(3))
			Expect(arguments[0]).To(Equal("aaaa"))
			Expect(arguments[1]).To(Equal("bbb"))
			Expect(arguments[2]).To(Equal("ccc"))
		})
	})
})