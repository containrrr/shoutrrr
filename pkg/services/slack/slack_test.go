package slack_test

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/nicholas-fedor/shoutrrr/internal/testutils"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/slack"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
)

const (
	TestWebhookURL = "https://hooks.slack.com/services/AAAAAAAAA/BBBBBBBBB/123456789123456789123456"
)

func TestSlack(t *testing.T) {
	format.CharactersAroundMismatchToInclude = 20

	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Shoutrrr Slack Suite")
}

var (
	service     *slack.Service
	envSlackURL *url.URL
	logger      *log.Logger
	_           = ginkgo.BeforeSuite(func() {
		service = &slack.Service{}
		logger = log.New(ginkgo.GinkgoWriter, "Test", log.LstdFlags)
		envSlackURL, _ = url.Parse(os.Getenv("SHOUTRRR_SLACK_URL"))
	})
)

var _ = ginkgo.Describe("the slack service", func() {
	ginkgo.When("running integration tests", func() {
		ginkgo.It("should not error out", func() {
			if envSlackURL.String() == "" {
				return
			}

			serviceURL, _ := url.Parse(envSlackURL.String())
			err := service.Initialize(serviceURL, testutils.TestLogger())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			err = service.Send("This is an integration test message", nil)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
		ginkgo.It("returns the correct service identifier", func() {
			gomega.Expect(service.GetID()).To(gomega.Equal("slack"))
		})
	})

	// xoxb:123456789012-1234567890123-4mt0t4l1YL3g1T5L4cK70k3N

	ginkgo.When("given a token with a malformed part", func() {
		ginkgo.It("should return an error if part A is not 9 letters", func() {
			expectErrorMessageGivenURL(slack.ErrorInvalidToken, "slack://lol@12345678/123456789/123456789123456789123456")
		})
		ginkgo.It("should return an error if part B is not 9 letters", func() {
			expectErrorMessageGivenURL(slack.ErrorInvalidToken, "slack://lol@123456789/12345678/123456789123456789123456")
		})
		ginkgo.It("should return an error if part C is not 24 letters", func() {
			expectErrorMessageGivenURL(slack.ErrorInvalidToken, "slack://123456789/123456789/12345678912345678912345")
		})
	})
	ginkgo.When("given a token missing a part", func() {
		ginkgo.It("should return an error if the missing part is A", func() {
			expectErrorMessageGivenURL(slack.ErrorInvalidToken, "slack://lol@/123456789/123456789123456789123456")
		})
		ginkgo.It("should return an error if the missing part is B", func() {
			expectErrorMessageGivenURL(slack.ErrorInvalidToken, "slack://lol@123456789//123456789")
		})
		ginkgo.It("should return an error if the missing part is C", func() {
			expectErrorMessageGivenURL(slack.ErrorInvalidToken, "slack://lol@123456789/123456789/")
		})
	})
	ginkgo.Describe("the slack config", func() {
		ginkgo.When("parsing the configuration URL", func() {
			ginkgo.When("given a config using the legacy format", func() {
				ginkgo.It("should be converted to the new format after de-/serialization", func() {
					oldURL := "slack://testbot@AAAAAAAAA/BBBBBBBBB/123456789123456789123456?color=3f00fe&title=Test+title"
					newURL := "slack://hook:AAAAAAAAA-BBBBBBBBB-123456789123456789123456@webhook?botname=testbot&color=3f00fe&title=Test+title"

					config := &slack.Config{}
					err := config.SetURL(testutils.URLMust(oldURL))
					gomega.Expect(err).NotTo(gomega.HaveOccurred(), "verifying")

					gomega.Expect(config.GetURL().String()).To(gomega.Equal(newURL))
				})
			})
		})
		ginkgo.When("the URL contains an invalid property", func() {
			testURL := testutils.URLMust("slack://hook:AAAAAAAAA-BBBBBBBBB-123456789123456789123456@webhook?bass=dirty")
			err := (&slack.Config{}).SetURL(testURL)
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
		ginkgo.It("should be identical after de-/serialization", func() {
			testURL := "slack://hook:AAAAAAAAA-BBBBBBBBB-123456789123456789123456@webhook?botname=testbot&color=3f00fe&title=Test+title"

			config := &slack.Config{}
			err := config.SetURL(testutils.URLMust(testURL))
			gomega.Expect(err).NotTo(gomega.HaveOccurred(), "verifying")

			outputURL := config.GetURL()
			gomega.Expect(outputURL.String()).To(gomega.Equal(testURL))
		})
		ginkgo.When("generating a config object", func() {
			ginkgo.It("should use the default botname if the argument list contains three strings", func() {
				slackURL, _ := url.Parse("slack://AAAAAAAAA/BBBBBBBBB/123456789123456789123456")
				config, configError := slack.CreateConfigFromURL(slackURL)

				gomega.Expect(configError).NotTo(gomega.HaveOccurred())
				gomega.Expect(config.BotName).To(gomega.BeEmpty())
			})
			ginkgo.It("should set the botname if the argument list is three", func() {
				slackURL, _ := url.Parse("slack://testbot@AAAAAAAAA/BBBBBBBBB/123456789123456789123456")
				config, configError := slack.CreateConfigFromURL(slackURL)

				gomega.Expect(configError).NotTo(gomega.HaveOccurred())
				gomega.Expect(config.BotName).To(gomega.Equal("testbot"))
			})
			ginkgo.It("should return an error if the argument list is shorter than three", func() {
				slackURL, _ := url.Parse("slack://AAAAAAAA")

				_, configError := slack.CreateConfigFromURL(slackURL)
				gomega.Expect(configError).To(gomega.HaveOccurred())
			})
		})
		ginkgo.When("getting credentials from token", func() {
			ginkgo.It("should return a valid webhook URL for the given token", func() {
				token := tokenMust("AAAAAAAAA/BBBBBBBBB/123456789123456789123456")
				gomega.Expect(token.WebhookURL()).To(gomega.Equal(TestWebhookURL))
			})
			ginkgo.It("should return a valid authorization header value for the given token", func() {
				token := tokenMust("xoxb:AAAAAAAAA-BBBBBBBBB-123456789123456789123456")
				expected := "Bearer xoxb-AAAAAAAAA-BBBBBBBBB-123456789123456789123456"
				gomega.Expect(token.Authorization()).To(gomega.Equal(expected))
			})
		})
	})

	ginkgo.Describe("creating the payload", func() {
		ginkgo.Describe("the icon fields", func() {
			payload := slack.MessagePayload{}
			ginkgo.It("should set IconURL when the configured icon looks like an URL", func() {
				payload.SetIcon("https://example.com/logo.png")
				gomega.Expect(payload.IconURL).To(gomega.Equal("https://example.com/logo.png"))
				gomega.Expect(payload.IconEmoji).To(gomega.BeEmpty())
			})
			ginkgo.It("should set IconEmoji when the configured icon does not look like an URL", func() {
				payload.SetIcon("tanabata_tree")
				gomega.Expect(payload.IconEmoji).To(gomega.Equal("tanabata_tree"))
				gomega.Expect(payload.IconURL).To(gomega.BeEmpty())
			})
			ginkgo.It("should clear both fields when icon is empty", func() {
				payload.SetIcon("")
				gomega.Expect(payload.IconEmoji).To(gomega.BeEmpty())
				gomega.Expect(payload.IconURL).To(gomega.BeEmpty())
			})
		})
		ginkgo.When("when more than 99 lines are being sent", func() {
			ginkgo.It("should append the exceeding lines to the last attachment", func() {
				config := slack.Config{}
				sb := strings.Builder{}
				for i := 1; i <= 110; i++ {
					sb.WriteString(fmt.Sprintf("Line %d\n", i))
				}
				payload := slack.CreateJSONPayload(&config, sb.String()).(slack.MessagePayload)
				atts := payload.Attachments

				fmt.Printf("\nLines: %d, Last: %#v\n", len(atts), atts[len(atts)-1])

				gomega.Expect(atts).To(gomega.HaveLen(100))
				gomega.Expect(atts[len(atts)-1].Text).To(gomega.ContainSubstring("Line 110"))
			})
		})
		ginkgo.When("when the last message line ends with a newline", func() {
			ginkgo.It("should not send an empty attachment", func() {
				payload := slack.CreateJSONPayload(&slack.Config{}, "One\nTwo\nThree\n").(slack.MessagePayload)
				atts := payload.Attachments
				gomega.Expect(atts[len(atts)-1].Text).NotTo(gomega.BeEmpty())
			})
		})
	})

	ginkgo.Describe("sending the payload", func() {
		ginkgo.When("sending via webhook URL", func() {
			var err error
			ginkgo.BeforeEach(func() {
				httpmock.Activate()
			})
			ginkgo.AfterEach(func() {
				httpmock.DeactivateAndReset()
			})

			ginkgo.It("should not report an error if the server accepts the payload", func() {
				serviceURL, _ := url.Parse("slack://testbot@AAAAAAAAA/BBBBBBBBB/123456789123456789123456")
				err = service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				httpmock.RegisterResponder("POST", TestWebhookURL, httpmock.NewStringResponder(200, ""))

				err = service.Send("Message", nil)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
			ginkgo.It("should not panic if an error occurs when sending the payload", func() {
				serviceURL, _ := url.Parse("slack://testbot@AAAAAAAAA/BBBBBBBBB/123456789123456789123456")
				err = service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				httpmock.RegisterResponder("POST", TestWebhookURL, httpmock.NewErrorResponder(errors.New("dummy error")))

				err = service.Send("Message", nil)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
		ginkgo.When("sending via bot API", func() {
			var err error
			ginkgo.BeforeEach(func() {
				httpmock.Activate()
			})
			ginkgo.AfterEach(func() {
				httpmock.DeactivateAndReset()
			})

			ginkgo.It("should not report an error if the server accepts the payload", func() {
				serviceURL := testutils.URLMust("slack://xoxb:123456789012-1234567890123-4mt0t4l1YL3g1T5L4cK70k3N@C0123456789")
				err = service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				targetURL := "https://slack.com/api/chat.postMessage"
				httpmock.RegisterResponder("POST", targetURL, testutils.JSONRespondMust(200, slack.APIResponse{
					Ok: true,
				}))

				err = service.Send("Message", nil)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
			ginkgo.It("should not panic if an error occurs when sending the payload", func() {
				serviceURL := testutils.URLMust("slack://xoxb:123456789012-1234567890123-4mt0t4l1YL3g1T5L4cK70k3N@C0123456789")
				err = service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				targetURL := "https://slack.com/api/chat.postMessage"
				httpmock.RegisterResponder("POST", targetURL, testutils.JSONRespondMust(200, slack.APIResponse{
					Error: "someone turned off the internet",
				}))

				err = service.Send("Message", nil)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})
})

func tokenMust(rawToken string) *slack.Token {
	token, err := slack.ParseToken(rawToken)
	gomega.ExpectWithOffset(1, err).NotTo(gomega.HaveOccurred())

	return token
}

func expectErrorMessageGivenURL(expected error, rawURL string) {
	err := service.Initialize(testutils.URLMust(rawURL), testutils.TestLogger())
	gomega.ExpectWithOffset(1, err).To(gomega.HaveOccurred())
	gomega.ExpectWithOffset(1, err).To(gomega.Equal(expected))
}
