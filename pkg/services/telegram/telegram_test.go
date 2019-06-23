package telegram_test

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"testing"

	. "github.com/containrrr/shoutrrr/pkg/services/telegram"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTelegram(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Telegram Suite")
}

var _ = Describe("the telegram plugin", func() {
	var telegram *Service
	var envTelegramUrl string
	var logger *log.Logger

	BeforeSuite(func() {
		envTelegramUrl = os.Getenv("SHOUTRRR_TELEGRAM_URL")
		logger = log.New(GinkgoWriter, "Test", log.LstdFlags)
	})

	BeforeEach(func() {
		telegram = &Service{}
	})


	When("running integration tests", func() {
		It("should not error out", func() {
			if envTelegramUrl == "" {
				return
			}
			serviceURL, _ := url.Parse(envTelegramUrl)
			telegram.Initialize(telegram.NewConfig(), serviceURL, logger)
			err := telegram.Send("This is an integration test message", nil)
			Expect(err).NotTo(HaveOccurred())
		})
		When("given a message that exceeds the max length", func() {
			It("should generate an error", func() {
				if envTelegramUrl == "" {
					return
				}
				hundredChars := "this string is exactly (to the letter) a hundred characters long which will make the send func error"
				serviceURL, _ := url.Parse("telegram://12345:mock-token/channel-1")
				builder := strings.Builder{}
				for i := 0; i < 42; i++ {
					builder.WriteString(hundredChars)
				}

				telegram.Initialize(telegram.NewConfig(), serviceURL, logger)
				err := telegram.Send(builder.String(), nil)
				Expect(err).To(HaveOccurred())
			})
		})
		When("given a valid request with a faked token", func() {
			if envTelegramUrl == "" {
				return
			}
			It("should generate a 401", func() {
				serviceURL, _ := url.Parse("telegram://000000000:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA@telegram/?channels=channel-id")
				message := "this is a perfectly valid message"

				telegram.Initialize(telegram.NewConfig(), serviceURL, logger)
				err := telegram.Send( message, nil)
				Expect(err).To(HaveOccurred())
				fmt.Println(err.Error())
				Expect(strings.Contains(err.Error(), "401 Unauthorized")).To(BeTrue())
			})
		})
	})

	Describe("creating configurations", func() {
		When("given an url", func() {
			It("should return an error if no arguments where supplied", func() {
				expectErrorAndEmptyObject(telegram, "telegram://", logger)
			})
			It("should return an error if the token has an invalid format", func() {
				expectErrorAndEmptyObject(telegram, "telegram://invalid-token", logger)
			})
			It("should return an error if only the api token where supplied", func() {
				expectErrorAndEmptyObject(telegram, "telegram://12345:mock-token@telegram", logger)
			})
			When("the url is valid", func() {
				var config *Config
				var err error

				BeforeEach(func() {
					serviceURL, _ := url.Parse("telegram://12345:mock-token@telegram/?channels=channel-1,channel-2,channel-3")
					config = new(Config)
					err = telegram.Initialize(config, serviceURL, logger)
				})

				It("should create a config object", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(config != nil).To(BeTrue())
				})
				It("should create a config object containing the API Token", func() {

					Expect(err).NotTo(HaveOccurred())
					Expect(config.Token).To(Equal("12345:mock-token"))
				})
				It("should add every subsequent argument as a channel id", func() {
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
})

func expectErrorAndEmptyObject(telegram *Service, rawURL string, logger *log.Logger) {
	serviceURL, _ := url.Parse(rawURL)
	config := new(Config)
	err := telegram.Initialize(config, serviceURL, logger)
	Expect(err).To(HaveOccurred())
	fmt.Printf("Token: \"%+v\" \"%s\" \n", config.Token, config.Token)
	Expect(config.Token).To(BeEmpty())
	Expect(len(config.Channels)).To(BeZero())
}