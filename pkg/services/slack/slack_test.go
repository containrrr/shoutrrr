package slack_test

import (
	"errors"
	. "github.com/containrrr/shoutrrr/pkg/services/slack"
	"github.com/containrrr/shoutrrr/pkg/util"
	"github.com/jarcoal/httpmock"
	"log"

	"net/url"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gomegaformat "github.com/onsi/gomega/format"
)

func TestSlack(t *testing.T) {
	gomegaformat.CharactersAroundMismatchToInclude = 20
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Slack Suite")
}

var (
	service     *Service
	envSlackURL *url.URL
	logger      *log.Logger
)

var _ = Describe("the slack service", func() {

	BeforeSuite(func() {
		service = &Service{}
		logger = log.New(GinkgoWriter, "Test", log.LstdFlags)
		envSlackURL, _ = url.Parse(os.Getenv("SHOUTRRR_SLACK_URL"))
	})

	When("running integration tests", func() {
		It("should not error out", func() {
			if envSlackURL.String() == "" {
				return
			}

			serviceURL, _ := url.Parse(envSlackURL.String())
			err := service.Initialize(serviceURL, util.TestLogger())
			Expect(err).NotTo(HaveOccurred())

			err = service.Send("This is an integration test message", nil)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	// xoxb:123456789012-1234567890123-4mt0t4l1YL3g1T5L4cK70k3N

	When("given a token with a malformed part", func() {
		It("should return an error if part A is not 9 letters", func() {
			expectErrorMessageGivenURL(ErrorInvalidToken, "slack://lol@12345678/123456789/123456789123456789123456")
		})
		It("should return an error if part B is not 9 letters", func() {
			expectErrorMessageGivenURL(ErrorInvalidToken, "slack://lol@123456789/12345678/123456789123456789123456")
		})
		It("should return an error if part C is not 24 letters", func() {
			expectErrorMessageGivenURL(ErrorInvalidToken, "slack://123456789/123456789/12345678912345678912345")
		})
	})
	When("given a token missing a part", func() {
		It("should return an error if the missing part is A", func() {
			expectErrorMessageGivenURL(ErrorInvalidToken, "slack://lol@/123456789/123456789123456789123456")
		})
		It("should return an error if the missing part is B", func() {
			expectErrorMessageGivenURL(ErrorInvalidToken, "slack://lol@123456789//123456789")
		})
		It("should return an error if the missing part is C", func() {
			expectErrorMessageGivenURL(ErrorInvalidToken, "slack://lol@123456789/123456789/")
		})
	})
	Describe("the slack config", func() {
		When("parsing the configuration URL", func() {
			When("given a config using the legacy format", func() {
				It("should be converted to the new format after de-/serialization", func() {
					oldURL := "slack://testbot@AAAAAAAAA/BBBBBBBBB/123456789123456789123456?color=3f00fe&title=Test+title"
					newURL := "slack://hook:AAAAAAAAA-BBBBBBBBB-123456789123456789123456@webhook?botname=testbot&color=3f00fe&title=Test+title"

					config := &Config{}
					err := config.SetURL(urlMust(oldURL))
					Expect(err).NotTo(HaveOccurred(), "verifying")

					Expect(config.GetURL().String()).To(Equal(newURL))

				})
			})
		})
		It("should be identical after de-/serialization", func() {
			testURL := "slack://hook:AAAAAAAAA-BBBBBBBBB-123456789123456789123456@webhook?botname=testbot&color=3f00fe&title=Test+title"

			config := &Config{}
			err := config.SetURL(urlMust(testURL))
			Expect(err).NotTo(HaveOccurred(), "verifying")

			outputURL := config.GetURL()
			Expect(outputURL.String()).To(Equal(testURL))

		})
		When("generating a config object", func() {
			It("should use the default botname if the argument list contains three strings", func() {
				slackURL, _ := url.Parse("slack://AAAAAAAAA/BBBBBBBBB/123456789123456789123456")
				config, configError := CreateConfigFromURL(slackURL)

				Expect(configError).NotTo(HaveOccurred())
				Expect(config.BotName).To(BeEmpty())
			})
			It("should set the botname if the argument list is three", func() {
				slackURL, _ := url.Parse("slack://testbot@AAAAAAAAA/BBBBBBBBB/123456789123456789123456")
				config, configError := CreateConfigFromURL(slackURL)

				Expect(configError).NotTo(HaveOccurred())
				Expect(config.BotName).To(Equal("testbot"))
			})
			It("should return an error if the argument list is shorter than three", func() {
				slackURL, _ := url.Parse("slack://AAAAAAAA")

				_, configError := CreateConfigFromURL(slackURL)
				Expect(configError).To(HaveOccurred())
			})
		})
		When("getting credentials from token", func() {
			It("should return a valid webhook URL for the given token", func() {
				token := tokenMust("AAAAAAAAA/BBBBBBBBB/123456789123456789123456")
				expected := "https://hooks.slack.com/services/AAAAAAAAA/BBBBBBBBB/123456789123456789123456"
				Expect(token.WebhookURL()).To(Equal(expected))
			})
			It("should return a valid authorization header value for the given token", func() {
				token := tokenMust("xoxb:AAAAAAAAA-BBBBBBBBB-123456789123456789123456")
				expected := "Bearer xoxb-AAAAAAAAA-BBBBBBBBB-123456789123456789123456"
				Expect(token.Authorization()).To(Equal(expected))
			})
		})
	})

	Describe("sending the payload", func() {
		When("sending via webhook URL", func() {
			var err error
			BeforeEach(func() {
				httpmock.Activate()
			})
			AfterEach(func() {
				httpmock.DeactivateAndReset()
			})

			It("should not report an error if the server accepts the payload", func() {
				serviceURL, _ := url.Parse("slack://testbot@AAAAAAAAA/BBBBBBBBB/123456789123456789123456")
				err = service.Initialize(serviceURL, logger)
				Expect(err).NotTo(HaveOccurred())

				targetURL := "https://hooks.slack.com/services/AAAAAAAAA/BBBBBBBBB/123456789123456789123456"
				httpmock.RegisterResponder("POST", targetURL, httpmock.NewStringResponder(200, ""))

				err = service.Send("Message", nil)
				Expect(err).NotTo(HaveOccurred())
			})
			It("should not panic if an error occurs when sending the payload", func() {
				serviceURL, _ := url.Parse("slack://testbot@AAAAAAAAA/BBBBBBBBB/123456789123456789123456")
				err = service.Initialize(serviceURL, logger)
				Expect(err).NotTo(HaveOccurred())

				targetURL := "https://hooks.slack.com/services/AAAAAAAAA/BBBBBBBBB/123456789123456789123456"
				httpmock.RegisterResponder("POST", targetURL, httpmock.NewErrorResponder(errors.New("dummy error")))

				err = service.Send("Message", nil)
				Expect(err).To(HaveOccurred())
			})
		})
		When("sending via bot API", func() {
			var err error
			BeforeEach(func() {
				httpmock.Activate()
			})
			AfterEach(func() {
				httpmock.DeactivateAndReset()
			})

			It("should not report an error if the server accepts the payload", func() {
				serviceURL := urlMust("slack://xoxb:123456789012-1234567890123-4mt0t4l1YL3g1T5L4cK70k3N@C0123456789")
				err = service.Initialize(serviceURL, logger)
				Expect(err).NotTo(HaveOccurred())

				targetURL := "https://slack.com/api/chat.postMessage"
				httpmock.RegisterResponder("POST", targetURL, jsonRespondMust(200, APIResponse{
					Ok: true,
				}))

				err = service.Send("Message", nil)
				Expect(err).NotTo(HaveOccurred())
			})
			It("should not panic if an error occurs when sending the payload", func() {
				serviceURL := urlMust("slack://xoxb:123456789012-1234567890123-4mt0t4l1YL3g1T5L4cK70k3N@C0123456789")
				err = service.Initialize(serviceURL, logger)
				Expect(err).NotTo(HaveOccurred())

				targetURL := "https://slack.com/api/chat.postMessage"
				httpmock.RegisterResponder("POST", targetURL, jsonRespondMust(200, APIResponse{
					Error: "someone turned off the internet",
				}))

				err = service.Send("Message", nil)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})

func jsonRespondMust(code int, response APIResponse) httpmock.Responder {
	responder, err := httpmock.NewJsonResponder(code, response)
	ExpectWithOffset(1, err).NotTo(HaveOccurred(), "invalid test response struct")
	return responder
}

func urlMust(rawURL string) *url.URL {
	parsed, err := url.Parse(rawURL)
	ExpectWithOffset(1, err).NotTo(HaveOccurred(), "invalid test URL string")
	return parsed
}

func tokenMust(rawToken string) *Token {
	token, err := ParseToken(rawToken)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return token
}

func expectErrorMessageGivenURL(expected error, rawURL string) {
	err := service.Initialize(urlMust(rawURL), util.TestLogger())
	ExpectWithOffset(1, err).To(HaveOccurred())
	ExpectWithOffset(1, err).To(Equal(expected))
}
