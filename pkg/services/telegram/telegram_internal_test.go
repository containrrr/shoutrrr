package telegram

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("the telegram service", func() {
	var logger *log.Logger

	BeforeEach(func() {
		logger = log.New(GinkgoWriter, "Test", log.LstdFlags)
	})

	Describe("creating configurations", func() {
		When("given an url", func() {

			When("a parse mode is not supplied", func() {
				It("no parse_mode should be present in payload", func() {
					payload, err := getPayloadStringFromURL("telegram://12345:mock-token@telegram/?channels=channel-1", "Message", logger)
					Expect(err).NotTo(HaveOccurred())
					Expect(payload).NotTo(ContainSubstring("parse_mode"))
				})
			})

			When("a parse mode is supplied", func() {
				When("it's set to a valid mode and not None", func() {
					It("parse_mode should be present in payload", func() {
						payload, err := getPayloadStringFromURL("telegram://12345:mock-token@telegram/?channels=channel-1&parsemode=Markdown", "Message", logger)
						Expect(err).NotTo(HaveOccurred())
						Expect(payload).To(ContainSubstring("parse_mode"))
					})
				})
				When("it's set to None", func() {
					When("no title has been provided", func() {
						It("no parse_mode should be present in payload", func() {
							payload, err := getPayloadStringFromURL("telegram://12345:mock-token@telegram/?channels=channel-1&parsemode=None", "Message", logger)
							Expect(err).NotTo(HaveOccurred())
							Expect(payload).NotTo(ContainSubstring("parse_mode"))
						})
					})
					When("a title has been provided", func() {
						payload, err := getPayloadFromURL("telegram://12345:mock-token@telegram/?channels=channel-1&title=MessageTitle", `Oh wow! <3 Cool & stuff ->`, logger)
						Expect(err).NotTo(HaveOccurred())
						It("should have parse_mode set to HTML", func() {
							Expect(payload.ParseMode).To(Equal("HTML"))
						})
						It("should contain the title prepended in the message", func() {
							Expect(payload.Text).To(ContainSubstring("MessageTitle"))
						})
						It("should escape the message HTML tags", func() {
							Expect(payload.Text).To(ContainSubstring("&lt;3"))
							Expect(payload.Text).To(ContainSubstring("Cool &amp; stuff"))
							Expect(payload.Text).To(ContainSubstring("-&gt;"))
						})
					})
				})
			})

		})
	})
})

func getPayloadFromURL(testURL string, message string, logger *log.Logger) (SendMessagePayload, error) {
	telegram := &Service{}

	serviceURL, err := url.Parse(testURL)
	if err != nil {
		return SendMessagePayload{}, err
	}

	if err = telegram.Initialize(serviceURL, logger); err != nil {
		return SendMessagePayload{}, err
	}

	if len(telegram.config.Chats) < 1 {
		return SendMessagePayload{}, errors.New("no channels were supplied")
	}

	return createSendMessagePayload(message, telegram.config.Chats[0], telegram.config), nil

}

func getPayloadStringFromURL(testURL string, message string, logger *log.Logger) ([]byte, error) {
	payload, err := getPayloadFromURL(testURL, message, logger)
	if err != nil {
		return nil, err
	}

	raw, err := json.Marshal(payload)
	return raw, err
}
