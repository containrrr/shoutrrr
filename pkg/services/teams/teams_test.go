package teams

import (
	"errors"
	"log"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

const (
	extraIdValue     = "V2ESyij_gAljSoUQHvZoZYzlpAoAXExyOl26dlf1xHEx05"
	scopedWebhookURL = "https://test.webhook.office.com/webhookb2/11111111-4444-4444-8444-cccccccccccc@22222222-4444-4444-8444-cccccccccccc/IncomingWebhook/33333333012222222222333333333344/44444444-4444-4444-8444-cccccccccccc/" + extraIdValue
	scopedDomainHost = "test.webhook.office.com"
	testURLBase      = "teams://11111111-4444-4444-8444-cccccccccccc@22222222-4444-4444-8444-cccccccccccc/33333333012222222222333333333344/44444444-4444-4444-8444-cccccccccccc/" + extraIdValue
	scopedURLBase    = testURLBase + "?host=" + scopedDomainHost
)

var logger = log.New(ginkgo.GinkgoWriter, "Test", log.LstdFlags)

func TestTeams(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Shoutrrr Teams Suite")
}

var _ = ginkgo.Describe("the teams service", func() {
	ginkgo.When("creating the webhook URL", func() {
		ginkgo.It("should match the expected output for custom URLs", func() {
			config := Config{}
			config.setFromWebhookParts([5]string{
				"11111111-4444-4444-8444-cccccccccccc",
				"22222222-4444-4444-8444-cccccccccccc",
				"33333333012222222222333333333344",
				"44444444-4444-4444-8444-cccccccccccc",
				extraIdValue,
			})
			apiURL := BuildWebhookURL(
				scopedDomainHost,
				config.Group,
				config.Tenant,
				config.AltID,
				config.GroupOwner,
				config.ExtraID,
			)
			gomega.Expect(apiURL).To(gomega.Equal(scopedWebhookURL))

			parts, err := parseAndVerifyWebhookURL(apiURL)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(parts).To(gomega.Equal(config.WebhookParts()))
		})
	})

	ginkgo.Describe("creating a config", func() {
		ginkgo.When("parsing the configuration URL", func() {
			ginkgo.It("should be identical after de-/serialization", func() {
				testURL := testURLBase + "?color=aabbcc&host=test.webhook.office.com&title=Test+title"
				url, err := url.Parse(testURL)
				gomega.Expect(err).NotTo(gomega.HaveOccurred(), "parsing")

				config := &Config{}
				err = config.SetURL(url)
				gomega.Expect(err).NotTo(gomega.HaveOccurred(), "verifying")

				outputURL := config.GetURL()
				gomega.Expect(outputURL.String()).To(gomega.Equal(testURL))
			})
		})
	})

	ginkgo.Describe("converting custom URL to service URL", func() {
		ginkgo.When("an invalid custom URL is provided", func() {
			ginkgo.It("should return an error", func() {
				service := Service{}
				testURL := "teams+https://google.com/search?q=what+is+love"
				customURL, err := url.Parse(testURL)
				gomega.Expect(err).NotTo(gomega.HaveOccurred(), "parsing")

				_, err = service.GetConfigURLFromCustom(customURL)
				gomega.Expect(err).To(gomega.HaveOccurred(), "converting")
			})
		})
		ginkgo.When("a valid custom URL is provided", func() {
			ginkgo.It("should set the host field from the custom URL", func() {
				service := Service{}
				testURL := `teams+` + scopedWebhookURL
				customURL, err := url.Parse(testURL)
				gomega.Expect(err).NotTo(gomega.HaveOccurred(), "parsing")

				serviceURL, err := service.GetConfigURLFromCustom(customURL)
				gomega.Expect(err).NotTo(gomega.HaveOccurred(), "converting")
				gomega.Expect(serviceURL.String()).To(gomega.Equal(scopedURLBase))
			})
			ginkgo.It("should preserve the query params in the generated service URL", func() {
				service := Service{}
				testURL := "teams+" + scopedWebhookURL + "?color=f008c1&title=TheTitle"
				customURL, err := url.Parse(testURL)
				gomega.Expect(err).NotTo(gomega.HaveOccurred(), "parsing")

				serviceURL, err := service.GetConfigURLFromCustom(customURL)
				gomega.Expect(err).NotTo(gomega.HaveOccurred(), "converting")
				expectedURL := testURLBase + "?color=f008c1&host=test.webhook.office.com&title=TheTitle"
				gomega.Expect(serviceURL.String()).To(gomega.Equal(expectedURL))
			})
		})
	})

	ginkgo.Describe("sending the payload", func() {
		var err error
		var service Service
		ginkgo.BeforeEach(func() {
			httpmock.Activate()
		})
		ginkgo.AfterEach(func() {
			httpmock.DeactivateAndReset()
		})
		ginkgo.It("should not report an error if the server accepts the payload", func() {
			serviceURL, _ := url.Parse(scopedURLBase)
			err = service.Initialize(serviceURL, logger)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			httpmock.RegisterResponder(
				"POST",
				scopedWebhookURL,
				httpmock.NewStringResponder(200, ""),
			)
			err = service.Send("Message", nil)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
		ginkgo.It("should not panic if an error occurs when sending the payload", func() {
			serviceURL, _ := url.Parse(testURLBase + "?host=test.webhook.office.com")
			err = service.Initialize(serviceURL, logger)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			httpmock.RegisterResponder(
				"POST",
				scopedWebhookURL,
				httpmock.NewErrorResponder(errors.New("dummy error")),
			)
			err = service.Send("Message", nil)
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
	})

	ginkgo.It("should return the correct service ID", func() {
		service := &Service{}
		gomega.Expect(service.GetID()).To(gomega.Equal("teams"))
	})

	// Config tests
	ginkgo.Describe("the teams config", func() {
		ginkgo.Describe("setURL", func() {
			ginkgo.It("should set all fields correctly from URL", func() {
				config := &Config{}
				urlStr := testURLBase + "?title=Test&color=red&host=test.webhook.office.com"
				parsedURL, err := url.Parse(urlStr)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				err = config.SetURL(parsedURL)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				gomega.Expect(config.Group).To(gomega.Equal("11111111-4444-4444-8444-cccccccccccc"))
				gomega.Expect(config.Tenant).
					To(gomega.Equal("22222222-4444-4444-8444-cccccccccccc"))
				gomega.Expect(config.AltID).To(gomega.Equal("33333333012222222222333333333344"))
				gomega.Expect(config.GroupOwner).
					To(gomega.Equal("44444444-4444-4444-8444-cccccccccccc"))
				gomega.Expect(config.ExtraID).To(gomega.Equal(extraIdValue))
				gomega.Expect(config.Title).To(gomega.Equal("Test"))
				gomega.Expect(config.Color).To(gomega.Equal("red"))
				gomega.Expect(config.Host).To(gomega.Equal("test.webhook.office.com"))
			})

			ginkgo.It("should reject URLs missing the extraID", func() {
				config := &Config{}
				urlStr := "teams://11111111-4444-4444-8444-cccccccccccc@22222222-4444-4444-8444-cccccccccccc/33333333012222222222333333333344/44444444-4444-4444-8444-cccccccccccc?host=test.webhook.office.com"
				parsedURL, err := url.Parse(urlStr)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				err = config.SetURL(parsedURL)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})

			ginkgo.It("should require the host parameter", func() {
				config := &Config{}
				urlStr := testURLBase
				parsedURL, err := url.Parse(urlStr)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				err = config.SetURL(parsedURL)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Describe("getURL", func() {
			ginkgo.It("should generate correct URL with all parameters", func() {
				config := &Config{
					Group:      "11111111-4444-4444-8444-cccccccccccc",
					Tenant:     "22222222-4444-4444-8444-cccccccccccc",
					AltID:      "33333333012222222222333333333344",
					GroupOwner: "44444444-4444-4444-8444-cccccccccccc",
					ExtraID:    extraIdValue,
					Title:      "Test",
					Color:      "red",
					Host:       "test.webhook.office.com",
				}

				urlObj := config.GetURL()
				urlStr := urlObj.String()
				expectedURL := testURLBase + "?color=red&host=test.webhook.office.com&title=Test"
				gomega.Expect(urlStr).To(gomega.Equal(expectedURL))
			})
		})

		ginkgo.Describe("verifyWebhookParts", func() {
			ginkgo.It("should validate correct webhook parts", func() {
				parts := [5]string{
					"11111111-4444-4444-8444-cccccccccccc",
					"22222222-4444-4444-8444-cccccccccccc",
					"33333333012222222222333333333344",
					"44444444-4444-4444-8444-cccccccccccc",
					extraIdValue,
				}
				err := verifyWebhookParts(parts)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})

			ginkgo.It("should reject invalid group ID", func() {
				parts := [5]string{
					"invalid-id",
					"22222222-4444-4444-8444-cccccccccccc",
					"33333333012222222222333333333344",
					"44444444-4444-4444-8444-cccccccccccc",
					extraIdValue,
				}
				err := verifyWebhookParts(parts)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Describe("parseAndVerifyWebhookURL", func() {
			ginkgo.It("should correctly parse valid webhook URL", func() {
				webhookURL := scopedWebhookURL
				parts, err := parseAndVerifyWebhookURL(webhookURL)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(parts).To(gomega.Equal([5]string{
					"11111111-4444-4444-8444-cccccccccccc",
					"22222222-4444-4444-8444-cccccccccccc",
					"33333333012222222222333333333344",
					"44444444-4444-4444-8444-cccccccccccc",
					extraIdValue,
				}))
			})

			ginkgo.It("should reject invalid webhook URL", func() {
				webhookURL := "https://teams.microsoft.com/invalid/webhook/url"
				_, err := parseAndVerifyWebhookURL(webhookURL)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})
})
