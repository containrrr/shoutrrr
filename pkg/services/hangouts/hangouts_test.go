package hangouts

import (
	"net/url"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHangouts(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Hangouts Chat Suite")
}

var _ = Describe("the Hangouts Chat plugin URL building", func() {
	It("should build a valid Hangouts Chat Incoming Webhook URL", func() {
		configURL, _ := url.Parse("hangouts://chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=baz")

		config := Config{}
		config.SetURL(configURL)

		expectedURL := "https://chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=baz"
		Expect(config.URL.String()).To(Equal(expectedURL))
	})
})
