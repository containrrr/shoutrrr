package zulip

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/nicholas-fedor/shoutrrr/internal/testutils"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestZulip(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Shoutrrr Zulip Suite")
}

var (
	service     *Service
	envZulipURL *url.URL
)

var _ = ginkgo.BeforeSuite(func() {
	service = &Service{}
	envZulipURL, _ = url.Parse(os.Getenv("SHOUTRRR_ZULIP_URL"))
})

// Helper function to create Zulip URLs with optional overrides.
func createZulipURL(botMail, botKey, host, stream, topic string) *url.URL {
	query := url.Values{}
	if stream != "" {
		query.Set("stream", stream)
	}

	if topic != "" {
		query.Set("topic", topic)
	}

	u := &url.URL{
		Scheme:   "zulip",
		User:     url.UserPassword(botMail, botKey),
		Host:     host,
		RawQuery: query.Encode(),
	}

	return u
}

var _ = ginkgo.Describe("the zulip service", func() {
	ginkgo.When("running integration tests", func() {
		ginkgo.It("should not error out", func() {
			// Covers zulip.go:32-52 (Send), zulip.go:55-63 (Initialize), zulip.go:65-78 (doSend partially)
			if envZulipURL.String() == "" {
				return
			}
			serviceURL, _ := url.Parse(envZulipURL.String())
			err := service.Initialize(serviceURL, testutils.TestLogger())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			err = service.Send("This is an integration test message", nil)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
	})

	ginkgo.When("given a service url with missing parts", func() {
		ginkgo.It("should return an error if bot mail is missing", func() {
			// Covers zulip_config.go:55-58 (setURL BotMail check), zulip_errors.go:15 (MissingBotMail)
			zulipURL := createZulipURL("", "correcthorsebatterystable", "example.zulipchat.com", "foo", "bar")
			expectErrorMessageGivenURL(MissingBotMail, zulipURL)
		})
		ginkgo.It("should return an error if api key is missing", func() {
			// Covers zulip_config.go:60-63 (setURL BotKey check), zulip_errors.go:13 (MissingAPIKey)
			zulipURL := &url.URL{
				Scheme: "zulip",
				User:   url.User("bot-name@zulipchat.com"),
				Host:   "example.zulipchat.com",
				RawQuery: url.Values{
					"stream": []string{"foo"},
					"topic":  []string{"bar"},
				}.Encode(),
			}
			expectErrorMessageGivenURL(MissingAPIKey, zulipURL)
		})
		ginkgo.It("should return an error if host is missing", func() {
			// Covers zulip_config.go:65-68 (setURL Host check), zulip_errors.go:14 (MissingHost)
			zulipURL := createZulipURL("bot-name@zulipchat.com", "correcthorsebatterystable", "", "foo", "bar")
			expectErrorMessageGivenURL(MissingHost, zulipURL)
		})
	})
	ginkgo.When("given a valid service url is provided", func() {
		ginkgo.It("should not return an error", func() {
			// Covers zulip.go:55-63 (Initialize), zulip_config.go:49-72 (SetURL success case)
			zulipURL := createZulipURL("bot-name@zulipchat.com", "correcthorsebatterystable", "example.zulipchat.com", "foo", "bar")
			err := service.Initialize(zulipURL, testutils.TestLogger())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
	})
	ginkgo.When("sending a message", func() {
		ginkgo.It("should error if topic exceeds max length", func() {
			// Covers zulip.go:43-45 (topic length validation), zulip_errors.go:16 (TopicTooLong)
			zulipURL := createZulipURL("bot-name@zulipchat.com", "correcthorsebatterystable", "example.zulipchat.com", "foo", "")
			err := service.Initialize(zulipURL, testutils.TestLogger())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			longTopic := strings.Repeat("a", topicMaxLength+1) // 61 chars
			params := &types.Params{"topic": longTopic}
			err = service.Send("test message", params)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring(fmt.Sprintf(string(TopicTooLong), topicMaxLength, len([]rune(longTopic)))))
		})
		ginkgo.It("should error if message exceeds max size", func() {
			// Covers zulip.go:47-50 (message size validation)
			zulipURL := createZulipURL("bot-name@zulipchat.com", "correcthorsebatterystable", "example.zulipchat.com", "foo", "bar")
			err := service.Initialize(zulipURL, testutils.TestLogger())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			longMessage := strings.Repeat("a", contentMaxSize+1) // 10001 bytes
			err = service.Send(longMessage, nil)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(err.Error()).To(gomega.Equal(fmt.Sprintf("message exceeds max size (%d bytes): was %d bytes", contentMaxSize, len(longMessage))))
		})
		ginkgo.It("should override stream from params", func() {
			// Covers zulip.go:38-39 (Stream override from params)
			zulipURL := createZulipURL("bot-name@zulipchat.com", "correcthorsebatterystable", "example.zulipchat.com", "original", "")
			err := service.Initialize(zulipURL, testutils.TestLogger())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			params := &types.Params{"stream": "newstream"}
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			apiURL := service.getAPIURL(&Config{
				BotMail: "bot-name@zulipchat.com",
				BotKey:  "correcthorsebatterystable",
				Host:    "example.zulipchat.com",
				Stream:  "newstream",
			})
			httpmock.RegisterResponder("POST", apiURL, httpmock.NewStringResponder(http.StatusOK, ""))
			err = service.Send("test message", params)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
		ginkgo.It("should override topic from params", func() {
			// Covers zulip.go:41-42 (Topic override from params)
			zulipURL := createZulipURL("bot-name@zulipchat.com", "correcthorsebatterystable", "example.zulipchat.com", "foo", "original")
			err := service.Initialize(zulipURL, testutils.TestLogger())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			params := &types.Params{"topic": "newtopic"}
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			config := &Config{
				BotMail: "bot-name@zulipchat.com",
				BotKey:  "correcthorsebatterystable",
				Host:    "example.zulipchat.com",
				Stream:  "foo",
				Topic:   "newtopic",
			}
			apiURL := service.getAPIURL(config)
			httpmock.RegisterResponder("POST", apiURL, func(req *http.Request) (*http.Response, error) {
				gomega.Expect(req.FormValue("topic")).To(gomega.Equal("newtopic"))

				return httpmock.NewStringResponse(http.StatusOK, ""), nil
			})
			err = service.Send("test message", params)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
		ginkgo.It("should handle HTTP errors", func() {
			// Covers zulip.go:71-75 (doSend HTTP error handling)
			zulipURL := createZulipURL("bot-name@zulipchat.com", "correcthorsebatterystable", "example.zulipchat.com", "foo", "bar")
			err := service.Initialize(zulipURL, testutils.TestLogger())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			apiURL := service.getAPIURL(service.Config)
			httpmock.RegisterResponder("POST", apiURL, httpmock.NewStringResponder(http.StatusBadRequest, "Bad Request"))
			err = service.Send("test message", nil)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("response status code 400"))
		})
	})
	ginkgo.Describe("the zulip config", func() {
		ginkgo.When("cloning a config object", func() {
			ginkgo.It("the clone should have equal values", func() {
				// Covers zulip_config.go:75-84 (Clone equality)
				config1 := &Config{
					BotMail: "bot-name@zulipchat.com",
					BotKey:  "correcthorsebatterystable",
					Host:    "example.zulipchat.com",
					Stream:  "foo",
					Topic:   "bar",
				}
				config2 := config1.Clone()
				gomega.Expect(config1).To(gomega.Equal(config2))
			})
			ginkgo.It("the clone should not be the same struct", func() {
				// Covers zulip_config.go:75-84 (Clone identity)
				config1 := &Config{
					BotMail: "bot-name@zulipchat.com",
					BotKey:  "correcthorsebatterystable",
					Host:    "example.zulipchat.com",
					Stream:  "foo",
					Topic:   "bar",
				}
				config2 := config1.Clone()
				gomega.Expect(config1).NotTo(gomega.BeIdenticalTo(config2))
			})
		})
		ginkgo.When("generating a config object", func() {
			ginkgo.It("should generate a correct config object using CreateConfigFromURL", func() {
				// Covers zulip_config.go:92-98 (CreateConfigFromURL), zulip_config.go:49-72 (setURL)
				zulipURL := createZulipURL("bot-name@zulipchat.com", "correcthorsebatterystable", "example.zulipchat.com", "foo", "bar")
				serviceConfig, err := CreateConfigFromURL(zulipURL)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				config := &Config{
					BotMail: "bot-name@zulipchat.com",
					BotKey:  "correcthorsebatterystable",
					Host:    "example.zulipchat.com",
					Stream:  "foo",
					Topic:   "bar",
				}
				gomega.Expect(serviceConfig).To(gomega.Equal(config))
			})
			ginkgo.It("should update config correctly using SetURL", func() {
				// Covers zulip_config.go:27-29 (SetURL), zulip_config.go:49-72 (setURL)
				config := &Config{} // Start with empty config
				zulipURL := createZulipURL("bot-name@zulipchat.com", "correcthorsebatterystable", "example.zulipchat.com", "foo", "bar")
				err := config.SetURL(zulipURL)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				expected := &Config{
					BotMail: "bot-name@zulipchat.com",
					BotKey:  "correcthorsebatterystable",
					Host:    "example.zulipchat.com",
					Stream:  "foo",
					Topic:   "bar",
				}
				gomega.Expect(config).To(gomega.Equal(expected))
			})
		})
		ginkgo.When("given a config object with stream and topic", func() {
			ginkgo.It("should build the correct service url", func() {
				// Covers zulip_config.go:27-46 (GetURL with Topic)
				config := Config{
					BotMail: "bot-name@zulipchat.com",
					BotKey:  "correcthorsebatterystable",
					Host:    "example.zulipchat.com",
					Stream:  "foo",
					Topic:   "bar",
				}
				url := config.GetURL()
				gomega.Expect(url.String()).To(gomega.Equal("zulip://bot-name%40zulipchat.com:correcthorsebatterystable@example.zulipchat.com?stream=foo&topic=bar"))
			})
		})
		ginkgo.When("given a config object with stream but without topic", func() {
			ginkgo.It("should build the correct service url", func() {
				// Covers zulip_config.go:27-46 (GetURL without Topic)
				config := Config{
					BotMail: "bot-name@zulipchat.com",
					BotKey:  "correcthorsebatterystable",
					Host:    "example.zulipchat.com",
					Stream:  "foo",
				}
				url := config.GetURL()
				gomega.Expect(url.String()).To(gomega.Equal("zulip://bot-name%40zulipchat.com:correcthorsebatterystable@example.zulipchat.com?stream=foo"))
			})
		})
	})
	ginkgo.Describe("the zulip payload", func() {
		ginkgo.When("creating a payload with topic", func() {
			ginkgo.It("should include all fields", func() {
				// Covers zulip_payload.go:7-18 (CreatePayload with Topic)
				config := &Config{
					Stream: "foo",
					Topic:  "bar",
				}
				payload := CreatePayload(config, "test message")
				gomega.Expect(payload.Get("type")).To(gomega.Equal("stream"))
				gomega.Expect(payload.Get("to")).To(gomega.Equal("foo"))
				gomega.Expect(payload.Get("content")).To(gomega.Equal("test message"))
				gomega.Expect(payload.Get("topic")).To(gomega.Equal("bar"))
			})
		})
		ginkgo.When("creating a payload without topic", func() {
			ginkgo.It("should exclude topic field", func() {
				// Covers zulip_payload.go:7-18 (CreatePayload without Topic)
				config := &Config{
					Stream: "foo",
				}
				payload := CreatePayload(config, "test message")
				gomega.Expect(payload.Get("type")).To(gomega.Equal("stream"))
				gomega.Expect(payload.Get("to")).To(gomega.Equal("foo"))
				gomega.Expect(payload.Get("content")).To(gomega.Equal("test message"))
				gomega.Expect(payload.Get("topic")).To(gomega.Equal(""))
			})
		})
	})
	ginkgo.It("should return the correct service ID", func() {
		service := &Service{}
		gomega.Expect(service.GetID()).To(gomega.Equal("zulip"))
	})
})

func expectErrorMessageGivenURL(msg ErrorMessage, zulipURL *url.URL) {
	err := service.Initialize(zulipURL, testutils.TestLogger())
	gomega.Expect(err).To(gomega.HaveOccurred())
	gomega.Expect(err.Error()).To(gomega.Equal(string(msg)))
}
