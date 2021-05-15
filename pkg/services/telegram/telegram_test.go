package telegram_test

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"log"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/containrrr/shoutrrr/internal/testutils"
	. "github.com/containrrr/shoutrrr/pkg/services/telegram"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTelegram(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Telegram Suite")
}

var _ = Describe("the telegram service", func() {
	var telegram *Service
	var envTelegramURL string
	var logger *log.Logger

	BeforeSuite(func() {
		envTelegramURL = os.Getenv("SHOUTRRR_TELEGRAM_URL")
		logger = log.New(GinkgoWriter, "Test", log.LstdFlags)
	})

	BeforeEach(func() {
		telegram = &Service{}
	})

	When("running integration tests", func() {
		It("should not error out", func() {
			if envTelegramURL == "" {
				return
			}
			serviceURL, _ := url.Parse(envTelegramURL)
			err := telegram.Initialize(serviceURL, logger)
			Expect(err).NotTo(HaveOccurred())
			err = telegram.Send("This is an integration test message", nil)
			Expect(err).NotTo(HaveOccurred())
		})
		When("given a message that exceeds the max length", func() {
			It("should generate an error", func() {
				if envTelegramURL == "" {
					return
				}
				hundredChars := "this string is exactly (to the letter) a hundred characters long which will make the send func error"
				serviceURL, _ := url.Parse("telegram://12345:mock-token/?chats=channel-1")
				builder := strings.Builder{}
				for i := 0; i < 42; i++ {
					builder.WriteString(hundredChars)
				}

				err := telegram.Initialize(serviceURL, logger)
				Expect(err).NotTo(HaveOccurred())
				err = telegram.Send(builder.String(), nil)
				Expect(err).To(HaveOccurred())
			})
		})
		When("given a valid request with a faked token", func() {
			if envTelegramURL == "" {
				return
			}
			It("should generate a 401", func() {
				serviceURL, _ := url.Parse("telegram://000000000:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA@telegram/?chats=channel-id")
				message := "this is a perfectly valid message"

				err := telegram.Initialize(serviceURL, logger)
				Expect(err).NotTo(HaveOccurred())
				err = telegram.Send(message, nil)
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
					serviceURL, _ := url.Parse("telegram://12345:mock-token@telegram/?chats=channel-1,channel-2,channel-3")
					err = telegram.Initialize(serviceURL, logger)
					config = telegram.GetConfig()
				})

				It("should create a config object", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(config != nil).To(BeTrue())
				})
				It("should create a config object containing the API Token", func() {

					Expect(err).NotTo(HaveOccurred())
					Expect(config.Token).To(Equal("12345:mock-token"))
				})
				It("should add every chats query field as a chat ID", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(config.Chats).To(Equal([]string{
						"channel-1",
						"channel-2",
						"channel-3",
					}))
				})
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
			serviceURL, _ := url.Parse("telegram://12345:mock-token@telegram/?chats=channel-1,channel-2,channel-3")
			err = telegram.Initialize(serviceURL, logger)
			Expect(err).NotTo(HaveOccurred())

			setupResponder("sendMessage", telegram.GetConfig().Token, 200, "")

			err = telegram.Send("Message", nil)
			Expect(err).NotTo(HaveOccurred())
		})

	})

	It("should implement basic service API methods correctly", func() {
		testutils.TestConfigGetInvalidQueryValue(&Config{})
		testutils.TestConfigSetInvalidQueryValue(&Config{}, "telegram://12345:mock-token@telegram/?chats=channel-1&foo=bar")

		testutils.TestConfigGetEnumsCount(&Config{}, 1)
		testutils.TestConfigGetFieldsCount(&Config{}, 5)
	})
})

func expectErrorAndEmptyObject(telegram *Service, rawURL string, logger *log.Logger) {
	serviceURL, _ := url.Parse(rawURL)
	err := telegram.Initialize(serviceURL, logger)
	Expect(err).To(HaveOccurred())
	config := telegram.GetConfig()
	fmt.Printf("Token: \"%+v\" \"%s\" \n", config.Token, config.Token)
	Expect(config.Token).To(BeEmpty())
	Expect(len(config.Chats)).To(BeZero())
}

func setupResponder(endpoint string, token string, code int, body string) {
	targetUrl := fmt.Sprintf("https://api.telegram.org/bot%s/%s", token, endpoint)
	httpmock.RegisterResponder("POST", targetUrl, httpmock.NewStringResponder(code, body))
}
