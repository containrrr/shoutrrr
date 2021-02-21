package gotify

import (
	"errors"
	"github.com/jarcoal/httpmock"
	"log"
	"net/url"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGotify(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Gotify Suite")
}

var logger *log.Logger

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

	When("TLS is disabled", func() {
		It("should use http schema", func() {
			config := Config{
				Token:      "Aaa.bbb.ccc.ddd",
				Host:       "my.gotify.tld",
				DisableTLS: true,
			}
			url, err := buildURL(&config)
			Expect(err).To(BeNil())
			expectedURL := "http://my.gotify.tld/message?token=Aaa.bbb.ccc.ddd"
			Expect(url).To(Equal(expectedURL))
		})
	})

	When("a custom path is provided", func() {
		It("should add it to the URL", func() {
			config := Config{
				Token: "Aaa.bbb.ccc.ddd",
				Host:  "my.gotify.tld",
				Path:  "/gotify",
			}
			url, err := buildURL(&config)
			Expect(err).To(BeNil())
			expectedURL := "https://my.gotify.tld/gotify/message?token=Aaa.bbb.ccc.ddd"
			Expect(url).To(Equal(expectedURL))
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
	Describe("creating a config", func() {
		When("parsing the configuration URL", func() {
			It("should be identical after de-/serialization (with path)", func() {
				testURL := "gotify://my.gotify.tld/gotify/Aaa.bbb.ccc.ddd?title=Test+title"

				url, err := url.Parse(testURL)
				Expect(err).NotTo(HaveOccurred(), "parsing")

				config := &Config{}
				err = config.SetURL(url)
				Expect(err).NotTo(HaveOccurred(), "verifying")

				outputURL := config.GetURL()
				Expect(outputURL.String()).To(Equal(testURL))

			})
			It("should be identical after de-/serialization (without path)", func() {
				testURL := "gotify://my.gotify.tld/Aaa.bbb.ccc.ddd?disabletls=Yes&priority=1&title=Test+title"

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
			serviceURL, _ := url.Parse("gotify://my.gotify.tld/Aaa.bbb.ccc.ddd")
			err = service.Initialize(serviceURL, logger)
			Expect(err).NotTo(HaveOccurred())

			targetURL := "https://my.gotify.tld/message?token=Aaa.bbb.ccc.ddd"
			httpmock.RegisterResponder("POST", targetURL, httpmock.NewStringResponder(200, ""))

			err = service.Send("Message", nil)
			Expect(err).NotTo(HaveOccurred())
		})
		It("should not panic if an error occurs when sending the payload", func() {
			serviceURL, _ := url.Parse("gotify://my.gotify.tld/Aaa.bbb.ccc.ddd")
			err = service.Initialize(serviceURL, logger)
			Expect(err).NotTo(HaveOccurred())

			targetURL := "https://my.gotify.tld/message?token=Aaa.bbb.ccc.ddd"
			httpmock.RegisterResponder("POST", targetURL, httpmock.NewErrorResponder(errors.New("dummy error")))

			err = service.Send("Message", nil)
			Expect(err).To(HaveOccurred())
		})

	})
})
