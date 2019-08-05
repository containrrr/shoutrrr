package gotify

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGotify(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Gotify Suite")
}

var _ = Describe("the Gotify plugin URL building and token validation functions", func() {
	It("should build a valid gotify URL", func() {
		config := Config{
			Token: "Aaa.bbb.ccc.ddd",
			Host:  "my.gotify.tld",
		}
		url, err := buildURL(&config)
		Expect(err).To(BeNil())
		expectedURL := "https://my.gotify.tld/message?token=Aaa.bbb.ccc.ddd"
		Expect(url).To(Equal(expectedURL))
	})
	When("provided empty params", func() {
		It("should return 0", func() {
			params := make(map[string]string)
			priority := getPriority(params)
			Expect(priority).To(Equal(0))
		})
	})
	When("provided invalid params", func() {
		It("should return 0", func() {
			params := make(map[string]string)
			params["priority"] = "not an integer"
			priority := getPriority(params)
			Expect(priority).To(Equal(0))
		})
	})
	When("provided 42", func() {
		It("should return 42", func() {
			params := make(map[string]string)
			params["priority"] = "42"
			priority := getPriority(params)
			Expect(priority).To(Equal(42))
		})
	})
	When("provided a valid token", func() {
		It("should return true", func() {
			token := "Ahwbsdyhwwgarxd"
			Expect(isTokenValid(token)).To(BeTrue())
		})
	})
	When("provided a token with an invalid prefix", func() {
		It("should return false", func() {
			token := "Chwbsdyhwwgarxd"
			Expect(isTokenValid(token)).To(BeFalse())
		})
	})
	When("provided a token with an invalid length", func() {
		It("should return false", func() {
			token := "Chwbsdyhwwga"
			Expect(isTokenValid(token)).To(BeFalse())
		})
	})
})
