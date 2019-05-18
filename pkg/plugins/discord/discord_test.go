package discord_test

import (
	. "github.com/containrrr/shoutrrr/pkg/plugins/discord"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"strings"
	"testing"
)

func TestDiscord(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Discord Suite")
}

var (
	plugin        *DiscordPlugin
	envDiscordURL string
)

var _ = Describe("the discord plugin", func() {
	BeforeSuite(func() {
		plugin = &DiscordPlugin{}
		envDiscordURL = os.Getenv("SHOUTRRR_DISCORD_URL")
	})
	When("running integration tests", func() {
		It("should work without errors", func() {
			if envDiscordURL == "" {
				return
			}
			err := plugin.Send(envDiscordURL, "this is an integration test")
			Expect(err).NotTo(HaveOccurred())
		})
	})
	Describe("creating a config", func() {
		When("given an url and a message", func() {
			It("should return an error if no arguments where supplied", func() {
				url := "discord://"
				_, err := plugin.CreateConfigFromURL(url)
				Expect(err).To(HaveOccurred())
			})
			It("should not return an error if exactly two arguments are given", func() {
				url := "discord://channel/token"
				_, err := plugin.CreateConfigFromURL(url)
				Expect(err).NotTo(HaveOccurred())
			})
			It("should return an error if more than two arguments are given", func() {
				url := "discord://channel/token/illegal-argument"
				_, err := plugin.CreateConfigFromURL(url)
				Expect(err).To(HaveOccurred())
			})
		})
	})
	Describe("creating a json payload", func() {
		When("given a blank message", func() {
			It("should return an error", func() {
				_, err := CreateJsonToSend("")
				Expect(err).To(HaveOccurred())
			})
		})
		When("given a message that exceeds the max length", func() {
			It("should return an error", func() {
				hundredChars := "this string is exactly (to the letter) a hundred characters long which will make the send func error"
				builder := strings.Builder{}

				for i := 0; i < 42; i++ {
					builder.WriteString(hundredChars)
				}
				_, err := CreateJsonToSend(builder.String())
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
