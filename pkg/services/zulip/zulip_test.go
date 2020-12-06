package zulip_test

import (
	"github.com/containrrr/shoutrrr/pkg/services/zulip"
	. "github.com/containrrr/shoutrrr/pkg/services/zulip"
	"github.com/containrrr/shoutrrr/pkg/util"

	"net/url"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestZulip(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Zulip Suite")
}

var (
	service     *Service
	envZulipURL *url.URL
)

var _ = Describe("the zulip service", func() {

	BeforeSuite(func() {
		service = &Service{}
		envZulipURL, _ = url.Parse(os.Getenv("SHOUTRRR_ZULIP_URL"))

	})

	When("running integration tests", func() {
		It("should not error out", func() {
			if envZulipURL.String() == "" {
				return
			}

			serviceURL, _ := url.Parse(envZulipURL.String())
			err := service.Initialize(serviceURL, util.TestLogger())
			Expect(err).NotTo(HaveOccurred())
			err = service.Send("This is an integration test message", nil)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	When("given a service url with missing parts", func() {
		It("should return an error if bot mail is missing", func() {
			zulipURL, err := url.Parse("zulip://example.zulipchat.com?stream=foo&topic=bar")
			Expect(err).NotTo(HaveOccurred())
			expectErrorMessageGivenURL(
				MissingBotMail,
				zulipURL,
			)
		})
		It("should return an error if api key is missing", func() {
			zulipURL, err := url.Parse("zulip://bot-name%40zulipchat.com@example.zulipchat.com?stream=foo&topic=bar")
			Expect(err).NotTo(HaveOccurred())
			expectErrorMessageGivenURL(
				MissingAPIKey,
				zulipURL,
			)
		})
		It("should return an error if host is missing", func() {
			zulipURL, err := url.Parse("zulip://bot-name%40zulipchat.com:correcthorsebatterystable@?stream=foo&topic=bar")
			Expect(err).NotTo(HaveOccurred())
			expectErrorMessageGivenURL(
				MissingHost,
				zulipURL,
			)
		})
	})
	When("given a valid service url is provided", func() {
		It("should not return an error", func() {
			zulipURL, err := url.Parse("zulip://bot-name%40zulipchat.com:correcthorsebatterystable@example.zulipchat.com?stream=foo&topic=bar")
			Expect(err).NotTo(HaveOccurred())
			err = service.Initialize(zulipURL, util.TestLogger())
			Expect(err).NotTo(HaveOccurred())
		})
	})
	Describe("the zulip config", func() {
		When("cloning a config object", func() {
			It("the clone should have equal values", func() {
				config1 := &zulip.Config{
					BotMail: "bot-name@zulipchat.com",
					BotKey:  "correcthorsebatterystable",
					Host:    "example.zulipchat.com",
					Path:    "api/v1/messages",
					Stream:  "foo",
					Topic:   "bar",
				}

				config2 := config1.Clone()

				Expect(config1).To(Equal(config2))
			})
			It("the clone should not be the same struct", func() {
				config1 := &zulip.Config{
					BotMail: "bot-name@zulipchat.com",
					BotKey:  "correcthorsebatterystable",
					Host:    "example.zulipchat.com",
					Path:    "api/v1/messages",
					Stream:  "foo",
					Topic:   "bar",
				}

				config2 := config1.Clone()

				Expect(config1).NotTo(BeIdenticalTo(config2))
			})
		})
		When("generating a config object", func() {
			It("should generate a correct config object", func() {
				zulipURL, err := url.Parse("zulip://bot-name%40zulipchat.com:correcthorsebatterystable@example.zulipchat.com?stream=foo&topic=bar")
				Expect(err).NotTo(HaveOccurred())
				serviceConfig, err := CreateConfigFromURL(zulipURL)
				Expect(err).NotTo(HaveOccurred())

				config := &zulip.Config{
					BotMail: "bot-name@zulipchat.com",
					BotKey:  "correcthorsebatterystable",
					Host:    "example.zulipchat.com",
					Path:    "api/v1/messages",
					Stream:  "foo",
					Topic:   "bar",
				}
				Expect(serviceConfig).To(Equal(config))
			})
		})
		When("given a config object with stream and topic", func() {
			It("should build the correct service url", func() {
				config := zulip.Config{
					BotMail: "bot-name@zulipchat.com",
					BotKey:  "correcthorsebatterystable",
					Host:    "example.zulipchat.com",
					Stream:  "foo",
					Topic:   "bar",
				}
				url := config.GetURL(nil)
				Expect(url.String()).To(Equal("zulip://bot-name%40zulipchat.com:correcthorsebatterystable@example.zulipchat.com?stream=foo&topic=bar"))
			})
		})
		When("given a config object with stream but without topic", func() {
			It("should build the correct service url", func() {
				config := zulip.Config{
					BotMail: "bot-name@zulipchat.com",
					BotKey:  "correcthorsebatterystable",
					Host:    "example.zulipchat.com",
					Stream:  "foo",
				}
				url := config.GetURL(nil)
				Expect(url.String()).To(Equal("zulip://bot-name%40zulipchat.com:correcthorsebatterystable@example.zulipchat.com?stream=foo"))
			})
		})
	})
})

func expectErrorMessageGivenURL(msg ErrorMessage, zulipURL *url.URL) {
	err := service.Initialize(zulipURL, util.TestLogger())
	Expect(err).To(HaveOccurred())
	Expect(err.Error()).To(Equal(string(msg)))
}
