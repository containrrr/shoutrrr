package telegram

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/nicholas-fedor/shoutrrr/internal/testutils"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestTelegram(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Shoutrrr Telegram Suite")
}

var (
	envTelegramURL string
	logger         *log.Logger

	_ = ginkgo.BeforeSuite(func() {
		envTelegramURL = os.Getenv("SHOUTRRR_TELEGRAM_URL")
		logger = log.New(ginkgo.GinkgoWriter, "Test", log.LstdFlags)
	})
)

var _ = ginkgo.Describe("the telegram service", func() {
	var telegram *Service // No telegram. prefix needed

	ginkgo.BeforeEach(func() {
		telegram = &Service{}
	})

	ginkgo.When("running integration tests", func() {
		ginkgo.It("should not error out", func() {
			if envTelegramURL == "" {
				return
			}
			serviceURL, _ := url.Parse(envTelegramURL)
			err := telegram.Initialize(serviceURL, logger)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			err = telegram.Send("This is an integration test Message", nil)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
		ginkgo.When("given a Message that exceeds the max length", func() {
			ginkgo.It("should generate an error", func() {
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
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				err = telegram.Send(builder.String(), nil)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
		ginkgo.When("given a valid request with a faked token", func() {
			if envTelegramURL == "" {
				return
			}
			ginkgo.It("should generate a 401", func() {
				serviceURL, _ := url.Parse("telegram://000000000:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA@telegram/?chats=channel-id")
				message := "this is a perfectly valid Message"

				err := telegram.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				err = telegram.Send(message, nil)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(strings.Contains(err.Error(), "401 Unauthorized")).To(gomega.BeTrue())
			})
		})
	})

	ginkgo.Describe("creating configurations", func() {
		ginkgo.When("given an url", func() {
			ginkgo.It("should return an error if no arguments where supplied", func() {
				expectErrorAndEmptyObject(telegram, "telegram://", logger)
			})
			ginkgo.It("should return an error if the token has an invalid format", func() {
				expectErrorAndEmptyObject(telegram, "telegram://invalid-token", logger)
			})
			ginkgo.It("should return an error if only the api token where supplied", func() {
				expectErrorAndEmptyObject(telegram, "telegram://12345:mock-token@telegram", logger)
			})

			ginkgo.When("the url is valid", func() {
				var config *Config // No telegram. prefix
				var err error

				ginkgo.BeforeEach(func() {
					serviceURL, _ := url.Parse("telegram://12345:mock-token@telegram/?chats=channel-1,channel-2,channel-3")
					err = telegram.Initialize(serviceURL, logger)
					config = telegram.GetConfig()
				})

				ginkgo.It("should create a config object", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(config != nil).To(gomega.BeTrue())
				})
				ginkgo.It("should create a config object containing the API Token", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(config.Token).To(gomega.Equal("12345:mock-token"))
				})
				ginkgo.It("should add every chats query field as a chat ID", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(config.Chats).To(gomega.Equal([]string{
						"channel-1",
						"channel-2",
						"channel-3",
					}))
				})
			})
		})
	})

	ginkgo.Describe("sending the payload", func() {
		var err error
		ginkgo.BeforeEach(func() {
			httpmock.Activate()
		})
		ginkgo.AfterEach(func() {
			httpmock.DeactivateAndReset()
		})
		ginkgo.It("should not report an error if the server accepts the payload", func() {
			serviceURL, _ := url.Parse("telegram://12345:mock-token@telegram/?chats=channel-1,channel-2,channel-3")
			err = telegram.Initialize(serviceURL, logger)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			setupResponder("sendMessage", telegram.GetConfig().Token, 200, "")

			err = telegram.Send("Message", nil)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
	})

	ginkgo.It("should implement basic service API methods correctly", func() {
		serviceURL, _ := url.Parse("telegram://12345:mock-token@telegram/?chats=channel-1")
		err := telegram.Initialize(serviceURL, logger)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		config := telegram.GetConfig()
		testutils.TestConfigGetInvalidQueryValue(config)
		testutils.TestConfigSetInvalidQueryValue(config, "telegram://12345:mock-token@telegram/?chats=channel-1&foo=bar")
		testutils.TestConfigGetEnumsCount(config, 1)
		testutils.TestConfigGetFieldsCount(config, 6)
	})
	ginkgo.It("should return the correct service ID", func() {
		service := &Service{}
		gomega.Expect(service.GetID()).To(gomega.Equal("telegram"))
	})
})

func expectErrorAndEmptyObject(telegram *Service, rawURL string, logger *log.Logger) {
	serviceURL, _ := url.Parse(rawURL)
	err := telegram.Initialize(serviceURL, logger)
	gomega.Expect(err).To(gomega.HaveOccurred())

	config := telegram.GetConfig()
	gomega.Expect(config.Token).To(gomega.BeEmpty())
	gomega.Expect(len(config.Chats)).To(gomega.BeZero())
}

func setupResponder(endpoint string, token string, code int, body string) {
	targetUrl := fmt.Sprintf("https://api.telegram.org/bot%s/%s", token, endpoint)
	httpmock.RegisterResponder("POST", targetUrl, httpmock.NewStringResponder(code, body))
}
