package telegram

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("the telegram service", func() {
	var logger *log.Logger

	BeforeEach(func() {
		logger = log.New(GinkgoWriter, "Test", log.LstdFlags)
	})

	Describe("creating configurations", func() {
		When("given an url", func() {
			When("a parse mode is supplied", func() {
				When("it's set to a valid mode and not None", func() {
					It("parse_mode should be present in payload", func() {
						payload, err := getPayloadStringFromURL("telegram://12345:mock-token@telegram/?channels=channel-1&parsemode=Markdown", "Message", logger)
						Expect(err).NotTo(HaveOccurred())
						Expect(payload).To(ContainSubstring("parse_mode"))
					})
				})
				When("it's set to None", func() {
					payload, err := getSinglePayloadFromURL("telegram://12345:mock-token@telegram/?channels=channel-1&title=MessageTitle", `Oh wow! <3 Cool & stuff ->`, logger)
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

func getPayloadsFromURL(testURL string, message string, logger *log.Logger) ([]SendMessagePayload, int, error) {
	var payloads []SendMessagePayload
	telegram := &Service{}

	serviceURL, err := url.Parse(testURL)
	if err != nil {
		return payloads, 0, err
	}

	if err = telegram.Initialize(serviceURL, logger); err != nil {
		return payloads, 0, err
	}

	if len(telegram.config.Chats) < 1 {
		return payloads, 0, errors.New("no channels were supplied")
	}

	messages, omitted := splitMessages(telegram.config, message)

	payloads = make([]SendMessagePayload, len(messages))
	for i, msg := range messages {
		payloads[i] = createSendMessagePayload(msg, telegram.config.Chats[0], telegram.config)
	}
	return payloads, omitted, nil

}

func getSinglePayloadFromURL(testURL string, message string, logger *log.Logger) (SendMessagePayload, error) {
	payloads, _, err := getPayloadsFromURL(testURL, message, logger)
	return payloads[0], err
}

func getPayloadStringFromURL(testURL string, message string, logger *log.Logger) ([]byte, error) {
	payload, err := getSinglePayloadFromURL(testURL, message, logger)
	if err != nil {
		return nil, err
	}

	raw, err := json.Marshal(payload)
	return raw, err
}
