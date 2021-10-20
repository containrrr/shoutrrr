package googlechat

import (
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGooglechat(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Google Chat Suite")
}

var _ = Describe("Google Chat Service", func() {
	It("should build a valid Google Chat Incoming Webhook URL", func() {
		configURL, _ := url.Parse("googlechat://chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=baz")

		config := Config{}
		Expect(config.SetURL(configURL)).To(Succeed())

		expectedURL := "https://chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=baz"
		Expect(getAPIURL(&config).String()).To(Equal(expectedURL))
	})
	When("parsing the configuration URL", func() {
		It("should be identical after de-/serialization", func() {
			testURL := "googlechat://chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=baz"

			url, err := url.Parse(testURL)
			Expect(err).NotTo(HaveOccurred(), "parsing")

			config := &Config{}
			err = config.SetURL(url)
			Expect(err).NotTo(HaveOccurred(), "verifying")

			outputURL := config.GetURL()

			Expect(outputURL.String()).To(Equal(testURL))

		})
	})

	Describe("sending the payload", func() {
		var err error
		BeforeEach(func() {
			httpmock.Activate()
		})
		AfterEach(func() {
			httpmock.DeactivateAndReset()
		})
		It("should not report an error if the server accepts the payload", func() {
			config := Config{
				Host:  "chat.googleapis.com",
				Path:  "v1/spaces/FOO/messages",
				Key:   "bar",
				Token: "baz",
			}
			serviceURL := config.GetURL()
			service := Service{}
			err = service.Initialize(serviceURL, nil)
			Expect(err).NotTo(HaveOccurred())

			httpmock.RegisterResponder("POST", "https://chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=baz", httpmock.NewStringResponder(200, `{}`))

			err = service.Send("Message", nil)
			Expect(err).NotTo(HaveOccurred())
		})

	})
})
