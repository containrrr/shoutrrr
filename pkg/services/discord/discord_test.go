package discord_test

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/nicholas-fedor/shoutrrr/internal/testutils"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/discord"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestDiscord(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Shoutrrr Discord Suite")
}

var (
	dummyColors   = [types.MessageLevelCount]uint{}
	service       *discord.Service
	envDiscordURL *url.URL
	logger        *log.Logger
	_             = ginkgo.BeforeSuite(func() {
		service = &discord.Service{}
		envDiscordURL, _ = url.Parse(os.Getenv("SHOUTRRR_DISCORD_URL"))
		logger = log.New(ginkgo.GinkgoWriter, "Test", log.LstdFlags)
	})
)

var _ = ginkgo.Describe("the discord service", func() {
	ginkgo.When("running integration tests", func() {
		ginkgo.It("should work without errors", func() {
			if envDiscordURL.String() == "" {
				return
			}

			serviceURL, _ := url.Parse(envDiscordURL.String())
			err := service.Initialize(serviceURL, testutils.TestLogger())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			err = service.Send(
				"this is an integration test",
				nil,
			)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
	})
	ginkgo.Describe("the service", func() {
		ginkgo.It("should implement Service interface", func() {
			var impl types.Service = service
			gomega.Expect(impl).ToNot(gomega.BeNil())
		})
		ginkgo.It("returns the correct service identifier", func() {
			// No initialization needed since GetID is static
			gomega.Expect(service.GetID()).To(gomega.Equal("discord"))
		})
	})
	ginkgo.Describe("creating a config", func() {
		ginkgo.When("given an url and a message", func() {
			ginkgo.It("should return an error if no arguments where supplied", func() {
				serviceURL, _ := url.Parse("discord://")
				err := service.Initialize(serviceURL, nil)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
			ginkgo.It("should not return an error if exactly two arguments are given", func() {
				serviceURL, _ := url.Parse("discord://dummyToken@dummyChannel")
				err := service.Initialize(serviceURL, nil)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
			ginkgo.It("should not return an error when given the raw path parameter", func() {
				serviceURL, _ := url.Parse("discord://dummyToken@dummyChannel/raw")
				err := service.Initialize(serviceURL, nil)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
			ginkgo.It("should set the JSON flag when given the raw path parameter", func() {
				serviceURL, _ := url.Parse("discord://dummyToken@dummyChannel/raw")
				config := discord.Config{}
				err := config.SetURL(serviceURL)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(config.JSON).To(gomega.BeTrue())
			})
			ginkgo.It("should not set the JSON flag when not provided raw path parameter", func() {
				serviceURL, _ := url.Parse("discord://dummyToken@dummyChannel")
				config := discord.Config{}
				err := config.SetURL(serviceURL)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(config.JSON).NotTo(gomega.BeTrue())
			})
			ginkgo.It("should return an error if more than two arguments are given", func() {
				serviceURL, _ := url.Parse("discord://dummyToken@dummyChannel/illegal-argument")
				err := service.Initialize(serviceURL, nil)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
		ginkgo.When("parsing the configuration URL", func() {
			ginkgo.It("should be identical after de-/serialization", func() {
				testURL := "discord://token@channel?avatar=TestBot.jpg&color=0x112233&colordebug=0x223344&colorerror=0x334455&colorinfo=0x445566&colorwarn=0x556677&splitlines=No&title=Test+Title&username=TestBot"

				url, err := url.Parse(testURL)
				gomega.Expect(err).NotTo(gomega.HaveOccurred(), "parsing")

				config := &discord.Config{}
				err = config.SetURL(url)
				gomega.Expect(err).NotTo(gomega.HaveOccurred(), "verifying")

				outputURL := config.GetURL()

				gomega.Expect(outputURL.String()).To(gomega.Equal(testURL))
			})
		})
	})
	ginkgo.Describe("creating a json payload", func() {
		ginkgo.When("given a blank message", func() {
			ginkgo.When("split lines is enabled", func() {
				ginkgo.It("should return an error", func() {
					// batches := CreateItemsFromPlain("", true)
					items := []types.MessageItem{}
					gomega.Expect(items).To(gomega.BeEmpty())
					_, err := discord.CreatePayloadFromItems(items, "title", dummyColors)
					gomega.Expect(err).To(gomega.HaveOccurred())
				})
			})
			ginkgo.When("split lines is disabled", func() {
				ginkgo.It("should return an error", func() {
					batches := discord.CreateItemsFromPlain("", false)
					items := batches[0]
					gomega.Expect(items).To(gomega.BeEmpty())
					_, err := discord.CreatePayloadFromItems(items, "title", dummyColors)
					gomega.Expect(err).To(gomega.HaveOccurred())
				})
			})
		})
		ginkgo.When("given a message that exceeds the max length", func() {
			ginkgo.It("should return a payload with chunked messages", func() {
				payload, err := buildPayloadFromHundreds(42, false, "Title", dummyColors)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())

				items := payload.Embeds

				gomega.Expect(items).To(gomega.HaveLen(3))

				gomega.Expect(items[0].Content).To(gomega.HaveLen(1994))
				gomega.Expect(items[1].Content).To(gomega.HaveLen(1999))
				gomega.Expect(items[2].Content).To(gomega.HaveLen(205))
			})
			ginkgo.It("omit characters above total max", func() {
				payload, err := buildPayloadFromHundreds(62, false, "", dummyColors)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())

				items := payload.Embeds

				gomega.Expect(items).To(gomega.HaveLen(4))
				gomega.Expect(items[0].Content).To(gomega.HaveLen(1994))
				gomega.Expect(items[1].Content).To(gomega.HaveLen(1999))
				gomega.Expect(len(items[2].Content)).To(gomega.Equal(1999))
				gomega.Expect(len(items[3].Content)).To(gomega.Equal(5))

				// gomega.Expect(meta.Footer.Text).To(ContainSubstring("200"))
			})
			ginkgo.When("no title is supplied and content fits", func() {
				ginkgo.It("should return a payload without a meta chunk", func() {
					payload, err := buildPayloadFromHundreds(42, false, "", dummyColors)
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(payload.Embeds[0].Footer).To(gomega.BeNil())
					gomega.Expect(payload.Embeds[0].Title).To(gomega.BeEmpty())
				})
			})
			ginkgo.When("title is supplied, but content fits", func() {
				ginkgo.It("should return a payload with a meta chunk", func() {
					payload, err := buildPayloadFromHundreds(42, false, "Title", dummyColors)
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(payload.Embeds[0].Title).ToNot(gomega.BeEmpty())
				})
			})

			ginkgo.It("rich test 1", func() {
				testTime, _ := time.Parse(time.RFC3339, time.RFC3339)
				items := []types.MessageItem{
					{
						Text:      "Message",
						Timestamp: testTime,
						Level:     types.Warning,
					},
				}
				payload, err := discord.CreatePayloadFromItems(items, "Title", dummyColors)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())

				item := payload.Embeds[0]

				gomega.Expect(payload.Embeds).To(gomega.HaveLen(1))
				gomega.Expect(item.Footer.Text).To(gomega.Equal(types.Warning.String()))
				gomega.Expect(item.Title).To(gomega.Equal("Title"))
				gomega.Expect(item.Color).To(gomega.Equal(dummyColors[types.Warning]))
			})
		})
	})

	ginkgo.Describe("sending the payload", func() {
		dummyConfig := discord.Config{
			WebhookID: "1",
			Token:     "dummyToken",
		}
		var service discord.Service
		ginkgo.BeforeEach(func() {
			httpmock.Activate()
			service = discord.Service{}
			if err := service.Initialize(dummyConfig.GetURL(), logger); err != nil {
				panic(fmt.Errorf("service initialization failed: %w", err))
			}
		})
		ginkgo.AfterEach(func() {
			httpmock.DeactivateAndReset()
		})
		ginkgo.It("should not report an error if the server accepts the payload", func() {
			setupResponder(&dummyConfig, 204, "")

			gomega.Expect(service.Send("Message", nil)).To(gomega.Succeed())
		})
		ginkgo.It("should report an error if the server response is not OK", func() {
			setupResponder(&dummyConfig, 400, "")
			gomega.Expect(service.Initialize(dummyConfig.GetURL(), logger)).To(gomega.Succeed())
			gomega.Expect(service.Send("Message", nil)).NotTo(gomega.Succeed())
		})
		ginkgo.It("should report an error if the message is empty", func() {
			setupResponder(&dummyConfig, 204, "")
			gomega.Expect(service.Initialize(dummyConfig.GetURL(), logger)).To(gomega.Succeed())
			gomega.Expect(service.Send("", nil)).NotTo(gomega.Succeed())
		})
		ginkgo.When("using a custom json payload", func() {
			ginkgo.It("should report an error if the server response is not OK", func() {
				config := dummyConfig
				config.JSON = true
				setupResponder(&config, 400, "")
				gomega.Expect(service.Initialize(config.GetURL(), logger)).To(gomega.Succeed())
				gomega.Expect(service.Send("Message", nil)).NotTo(gomega.Succeed())
			})
		})
	})
})

func buildPayloadFromHundreds(hundreds int, split bool, title string, colors [types.MessageLevelCount]uint) (discord.WebhookPayload, error) {
	hundredChars := "this string is exactly (to the letter) a hundred characters long which will make the send func error"
	builder := strings.Builder{}

	for i := 0; i < hundreds; i++ {
		builder.WriteString(hundredChars)
	}

	batches := discord.CreateItemsFromPlain(builder.String(), split)
	items := batches[0]

	return discord.CreatePayloadFromItems(items, title, colors)
}

func setupResponder(config *discord.Config, code int, body string) {
	targetURL := discord.CreateAPIURLFromConfig(config)
	httpmock.RegisterResponder("POST", targetURL, httpmock.NewStringResponder(code, body))
}
