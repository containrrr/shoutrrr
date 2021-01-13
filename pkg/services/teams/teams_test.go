package teams

import (
	"errors"
	"github.com/jarcoal/httpmock"
	"log"
	"net/url"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	testWebhookURL = "https://outlook.office.com/webhook/11111111-4444-4444-8444-cccccccccccc@22222222-4444-4444-8444-cccccccccccc/IncomingWebhook/33333333012222222222333333333344/44444444-4444-4444-8444-cccccccccccc"
	testURLBase    = "teams://11111111-4444-4444-8444-cccccccccccc@22222222-4444-4444-8444-cccccccccccc/33333333012222222222333333333344/44444444-4444-4444-8444-cccccccccccc"
)

var logger = log.New(GinkgoWriter, "Test", log.LstdFlags)

func TestTeams(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Teams Suite")
}

var _ = Describe("the teams plugin", func() {
	When("creating the webhook URL", func() {
		It("should match the expected output", func() {
			config := Config{
				WebhookParts: [4]string{
					"11111111-4444-4444-8444-cccccccccccc",
					"22222222-4444-4444-8444-cccccccccccc",
					"33333333012222222222333333333344",
					"44444444-4444-4444-8444-cccccccccccc",
				},
			}
			apiURL := buildWebhookURL(config.WebhookParts)
			Expect(apiURL).To(Equal(testWebhookURL))

			parts, err := parseAndVerifyWebhookURL(apiURL)
			Expect(err).ToNot(HaveOccurred())
			Expect(parts).To(Equal(config.WebhookParts))
		})
	})

	Describe("creating a config", func() {
		When("parsing the configuration URL", func() {
			It("should be identical after de-/serialization", func() {
				testURL := testURLBase + "?color=aabbcc&title=Test title"

				url, err := url.Parse(testURL)
				Expect(err).NotTo(HaveOccurred(), "parsing")

				config := &Config{}
				err = config.SetURL(url)
				Expect(err).NotTo(HaveOccurred(), "verifying")

				outputURL := config.GetURL()
				Expect(outputURL.String()).To(Equal(testURL))

			})
		})
	})

	Describe("sending the payload", func() {
		var err error
		var service Service
		BeforeEach(func() {
			httpmock.Activate()
		})
		AfterEach(func() {
			httpmock.DeactivateAndReset()
		})
		It("should not report an error if the server accepts the payload", func() {
			serviceURL, _ := url.Parse(testURLBase)
			err = service.Initialize(serviceURL, logger)
			Expect(err).NotTo(HaveOccurred())

			httpmock.RegisterResponder("POST", testWebhookURL, httpmock.NewStringResponder(200, ""))

			err = service.Send("Message", nil)
			Expect(err).NotTo(HaveOccurred())
		})
		It("should not panic if an error occurs when sending the payload", func() {
			serviceURL, _ := url.Parse(testURLBase)
			err = service.Initialize(serviceURL, logger)
			Expect(err).NotTo(HaveOccurred())

			httpmock.RegisterResponder("POST", testWebhookURL, httpmock.NewErrorResponder(errors.New("dummy error")))

			err = service.Send("Message", nil)
			Expect(err).To(HaveOccurred())
		})

	})

})
