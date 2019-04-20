package telegram_test

import (
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"testing"

	. "github.com/containrrr/shoutrrr/pkg/plugins/telegram"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTelegram(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Telegram Suite")
}

var _ = Describe("the telegram plugin", func() {
	var telegram *TelegramPlugin
	var envTelegramUrl string

	BeforeSuite(func() {
		logrus.SetLevel(logrus.DebugLevel)
		telegram = &TelegramPlugin{}
		envTelegramUrl = os.Getenv("SHOUTRRR_TELEGRAM_URL")

	})


	When("running integration tests", func() {
		It("should not error out", func() {
			if envTelegramUrl == "" {
				return
			}
			err := telegram.Send(envTelegramUrl, "This is an integration test message")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("sending a message", func() {
		When("given a valid request with a faked token", func() {
			It("should generate a 400", func() {
				url := "telegram://703391768:AAEWjOpAH_szG7Ym-WaaiPp6emexFc13uf0/channel-id"
				message := "this is a perfectly valid message"
				err := telegram.Send(url, message)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "400 Bad Request")).To(BeTrue())
			})
		})
		When("given a message that exceeds the max length", func() {
			It("should generate an error", func() {
				hundredChars := "this string is exactly (to the letter) a hundred characters long which will make the send func error"
				url := "telegram://12345:mock-token/channel-1"
				builder := strings.Builder{}
				for i := 0; i < 42; i++ {
					builder.WriteString(hundredChars)
				}

				err := telegram.Send(url, builder.String())
				Expect(err).To(HaveOccurred())
			})
		})
	})
	Describe("creating configurations", func() {
		When("given an url", func() {
			It("should return an error if no arguments where supplied", func() {
				url := "telegram://"
				config, err := telegram.CreateConfigFromUrl(url)
				Expect(err).To(HaveOccurred())
				Expect(config == nil).To(BeTrue())
			})
			It("should return an error if the token has an invalid format", func() {
				url := "telegram://invalid-token"
				config, err := telegram.CreateConfigFromUrl(url)
				Expect(err).To(HaveOccurred())
				Expect(config == nil).To(BeTrue())
			})
			It("should return an error if only the api token where supplied", func() {
				url := "telegram://12345:mock-token"
				config, err := telegram.CreateConfigFromUrl(url)
				Expect(err).To(HaveOccurred())
				Expect(config == nil).To(BeTrue())
			})
			It("should create a config object", func() {
				url := "telegram://12345:mock-token/channel-1/channel-2/channel-3"
				config, err := telegram.CreateConfigFromUrl(url)
				Expect(err).NotTo(HaveOccurred())
				Expect(config != nil).To(BeTrue())
			})
			It("should create a config object containing the API Token", func() {
				url := "telegram://12345:mock-token/channel-1/channel-2/channel-3"
				config, err := telegram.CreateConfigFromUrl(url)
				Expect(err).NotTo(HaveOccurred())
				Expect(config.ApiToken).To(Equal("12345:mock-token"))
			})
			It("should add every subsequent argument as a channel id", func() {
				url := "telegram://12345:mock-token/channel-1/channel-2/channel-3"
				config, err := telegram.CreateConfigFromUrl(url)
				Expect(err).NotTo(HaveOccurred())
				Expect(config.Channels).To(Equal([]string {
					"channel-1",
					"channel-2",
					"channel-3",
				}))
			})
		})
	})
})

