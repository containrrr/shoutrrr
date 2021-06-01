package googlechat

import (
	"net/url"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGooglechat(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Google Chat Suite")
}

var _ = Describe("the Googlechat Chat plugin URL building", func() {
	It("should build a valid Google Chat Incoming Webhook URL", func() {
		configURL, _ := url.Parse("googlechat://chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=baz")

		config := Config{}
		config.SetURL(configURL)

		expectedURL := "https://chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=baz"
		Expect(getAPIURL(&config).String()).To(Equal(expectedURL))
	})
})
