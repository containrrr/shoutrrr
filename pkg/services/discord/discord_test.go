package discord_test

import (
	"github.com/containrrr/shoutrrr/pkg/services"
	. "github.com/containrrr/shoutrrr/pkg/services/discord"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestDiscord(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Discord Suite")
}

var (
	service       *Service
	envDiscordURL *url.URL
)

var _ = Describe("the discord service", func() {
	BeforeSuite(func() {
		service = &Service{}
		envDiscordURL, _ = url.Parse(os.Getenv("SHOUTRRR_DISCORD_URL"))
	})
	When("running integration tests", func() {
		It("should work without errors", func() {
			if envDiscordURL.String() == "" {
				return
			}

			serviceURL, _ := url.Parse(envDiscordURL.String())
			opts := services.GetDefaultOpts()
			err := service.Send(
				serviceURL,
				"this is an integration test",
				opts,
				)
			Expect(err).NotTo(HaveOccurred())
		})
	})
	Describe("creating a config", func() {
		When("given an url and a message", func() {
			It("should return an error if no arguments where supplied", func() {
				serviceURL, _ := url.Parse("discord://")
				_, err := service.CreateConfigFromURL(serviceURL)
				Expect(err).To(HaveOccurred())
			})
			It("should not return an error if exactly two arguments are given", func() {
				serviceURL, _ := url.Parse("discord://dummyToken@dummyChannel")
				_, err := service.CreateConfigFromURL(serviceURL)
				Expect(err).NotTo(HaveOccurred())
			})
			It("should return an error if more than two arguments are given", func() {
				serviceURL, _ := url.Parse("discord://dummyToken@dummyChannel/illegal-argument")
				_, err := service.CreateConfigFromURL(serviceURL)
				Expect(err).To(HaveOccurred())
			})
		})
	})
	Describe("creating a json payload", func() {
		When("given a blank message", func() {
			It("should return an error", func() {
				_, err := CreateJSONToSend("")
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
				_, err := CreateJSONToSend(builder.String())
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
