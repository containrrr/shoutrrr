package slack_test

import (
	. "github.com/containrrr/shoutrrr/pkg/services/slack"
	"github.com/containrrr/shoutrrr/pkg/util"

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
)

var _ = Describe("the slack service", func() {

	BeforeSuite(func() {
		service = &Service{}
		envSlackURL, _ = url.Parse(os.Getenv("SHOUTRRR_SLACK_URL"))

	})

	When("running integration tests", func() {
		It("should not error out", func() {
			if envSlackURL.String() == "" {
				return
			}

			serviceURL, _ := url.Parse(envSlackURL.String())
			service.Initialize(serviceURL, util.TestLogger())
			err := service.Send("This is an integration test message", nil)
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
		When("generating a config object", func() {
			It("should use the default botname if the argument list contains three strings", func() {
				slackURL, _ := url.Parse("slack://AAAAAAAAA/BBBBBBBBB/123456789123456789123456")
				config, configError := CreateConfigFromURL(slackURL)

				Expect(config.BotName).To(Equal(DefaultUser))
				Expect(configError).NotTo(HaveOccurred())
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
})

func expectErrorMessageGivenURL(msg ErrorMessage, slackURL *url.URL) {
	err := service.Initialize(slackURL, util.TestLogger())
	Expect(err).To(HaveOccurred())
	Expect(err.Error()).To(Equal(string(msg)))
}
