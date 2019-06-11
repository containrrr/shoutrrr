package slack_test

import (
	. "github.com/containrrr/shoutrrr/pkg/services/slack"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/url"
	"os"
	"testing"
)

func TestShoutrrr(t *testing.T) {
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

			err := service.Send(serviceURL, "This is an integration test message",nil)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	When("given a token with a malformed part", func() {
		It("should return an error if part A is not 9 letters", func() {
			slackURL, _ := url.Parse("slack://lol@12345678/123456789/123456789123456789123456")
			expectErrorMessageGivenUrl(
				TokenAMalformed,
				slackURL,
			)
		})
		It("should return an error if part B is not 9 letters", func() {
			slackURL, _ := url.Parse("slack://lol@123456789/12345678/123456789123456789123456")
			expectErrorMessageGivenUrl(
				TokenBMalformed,
				slackURL,
			)
		})
		It("should return an error if part C is not 24 letters", func() {
			slackURL, _ := url.Parse("slack://123456789/123456789/12345678912345678912345")
			expectErrorMessageGivenUrl(
				TokenCMalformed,
				slackURL,
			)
		})
	})
	When("given a token missing a part", func() {
		It("should return an error if the missing part is A", func() {
			slackURL, _ := url.Parse("slack://lol@/123456789/123456789123456789123456")
			expectErrorMessageGivenUrl(
				TokenAMissing,
				slackURL,
			)
		})
		It("should return an error if the missing part is B", func() {
			slackURL, _ := url.Parse("slack://lol@123456789//123456789")
			expectErrorMessageGivenUrl(
				TokenBMissing,
				slackURL,
			)

		})
		It("should return an error if the missing part is C", func() {
			slackURL, _ := url.Parse("slack://lol@123456789/123456789/")
			expectErrorMessageGivenUrl(
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

func expectErrorMessageGivenUrl(msg ErrorMessage, slackUrl *url.URL) {
	err := service.Send(slackUrl, "Hello", nil)
	Expect(err).To(HaveOccurred())
	Expect(err.Error()).To(Equal(string(msg)))
}
