package teams

import (
	"errors"
	"log"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	legacyWebhookURL = "https://outlook.office.com/webhook/11111111-4444-4444-8444-cccccccccccc@22222222-4444-4444-8444-cccccccccccc/IncomingWebhook/33333333012222222222333333333344/44444444-4444-4444-8444-cccccccccccc"
	scopedWebhookURL = "https://test.webhook.office.com/webhookb2/11111111-4444-4444-8444-cccccccccccc@22222222-4444-4444-8444-cccccccccccc/IncomingWebhook/33333333012222222222333333333344/44444444-4444-4444-8444-cccccccccccc"
	scopedDomainHost = "test.webhook.office.com"
	testURLBase      = "teams://11111111-4444-4444-8444-cccccccccccc@22222222-4444-4444-8444-cccccccccccc/33333333012222222222333333333344/44444444-4444-4444-8444-cccccccccccc"
	scopedURLBase    = testURLBase + `?host=` + scopedDomainHost
)

var logger = log.New(GinkgoWriter, "Test", log.LstdFlags)

func TestTeams(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Teams Service Suite")
}

var _ = Describe("the teams service", func() {
	When("creating the webhook URL", func() {
		It("should match the expected output for legacy URLs", func() {
			config := Config{}
			config.setFromWebhookParts([4]string{
				"11111111-4444-4444-8444-cccccccccccc",
				"22222222-4444-4444-8444-cccccccccccc",
				"33333333012222222222333333333344",
				"44444444-4444-4444-8444-cccccccccccc",
			})
			apiURL := buildWebhookURL(LegacyHost, config.Group, config.Tenant, config.AltID, config.GroupOwner)
			Expect(apiURL).To(Equal(legacyWebhookURL))

			parts, err := parseAndVerifyWebhookURL(apiURL)
			Expect(err).ToNot(HaveOccurred())
			Expect(parts).To(Equal(config.webhookParts()))
		})
		It("should match the expected output for custom URLs", func() {
			config := Config{}
			config.setFromWebhookParts([4]string{
				"11111111-4444-4444-8444-cccccccccccc",
				"22222222-4444-4444-8444-cccccccccccc",
				"33333333012222222222333333333344",
				"44444444-4444-4444-8444-cccccccccccc",
			})
			apiURL := buildWebhookURL(scopedDomainHost, config.Group, config.Tenant, config.AltID, config.GroupOwner)
			Expect(apiURL).To(Equal(scopedWebhookURL))

			parts, err := parseAndVerifyWebhookURL(apiURL)
			Expect(err).ToNot(HaveOccurred())
			Expect(parts).To(Equal(config.webhookParts()))
		})
	})

	Describe("creating a config", func() {
		When("parsing the configuration URL", func() {
			It("should be identical after de-/serialization", func() {
				testURL := testURLBase + "?color=aabbcc&host=test.outlook.office.com&title=Test+title"

				url, err := url.Parse(testURL)
				Expect(err).NotTo(HaveOccurred(), "parsing")

				config := &Config{Host: LegacyHost}
				err = config.SetURL(url)
				Expect(err).NotTo(HaveOccurred(), "verifying")

				outputURL := config.GetURL()
				Expect(outputURL.String()).To(Equal(testURL))
			})
		})
	})

	Describe("converting custom URL to service URL", func() {
		When("an invalid custom URL is provided", func() {
			It("should return an error", func() {
				service := Service{}
				testURL := "teams+https://google.com/search?q=what+is+love"

				customURL, err := url.Parse(testURL)
				Expect(err).NotTo(HaveOccurred(), "parsing")

				_, err = service.GetConfigURLFromCustom(customURL)
				Expect(err).To(HaveOccurred(), "converting")
			})
		})
		When("a valid custom URL is provided", func() {
			It("should set the host field from the custom URL", func() {
				service := Service{}
				testURL := `teams+` + scopedWebhookURL

				customURL, err := url.Parse(testURL)
				Expect(err).NotTo(HaveOccurred(), "parsing")

				serviceURL, err := service.GetConfigURLFromCustom(customURL)
				Expect(err).NotTo(HaveOccurred(), "converting")

				Expect(serviceURL.String()).To(Equal(scopedURLBase))
			})
			It("should preserve the query params in the generated service URL", func() {
				service := Service{}
				testURL := "teams+" + legacyWebhookURL + "?color=f008c1&title=TheTitle"

				customURL, err := url.Parse(testURL)
				Expect(err).NotTo(HaveOccurred(), "parsing")

				serviceURL, err := service.GetConfigURLFromCustom(customURL)
				Expect(err).NotTo(HaveOccurred(), "converting")

				Expect(serviceURL.String()).To(Equal(testURLBase + "?color=f008c1&title=TheTitle"))
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
			serviceURL, _ := url.Parse(scopedURLBase)
			err = service.Initialize(serviceURL, logger)
			Expect(err).NotTo(HaveOccurred())

			httpmock.RegisterResponder("POST", scopedWebhookURL, httpmock.NewStringResponder(200, ""))

			err = service.Send("Message", nil)
			Expect(err).NotTo(HaveOccurred())
		})
		It("should not panic if an error occurs when sending the payload", func() {
			serviceURL, _ := url.Parse(testURLBase)
			err = service.Initialize(serviceURL, logger)
			Expect(err).NotTo(HaveOccurred())

			httpmock.RegisterResponder("POST", legacyWebhookURL, httpmock.NewErrorResponder(errors.New("dummy error")))

			err = service.Send("Message", nil)
			Expect(err).To(HaveOccurred())
		})

	})

})
