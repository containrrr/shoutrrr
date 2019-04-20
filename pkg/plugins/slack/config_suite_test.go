package slack_test

import (
    "github.com/containrrr/shoutrrr/pkg/plugins/slack"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "testing"
)

func TestShoutrrr(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Shoutrrr Suite")
}

const defaultUser = "Shoutrrr"

var _ = Describe("the slack config", func() {
    When("generating a config object", func() {
        It("should use the default botname if the argument list contains three strings", func() {
            url := "slack://AAAAAAAAA/BBBBBBBBB/123456789123456789123456"
            config, configError := slack.CreateConfigFromUrl(url)

            Expect(config.Botname).To(Equal(defaultUser))
            Expect(configError).NotTo(HaveOccurred())

        })
        It("should set the botname if the argument list is larger than three", func() {
            url := "slack://testbot/AAAAAAAAA/BBBBBBBBB/123456789123456789123456"
            config, configError := slack.CreateConfigFromUrl(url)

            Expect(configError).NotTo(HaveOccurred())
            Expect(config.Botname).To(Equal("testbot"))
        })
        It("should return an error if the argument list is shorter than three", func() {
            url := "slack://AAAAAAAA"
            _, configError := slack.CreateConfigFromUrl(url)
            Expect(configError).To(HaveOccurred())
        })
    })
})