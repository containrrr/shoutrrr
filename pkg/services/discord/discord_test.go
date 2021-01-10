package discord_test

import (
	. "github.com/containrrr/shoutrrr/pkg/services/discord"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util"
	"github.com/jarcoal/httpmock"
	"log"
	"time"

	"net/url"
	"os"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDiscord(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Discord Suite")
}

var (
	dummyColors   = [types.MessageLevelCount]int{}
	service       *Service
	envDiscordURL *url.URL
	logger        *log.Logger
)

var _ = Describe("the discord service", func() {
	BeforeSuite(func() {
		service = &Service{}
		envDiscordURL, _ = url.Parse(os.Getenv("SHOUTRRR_DISCORD_URL"))
		logger = log.New(GinkgoWriter, "Test", log.LstdFlags)
	})
	When("running integration tests", func() {
		It("should work without errors", func() {
			if envDiscordURL.String() == "" {
				return
			}

			serviceURL, _ := url.Parse(envDiscordURL.String())
			err := service.Initialize(serviceURL, util.TestLogger())
			Expect(err).NotTo(HaveOccurred())

			err = service.Send(
				"this is an integration test",
				nil,
			)
			Expect(err).NotTo(HaveOccurred())
		})
	})
	Describe("the service", func() {
		It("should implement Service interface", func() {
			var impl types.Service = service
			Expect(impl).ToNot(BeNil())
		})
	})
	Describe("creating a config", func() {
		When("given an url and a message", func() {
			It("should return an error if no arguments where supplied", func() {
				serviceURL, _ := url.Parse("discord://")
				err := service.Initialize(serviceURL, nil)
				Expect(err).To(HaveOccurred())
			})
			It("should not return an error if exactly two arguments are given", func() {
				serviceURL, _ := url.Parse("discord://dummyToken@dummyChannel")
				err := service.Initialize(serviceURL, nil)
				Expect(err).NotTo(HaveOccurred())
			})
			It("should not return an error when given the raw path parameter", func() {
				serviceURL, _ := url.Parse("discord://dummyToken@dummyChannel/raw")
				err := service.Initialize(serviceURL, nil)
				Expect(err).NotTo(HaveOccurred())
			})
			It("should set the JSON flag when given the raw path parameter", func() {
				serviceURL, _ := url.Parse("discord://dummyToken@dummyChannel/raw")
				config := Config{}
				err := config.SetURL(serviceURL)
				Expect(err).NotTo(HaveOccurred())
				Expect(config.JSON).To(BeTrue())
			})
			It("should not set the JSON flag when not provided raw path parameter", func() {
				serviceURL, _ := url.Parse("discord://dummyToken@dummyChannel")
				config := Config{}
				err := config.SetURL(serviceURL)
				Expect(err).NotTo(HaveOccurred())
				Expect(config.JSON).NotTo(BeTrue())
			})
			It("should return an error if more than two arguments are given", func() {
				serviceURL, _ := url.Parse("discord://dummyToken@dummyChannel/illegal-argument")
				err := service.Initialize(serviceURL, nil)
				Expect(err).To(HaveOccurred())
			})
		})
		When("parsing the configuration URL", func() {
			It("should be identical after de-/serialization", func() {
				testURL := "discord://token@channel?avatar=TestBot.jpg&color=112233&colordebug=223344&colorerror=334455&colorinfo=445566&colorwarn=556677&splitlines=No&title=Test Title&username=TestBot"

				url, err := url.Parse(testURL)
				Expect(err).NotTo(HaveOccurred(), "parsing")

				config := &Config{}
				err = config.SetURL(url)
				Expect(err).NotTo(HaveOccurred(), "verifying")

				outputURL := config.GetURL()

				Expect(outputURL.String()).To(Equal(testURL))

			})
		})
	})
	Describe("creating a json payload", func() {
		//When("given a blank message", func() {
		//	It("should return an error", func() {
		//		_, err := CreatePayloadFromItems("", false)
		//		Expect(err).To(HaveOccurred())
		//	})
		//})
		When("given a message that exceeds the max length", func() {
			It("should return a payload with chunked messages", func() {

				payload, err := buildPayloadFromHundreds(42, false, "", dummyColors)
				Expect(err).ToNot(HaveOccurred())

				meta := payload.Embeds[0]
				items := payload.Embeds[1:]

				Expect(items).To(HaveLen(3))

				Expect(items[0].Content).To(HaveLen(1994))
				Expect(items[1].Content).To(HaveLen(1999))
				Expect(items[2].Content).To(HaveLen(205))

				Expect(meta.Footer).To(BeNil())
			})
			It("omit characters above total max", func() {

				payload, err := buildPayloadFromHundreds(62, false, "", dummyColors)
				Expect(err).ToNot(HaveOccurred())

				meta := payload.Embeds[0]
				items := payload.Embeds[1:]

				Expect(items).To(HaveLen(4))
				Expect(items[0].Content).To(HaveLen(1994))
				Expect(items[1].Content).To(HaveLen(1999))
				Expect(len(items[2].Content)).To(Equal(1999))
				Expect(len(items[3].Content)).To(Equal(5))

				Expect(meta.Footer.Text).To(ContainSubstring("200"))
			})

			It("rich test 1", func() {

				testTime, _ := time.Parse(time.RFC3339, time.RFC3339)
				items := []types.MessageItem{
					{
						Text:      "Message",
						Timestamp: &testTime,
						Level:     types.Warning,
					},
				}
				payload, err := CreatePayloadFromItems(items, "Title", dummyColors, 0)
				Expect(err).ToNot(HaveOccurred())

				meta := payload.Embeds[0]
				item := payload.Embeds[1]

				Expect(payload.Embeds).To(HaveLen(2))
				Expect(item.Footer.Text).To(Equal(types.Warning.String()))
				Expect(meta.Title).To(Equal("Title"))
				Expect(item.Color).To(Equal(dummyColors[types.Warning]))
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
			config := Config{
				Channel: "1",
				Token:   "dummyToken",
			}
			serviceURL := config.GetURL()
			service := Service{}
			err = service.Initialize(serviceURL, logger)
			Expect(err).NotTo(HaveOccurred())

			setupResponder(&config, 204, "")

			err = service.Send("Message", nil)
			Expect(err).NotTo(HaveOccurred())
		})

	})
})

func buildPayloadFromHundreds(hundreds int, split bool, title string, colors [types.MessageLevelCount]int) (WebhookPayload, error) {
	hundredChars := "this string is exactly (to the letter) a hundred characters long which will make the send func error"
	builder := strings.Builder{}

	for i := 0; i < hundreds; i++ {
		builder.WriteString(hundredChars)
	}

	items, omitted := CreateItemsFromPlain(builder.String(), split)

	return CreatePayloadFromItems(items, title, colors, omitted)
}

func setupResponder(config *Config, code int, body string) {
	targetURL := CreateAPIURLFromConfig(config)
	httpmock.RegisterResponder("POST", targetURL, httpmock.NewStringResponder(code, body))
}
