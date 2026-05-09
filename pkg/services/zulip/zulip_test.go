package zulip_test

import (
	"errors"
	"github.com/containrrr/shoutrrr/internal/testutils"
	"github.com/containrrr/shoutrrr/pkg/services/zulip"
	. "github.com/containrrr/shoutrrr/pkg/services/zulip"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/jarcoal/httpmock"
	"io"
	"net/http"

	"net/url"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
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

var _ = BeforeSuite(func() {
	service = &Service{}
	envZulipURL, _ = url.Parse(os.Getenv("SHOUTRRR_ZULIP_URL"))
})

var _ = Describe("the zulip service", func() {

	When("running integration tests", func() {
		It("should not error out", func() {
			if envZulipURL.String() == "" {
				return
			}

			serviceURL, _ := url.Parse(envZulipURL.String())
			err := service.Initialize(serviceURL, testutils.TestLogger())
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
			err = service.Initialize(zulipURL, testutils.TestLogger())
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
				url := config.GetURL()
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
				url := config.GetURL()
				Expect(url.String()).To(Equal("zulip://bot-name%40zulipchat.com:correcthorsebatterystable@example.zulipchat.com?stream=foo"))
			})
		})
		When("given a service url with a non-standard port", func() {
			It("should preserve the port", func() {
				zulipURL, err := url.Parse("zulip://bot-name%40zulipchat.com:correcthorsebatterystable@example.zulipchat.com:8443?stream=foo&topic=bar")
				Expect(err).NotTo(HaveOccurred())
				serviceConfig, err := CreateConfigFromURL(zulipURL)
				Expect(err).NotTo(HaveOccurred())

				Expect(serviceConfig.Host).To(Equal("example.zulipchat.com:8443"))
				Expect(serviceConfig.GetURL().String()).To(Equal("zulip://bot-name%40zulipchat.com:correcthorsebatterystable@example.zulipchat.com:8443?stream=foo&topic=bar"))
			})
		})
	})
	Describe("sending messages", func() {
		const apiURL = "https://bot-name%40zulipchat.com:correcthorsebatterystable@example.zulipchat.com/api/v1/messages"

		BeforeEach(func() {
			httpmock.Activate()
		})

		AfterEach(func() {
			httpmock.DeactivateAndReset()
		})

		When("the topic is too long", func() {
			It("should return an error before posting", func() {
				zulipURL, err := url.Parse("zulip://bot-name%40zulipchat.com:correcthorsebatterystable@example.zulipchat.com?stream=foo&topic=abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghi")
				Expect(err).NotTo(HaveOccurred())

				service := &Service{}
				err = service.Initialize(zulipURL, testutils.TestLogger())
				Expect(err).NotTo(HaveOccurred())

				err = service.Send("This is a message", nil)
				Expect(err).To(MatchError("topic exceeds max length (60 characters): was 61 characters"))
				Expect(httpmock.GetTotalCallCount()).To(Equal(0))
			})
		})

		When("the message is too large", func() {
			It("should return an error before posting", func() {
				zulipURL, err := url.Parse("zulip://bot-name%40zulipchat.com:correcthorsebatterystable@example.zulipchat.com?stream=foo&topic=bar")
				Expect(err).NotTo(HaveOccurred())

				service := &Service{}
				err = service.Initialize(zulipURL, testutils.TestLogger())
				Expect(err).NotTo(HaveOccurred())

				err = service.Send(string(make([]byte, 10001)), nil)
				Expect(err).To(MatchError("message exceeds max size (10000 bytes): was 10001 bytes"))
				Expect(httpmock.GetTotalCallCount()).To(Equal(0))
			})
		})

		When("send-time params override stream and topic", func() {
			It("should post the overridden form values", func() {
				zulipURL, err := url.Parse("zulip://bot-name%40zulipchat.com:correcthorsebatterystable@example.zulipchat.com?stream=foo&topic=bar")
				Expect(err).NotTo(HaveOccurred())

				service := &Service{}
				err = service.Initialize(zulipURL, testutils.TestLogger())
				Expect(err).NotTo(HaveOccurred())

				httpmock.RegisterResponder(http.MethodPost, apiURL, func(req *http.Request) (*http.Response, error) {
					body, err := io.ReadAll(req.Body)
					Expect(err).NotTo(HaveOccurred())
					Expect(req.Header.Get("Content-Type")).To(Equal("application/x-www-form-urlencoded"))

					form, err := url.ParseQuery(string(body))
					Expect(err).NotTo(HaveOccurred())
					Expect(form.Get("type")).To(Equal("stream"))
					Expect(form.Get("to")).To(Equal("overridden-stream"))
					Expect(form.Get("topic")).To(Equal("overridden-topic"))
					Expect(form.Get("content")).To(Equal("This is a message"))

					return httpmock.NewStringResponse(http.StatusOK, ""), nil
				})

				params := types.Params{"stream": "overridden-stream", "topic": "overridden-topic"}
				err = service.Send("This is a message", &params)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("the Zulip API rejects the notification", func() {
			It("should report the response status", func() {
				zulipURL, err := url.Parse("zulip://bot-name%40zulipchat.com:correcthorsebatterystable@example.zulipchat.com?stream=foo&topic=bar")
				Expect(err).NotTo(HaveOccurred())

				service := &Service{}
				err = service.Initialize(zulipURL, testutils.TestLogger())
				Expect(err).NotTo(HaveOccurred())

				httpmock.RegisterResponder(http.MethodPost, apiURL, httpmock.NewStringResponder(http.StatusBadRequest, "bad payload"))

				err = service.Send("This is a message", nil)
				Expect(err).To(MatchError("failed to send zulip message: response status code 400 Bad Request"))
			})
		})

		When("the Zulip API cannot be reached", func() {
			It("should report the transport error", func() {
				zulipURL, err := url.Parse("zulip://bot-name%40zulipchat.com:correcthorsebatterystable@example.zulipchat.com?stream=foo&topic=bar")
				Expect(err).NotTo(HaveOccurred())

				service := &Service{}
				err = service.Initialize(zulipURL, testutils.TestLogger())
				Expect(err).NotTo(HaveOccurred())

				httpmock.RegisterResponder(http.MethodPost, apiURL, httpmock.NewErrorResponder(errors.New("network down")))

				err = service.Send("This is a message", nil)
				Expect(err).To(MatchError(ContainSubstring("failed to send zulip message")))
				Expect(err).To(MatchError(ContainSubstring("network down")))
			})
		})
	})
})

func expectErrorMessageGivenURL(msg ErrorMessage, zulipURL *url.URL) {
	err := service.Initialize(zulipURL, testutils.TestLogger())
	Expect(err).To(HaveOccurred())
	Expect(err.Error()).To(Equal(string(msg)))
}
