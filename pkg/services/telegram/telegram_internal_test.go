package telegram

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
)

var _ = ginkgo.Describe("the telegram service", func() {
	var logger *log.Logger

	ginkgo.BeforeEach(func() {
		logger = log.New(ginkgo.GinkgoWriter, "Test", log.LstdFlags)
	})

	ginkgo.Describe("creating configurations", func() {
		ginkgo.When("given an url", func() {
			ginkgo.When("a parse mode is not supplied", func() {
				ginkgo.It("no parse_mode should be present in payload", func() {
					payload, err := getPayloadStringFromURL("telegram://12345:mock-token@telegram/?channels=channel-1", "Message", logger)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(payload).NotTo(gomega.ContainSubstring("parse_mode"))
				})
			})

			ginkgo.When("a parse mode is supplied", func() {
				ginkgo.When("it's set to a valid mode and not None", func() {
					ginkgo.It("parse_mode should be present in payload", func() {
						payload, err := getPayloadStringFromURL("telegram://12345:mock-token@telegram/?channels=channel-1&parsemode=Markdown", "Message", logger)
						gomega.Expect(err).NotTo(gomega.HaveOccurred())
						gomega.Expect(payload).To(gomega.ContainSubstring("parse_mode"))
					})
				})
				ginkgo.When("it's set to None", func() {
					ginkgo.When("no title has been provided", func() {
						ginkgo.It("no parse_mode should be present in payload", func() {
							payload, err := getPayloadStringFromURL("telegram://12345:mock-token@telegram/?channels=channel-1&parsemode=None", "Message", logger)
							gomega.Expect(err).NotTo(gomega.HaveOccurred())
							gomega.Expect(payload).NotTo(gomega.ContainSubstring("parse_mode"))
						})
					})
					ginkgo.When("a title has been provided", func() {
						payload, err := getPayloadFromURL("telegram://12345:mock-token@telegram/?channels=channel-1&title=MessageTitle", `Oh wow! <3 Cool & stuff ->`, logger)
						gomega.Expect(err).NotTo(gomega.HaveOccurred())
						ginkgo.It("should have parse_mode set to HTML", func() {
							gomega.Expect(payload.ParseMode).To(gomega.Equal("HTML"))
						})
						ginkgo.It("should contain the title prepended in the message", func() {
							gomega.Expect(payload.Text).To(gomega.ContainSubstring("MessageTitle"))
						})
						ginkgo.It("should escape the message HTML tags", func() {
							gomega.Expect(payload.Text).To(gomega.ContainSubstring("&lt;3"))
							gomega.Expect(payload.Text).To(gomega.ContainSubstring("Cool &amp; stuff"))
							gomega.Expect(payload.Text).To(gomega.ContainSubstring("-&gt;"))
						})
					})
				})
			})

			ginkgo.When("parsing URL that might have a message thread id", func() {
				ginkgo.When("no thread id is provided", func() {
					payload, err := getPayloadFromURL("telegram://12345:mock-token@telegram/?channels=channel-1&title=MessageTitle", `Oh wow! <3 Cool & stuff ->`, logger)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					ginkgo.It("should have message_thread_id set to nil", func() {
						gomega.Expect(payload.MessageThreadID).To(gomega.BeNil())
					})
				})
				ginkgo.When("a numeric thread id is provided", func() {
					payload, err := getPayloadFromURL("telegram://12345:mock-token@telegram/?channels=channel-1:10&title=MessageTitle", `Oh wow! <3 Cool & stuff ->`, logger)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					ginkgo.It("should have message_thread_id set to 10", func() {
						gomega.Expect(payload.MessageThreadID).To(gstruct.PointTo(gomega.Equal(10)))
					})
				})
				ginkgo.When("non-numeric thread id is provided", func() {
					payload, err := getPayloadFromURL("telegram://12345:mock-token@telegram/?channels=channel-1:invalid&title=MessageTitle", `Oh wow! <3 Cool & stuff ->`, logger)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					ginkgo.It("should have message_thread_id set to nil", func() {
						gomega.Expect(payload.MessageThreadID).To(gomega.BeNil())
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

	if len(telegram.Config.Chats) < 1 {
		return SendMessagePayload{}, errors.New("no channels were supplied")
	}

	return createSendMessagePayload(message, telegram.Config.Chats[0], telegram.Config), nil
}

func getPayloadStringFromURL(testURL string, message string, logger *log.Logger) ([]byte, error) {
	payload, err := getPayloadFromURL(testURL, message, logger)
	if err != nil {
		return nil, err
	}

	raw, err := json.Marshal(payload)

	return raw, err
}
