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
)

func TestSlack(t *testing.T) {
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

	When("given a token with a malformed part", func() {
		It("should return an error if part A is not 9 letters", func() {
			slackURL, err := url.Parse("slack://lol@12345678/123456789/123456789123456789123456")
			Expect(err).NotTo(HaveOccurred())
			expectErrorMessageGivenURL(
				TokenAMalformed,
				slackURL,
			)
		})
		It("should return an error if part B is not 9 letters", func() {
			slackURL, err := url.Parse("slack://lol@123456789/12345678/123456789123456789123456")
			Expect(err).NotTo(HaveOccurred())
			expectErrorMessageGivenURL(
				TokenBMalformed,
				slackURL,
			)
		})
		It("should return an error if part C is not 24 letters", func() {
			slackURL, err := url.Parse("slack://123456789/123456789/12345678912345678912345")
			Expect(err).NotTo(HaveOccurred())
			expectErrorMessageGivenURL(
				TokenCMalformed,
				slackURL,
			)
		})
	})
	When("given a token missing a part", func() {
		It("should return an error if the missing part is A", func() {
			slackURL, err := url.Parse("slack://lol@/123456789/123456789123456789123456")
			Expect(err).NotTo(HaveOccurred())
			expectErrorMessageGivenURL(
				TokenAMissing,
				slackURL,
			)
		})
		It("should return an error if the missing part is B", func() {
			slackURL, err := url.Parse("slack://lol@123456789//123456789")
			Expect(err).NotTo(HaveOccurred())
			expectErrorMessageGivenURL(
				TokenBMissing,
				slackURL,
			)

		})
		It("should return an error if the missing part is C", func() {
			slackURL, err := url.Parse("slack://lol@123456789/123456789/")
			Expect(err).NotTo(HaveOccurred())
			expectErrorMessageGivenURL(
				TokenCMissing,
				slackURL,
			)
		})
	})
	Describe("the slack config", func() {
		When("parsing the configuration URL", func() {
			It("should be identical after de-/serialization", func() {
				testURL := "slack://testbot@AAAAAAAAA/BBBBBBBBB/123456789123456789123456?color=3f00fe&title=Test title"

				url, err := url.Parse(testURL)
				Expect(err).NotTo(HaveOccurred(), "parsing")

				config := &Config{}
				err = config.SetURL(url)
				Expect(err).NotTo(HaveOccurred(), "verifying")

				outputURL := config.GetURL()
				Expect(outputURL.String()).To(Equal(testURL))

			})
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
})

func expectErrorMessageGivenURL(msg ErrorMessage, slackURL *url.URL) {
	err := service.Initialize(slackURL, util.TestLogger())
	Expect(err).To(HaveOccurred())
	Expect(err.Error()).To(Equal(string(msg)))
}
