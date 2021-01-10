package discord_test

import (
	. "github.com/containrrr/shoutrrr/pkg/services/discord"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util"

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

				items, omitted := buildItemsFromHundreds(42, false)

				Expect(items).To(HaveLen(3))

				Expect(items[0].Text).To(HaveLen(1994))
				Expect(items[1].Text).To(HaveLen(1999))
				Expect(items[2].Text).To(HaveLen(205))

				Expect(omitted).To(Equal(0))
			})
			It("omit characters above total max", func() {

				items, omitted := buildItemsFromHundreds(62, false)

				Expect(items).To(HaveLen(4))

				Expect(items[0].Text).To(HaveLen(1994))
				Expect(items[1].Text).To(HaveLen(1999))
				Expect(len(items[2].Text)).To(Equal(1999))
				Expect(len(items[3].Text)).To(Equal(5))

				Expect(omitted).To(Equal(200))
			})
		})
	})
})

func buildItemsFromHundreds(hundreds int, split bool) (items []types.MessageItem, omitted int) {
	hundredChars := "this string is exactly (to the letter) a hundred characters long which will make the send func error"
	builder := strings.Builder{}

	for i := 0; i < hundreds; i++ {
		builder.WriteString(hundredChars)
	}

	return CreateItemsFromPlain(builder.String(), split)
}
