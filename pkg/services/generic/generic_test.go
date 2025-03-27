package generic_test

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/nicholas-fedor/shoutrrr/internal/testutils"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/generic"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// Test constants.
const (
	TestWebhookURL = "https://host.tld/webhook" // Default test webhook URL
)

// TestGeneric runs the Ginkgo test suite for the generic package.
func TestGeneric(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Shoutrrr Generic Webhook Suite")
}

var (
	service       *generic.Service
	logger        *log.Logger
	envGenericURL *url.URL
	_             = ginkgo.BeforeSuite(func() {
		service = &generic.Service{}
		logger = log.New(ginkgo.GinkgoWriter, "Test", log.LstdFlags)
		var err error
		envGenericURL, err = url.Parse(os.Getenv("SHOUTRRR_GENERIC_URL"))
		if err != nil {
			envGenericURL = &url.URL{} // Default to empty URL if parsing fails
		}
	})
)

var _ = ginkgo.Describe("the generic service", func() {
	ginkgo.When("running integration tests", func() {
		ginkgo.It("sends a message successfully with a valid ENV URL", func() {
			if envGenericURL.String() == "" {
				ginkgo.Skip("No integration test ENV URL was set")

				return
			}
			serviceURL := testutils.URLMust(envGenericURL.String())
			err := service.Initialize(serviceURL, testutils.TestLogger())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			err = service.Send("This is an integration test message", nil)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
	})

	ginkgo.Describe("the service", func() {
		ginkgo.BeforeEach(func() {
			service = &generic.Service{}
			service.SetLogger(logger)
		})
		ginkgo.It("returns the correct service identifier", func() {
			gomega.Expect(service.GetID()).To(gomega.Equal("generic"))
		})
	})

	ginkgo.When("parsing a custom URL", func() {
		ginkgo.BeforeEach(func() {
			service = &generic.Service{}
			service.SetLogger(logger)
		})
		ginkgo.It("correctly sets webhook URL from custom URL", func() {
			customURL := testutils.URLMust("generic+https://test.tld")
			serviceURL, err := service.GetConfigURLFromCustom(customURL)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			err = service.Initialize(serviceURL, logger)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(service.Config.WebhookURL().String()).To(gomega.Equal("https://test.tld"))
		})

		ginkgo.When("a HTTP URL is provided via query parameter", func() {
			ginkgo.It("disables TLS", func() {
				config := &generic.Config{}
				err := config.SetURL(testutils.URLMust("generic://example.com?disabletls=yes"))
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(config.DisableTLS).To(gomega.BeTrue())
			})
		})
		ginkgo.When("a HTTPS URL is provided", func() {
			ginkgo.It("enables TLS", func() {
				config := &generic.Config{}
				err := config.SetURL(testutils.URLMust("generic://example.com"))
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(config.DisableTLS).To(gomega.BeFalse())
			})
		})
		ginkgo.It("escapes conflicting custom query keys", func() {
			serviceURL := testutils.URLMust("generic://example.com/?__template=passed")
			err := service.Initialize(serviceURL, logger)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(service.Config.Template).NotTo(gomega.Equal("passed"))
			whURL := service.Config.WebhookURL().String()
			gomega.Expect(whURL).To(gomega.Equal("https://example.com/?template=passed"))
			gomega.Expect(service.Config.GetURL().String()).To(gomega.Equal(serviceURL.String()))
		})
		ginkgo.It("handles both escaped and service prop versions of keys", func() {
			serviceURL := testutils.URLMust("generic://example.com/?__template=passed&template=captured")
			err := service.Initialize(serviceURL, logger)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(service.Config.Template).To(gomega.Equal("captured"))
			whURL := service.Config.WebhookURL().String()
			gomega.Expect(whURL).To(gomega.Equal("https://example.com/?template=passed"))
		})
	})

	ginkgo.When("retrieving the webhook URL", func() {
		ginkgo.BeforeEach(func() {
			service = &generic.Service{}
			service.SetLogger(logger)
		})
		ginkgo.It("builds a valid webhook URL", func() {
			serviceURL := testutils.URLMust("generic://example.com/path?foo=bar")
			err := service.Initialize(serviceURL, logger)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(service.Config.WebhookURL().String()).To(gomega.Equal("https://example.com/path?foo=bar"))
		})

		ginkgo.When("TLS is disabled", func() {
			ginkgo.It("uses http scheme", func() {
				serviceURL := testutils.URLMust("generic://test.tld?disabletls=yes")
				err := service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(service.Config.WebhookURL().Scheme).To(gomega.Equal("http"))
			})
		})
		ginkgo.When("TLS is not disabled", func() {
			ginkgo.It("uses https scheme", func() {
				serviceURL := testutils.URLMust("generic://test.tld")
				err := service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(service.Config.WebhookURL().Scheme).To(gomega.Equal("https"))
			})
		})
	})

	ginkgo.Describe("the generic config", func() {
		ginkgo.When("parsing the configuration URL", func() {
			ginkgo.It("is identical after de-/serialization", func() {
				testURL := "generic://user:pass@host.tld/api/v1/webhook?$context=inside-joke&@Authorization=frend&__title=w&contenttype=a%2Fb&template=f&title=t"
				expectedURL := "generic://user:pass@host.tld/api/v1/webhook?%24context=inside-joke&%40Authorization=frend&__title=w&contenttype=a%2Fb&template=f&title=t"
				serviceURL := testutils.URLMust(testURL)
				err := service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(service.Config.GetURL().String()).To(gomega.Equal(expectedURL))
			})
		})
	})

	ginkgo.Describe("building the payload", func() {
		ginkgo.BeforeEach(func() {
			service = &generic.Service{}
			service.SetLogger(logger)
		})
		ginkgo.When("no template is specified", func() {
			ginkgo.It("uses the message as payload", func() {
				serviceURL := testutils.URLMust("generic://host.tld/webhook")
				err := service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				payload, err := service.GetPayload(service.Config, types.Params{"message": "test message"})
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				contents, err := io.ReadAll(payload)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(string(contents)).To(gomega.Equal("test message"))
			})
		})
		ginkgo.When("template is specified as `JSON`", func() {
			ginkgo.It("creates a JSON object as the payload", func() {
				serviceURL := testutils.URLMust("generic://host.tld/webhook?template=JSON")
				err := service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				params := types.Params{"title": "test title", "message": "test message"}
				payload, err := service.GetPayload(service.Config, params)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				contents, err := io.ReadAll(payload)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(string(contents)).To(gomega.MatchJSON(`{
					"title": "test title",
					"message": "test message"
				}`))
			})
			ginkgo.When("alternate keys are specified", func() {
				ginkgo.It("creates a JSON object using the specified keys", func() {
					serviceURL := testutils.URLMust("generic://host.tld/webhook?template=JSON&messagekey=body&titlekey=header")
					err := service.Initialize(serviceURL, logger)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					params := types.Params{"header": "test title", "body": "test message"}
					payload, err := service.GetPayload(service.Config, params)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					contents, err := io.ReadAll(payload)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(string(contents)).To(gomega.MatchJSON(`{
						"header": "test title",
						"body": "test message"
					}`))
				})
			})
		})
		ginkgo.When("a valid template is specified", func() {
			ginkgo.It("applies the template to the message payload", func() {
				serviceURL := testutils.URLMust("generic://host.tld/webhook?template=news")
				err := service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				err = service.SetTemplateString("news", `{{.title}} ==> {{.message}}`)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				params := types.Params{"title": "BREAKING NEWS", "message": "it's today!"}
				payload, err := service.GetPayload(service.Config, params)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				contents, err := io.ReadAll(payload)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(string(contents)).To(gomega.Equal("BREAKING NEWS ==> it's today!"))
			})
			ginkgo.When("given nil params", func() {
				ginkgo.It("applies template with message data", func() {
					serviceURL := testutils.URLMust("generic://host.tld/webhook?template=arrows")
					err := service.Initialize(serviceURL, logger)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					err = service.SetTemplateString("arrows", `==> {{.message}} <==`)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					payload, err := service.GetPayload(service.Config, types.Params{"message": "LOOK AT ME"})
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					contents, err := io.ReadAll(payload)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(string(contents)).To(gomega.Equal("==> LOOK AT ME <=="))
				})
			})
		})
		ginkgo.When("an unknown template is specified", func() {
			ginkgo.It("returns an error", func() {
				serviceURL := testutils.URLMust("generic://host.tld/webhook?template=missing")
				err := service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				_, err = service.GetPayload(service.Config, nil)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("sending the payload", func() {
		ginkgo.BeforeEach(func() {
			httpmock.Activate()
			service = &generic.Service{}
			service.SetLogger(logger)
		})
		ginkgo.AfterEach(func() {
			httpmock.DeactivateAndReset()
		})

		ginkgo.When("sending via webhook URL", func() {
			ginkgo.It("succeeds if the server accepts the payload", func() {
				serviceURL := testutils.URLMust("generic://host.tld/webhook")
				err := service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				httpmock.RegisterResponder("POST", TestWebhookURL, httpmock.NewStringResponder(200, ""))
				err = service.Send("Message", nil)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
			ginkgo.It("reports an error if sending fails", func() {
				serviceURL := testutils.URLMust("generic://host.tld/webhook")
				err := service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				httpmock.RegisterResponder("POST", TestWebhookURL, httpmock.NewErrorResponder(errors.New("dummy error")))
				err = service.Send("Message", nil)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
			ginkgo.It("includes custom headers in the request", func() {
				serviceURL := testutils.URLMust("generic://host.tld/webhook?@authorization=frend")
				err := service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				httpmock.RegisterResponder("POST", TestWebhookURL,
					func(req *http.Request) (*http.Response, error) {
						gomega.Expect(req.Header.Get("Authorization")).To(gomega.Equal("frend"))

						return httpmock.NewStringResponse(200, ""), nil
					})
				err = service.Send("Message", nil)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
			ginkgo.It("includes extra data in JSON payload", func() {
				serviceURL := testutils.URLMust("generic://host.tld/webhook?template=json&$context=inside+joke")
				err := service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				httpmock.RegisterResponder("POST", TestWebhookURL,
					func(req *http.Request) (*http.Response, error) {
						body, err := io.ReadAll(req.Body)
						gomega.Expect(err).NotTo(gomega.HaveOccurred())
						gomega.Expect(string(body)).To(gomega.MatchJSON(`{"message":"Message","context":"inside joke"}`))

						return httpmock.NewStringResponse(200, ""), nil
					})
				err = service.Send("Message", nil)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
			ginkgo.It("uses the configured HTTP method", func() {
				serviceURL := testutils.URLMust("generic://host.tld/webhook?method=GET")
				err := service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				httpmock.RegisterResponder("GET", TestWebhookURL, httpmock.NewStringResponder(200, ""))
				err = service.Send("Message", nil)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
			ginkgo.It("does not mutate the given params", func() {
				serviceURL := testutils.URLMust("generic://host.tld/webhook?method=GET")
				err := service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				httpmock.RegisterResponder("GET", TestWebhookURL, httpmock.NewStringResponder(200, ""))
				params := types.Params{"title": "TITLE"}
				err = service.Send("Message", &params)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(params).To(gomega.Equal(types.Params{"title": "TITLE"}))
			})
		})
	})
})
