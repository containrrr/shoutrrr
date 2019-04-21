package slack_test

import (
	"os"
	"testing"
	. "github.com/containrrr/shoutrrr/pkg/plugins/slack"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestShoutrrr(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Slack Suite")
}

var (
	plugin      *SlackPlugin
	envSlackUrl string
)

var _ = Describe("the slack plugin", func() {

	BeforeSuite(func() {
		plugin = &SlackPlugin{}
		envSlackUrl = os.Getenv("SHOUTRRR_SLACK_URL")

	})

	When("running integration tests", func() {
		It("should not error out", func() {
			if envSlackUrl == "" {
				return
			}
			err := plugin.Send(envSlackUrl, "This is an integration test message")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	When("given a token with a malformed part", func() {
		It("should return an error if part A is not 9 letters", func() {
			expectErrorMessageGivenUrl(
				TokenAMalformed,
				"slack://lol/12345678/123456789/123456789123456789123456")
		})
		It("should return an error if part B is not 9 letters", func() {
			expectErrorMessageGivenUrl(
				TokenBMalformed,
				"slack://lol/123456789/12345678/123456789123456789123456")
		})
		It("should return an error if part C is not 24 letters", func() {
			expectErrorMessageGivenUrl(
				TokenCMalformed,
				"slack://123456789/123456789/12345678912345678912345")
		})
	})
	When("given a token missing a part", func() {
		It("should return an error if the missing part is A", func() {
			expectErrorMessageGivenUrl(
				TokenAMissing,
				"slack://lol//123456789/123456789123456789123456")
		})
		It("should return an error if the missing part is B", func() {
			expectErrorMessageGivenUrl(
				TokenBMissing,
				"slack://lol/123456789//123456789")

		})
		It("should return an error if the missing part is C", func() {
			expectErrorMessageGivenUrl(
				TokenCMissing,
				"slack://lol/123456789/123456789/")
		})
	})
	Describe("the slack config", func() {
		When("generating a config object", func() {
			It("should use the default botname if the argument list contains three strings", func() {
				url := "slack://AAAAAAAAA/BBBBBBBBB/123456789123456789123456"
				config, configError := CreateConfigFromUrl(url)

				Expect(config.Botname).To(Equal(DefaultUser))
				Expect(configError).NotTo(HaveOccurred())

			})
			It("should set the botname if the argument list is larger than three", func() {
				url := "slack://testbot/AAAAAAAAA/BBBBBBBBB/123456789123456789123456"
				config, configError := CreateConfigFromUrl(url)

				Expect(configError).NotTo(HaveOccurred())
				Expect(config.Botname).To(Equal("testbot"))
			})
			It("should return an error if the argument list is shorter than three", func() {
				url := "slack://AAAAAAAA"
				_, configError := CreateConfigFromUrl(url)
				Expect(configError).To(HaveOccurred())
			})
		})
	})
})

func expectErrorMessageGivenUrl(msg ErrorMessage, url string) {
	err := plugin.Send(url, "Hello")
	Expect(err).To(HaveOccurred())
	Expect(err.Error()).To(Equal(string(msg)))
}
