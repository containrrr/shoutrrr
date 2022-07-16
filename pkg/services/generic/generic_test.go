package generic

import (
	"errors"
	"io/ioutil"
	"log"
	"net/url"
	"testing"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/jarcoal/httpmock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGeneric(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Generic Webhook Suite")
}

var (
	logger  = log.New(GinkgoWriter, "Test", log.LstdFlags)
	service *Service
)

var _ = Describe("the Generic service", func() {
	BeforeEach(func() {
		service = &Service{}
		service.SetLogger(logger)
	})
	When("parsing a custom URL", func() {
		It("should strip generic prefix before parsing", func() {
			customURL, err := url.Parse("generic+https://test.tld")
			Expect(err).NotTo(HaveOccurred())
			actualURL, err := service.GetConfigURLFromCustom(customURL)
			Expect(err).NotTo(HaveOccurred())
			_, expectedURL := testCustomURL("https://test.tld")
			Expect(actualURL.String()).To(Equal(expectedURL.String()))
		})

		When("a HTTP URL is provided", func() {
			It("should disable TLS", func() {
				config, _ := testCustomURL("http://example.com")
				Expect(config.DisableTLS).To(BeTrue())
			})
		})
		When("a HTTPS URL is provided", func() {
			It("should enable TLS", func() {
				config, _ := testCustomURL("https://example.com")
				Expect(config.DisableTLS).To(BeFalse())
			})
		})
		It("should escape conflicting custom query keys", func() {
			expectedURL := "generic://example.com/?__template=passed"
			config, srvURL := testCustomURL("https://example.com/?template=passed")
			Expect(config.Template).NotTo(Equal("passed")) // captured
			whURL := config.WebhookURL().String()
			Expect(whURL).To(Equal("https://example.com/?template=passed"))
			Expect(srvURL.String()).To(Equal(expectedURL))

		})
		It("should handle both escaped and service prop version of keys", func() {
			config, _ := testServiceURL("generic://example.com/?__template=passed&template=captured")
			Expect(config.Template).To(Equal("captured"))
			whURL := config.WebhookURL().String()
			Expect(whURL).To(Equal("https://example.com/?template=passed"))
		})
	})
	When("retrieving the webhook URL", func() {
		It("should build a valid webhook URL", func() {
			expectedURL := "https://example.com/path?foo=bar"
			config, _ := testServiceURL("generic://example.com/path?foo=bar")
			Expect(config.WebhookURL().String()).To(Equal(expectedURL))
		})

		When("TLS is disabled", func() {
			It("should use http schema", func() {
				config := Config{
					webhookURL: &url.URL{
						Host: "test.tld",
					},
					DisableTLS: true,
				}
				Expect(config.WebhookURL().Scheme).To(Equal("http"))
			})
		})
		When("TLS is not disabled", func() {
			It("should use https schema", func() {
				config := Config{
					webhookURL: &url.URL{
						Host: "test.tld",
					},
					DisableTLS: false,
				}
				Expect(config.WebhookURL().Scheme).To(Equal("https"))
			})
		})
	})

	Describe("creating a config", func() {
		When("creating a default config", func() {
			It("should not return an error", func() {
				config := &Config{}
				pkr := format.NewPropKeyResolver(config)
				err := pkr.SetDefaultProps(config)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		When("parsing the configuration URL", func() {
			It("should be identical after de-/serialization", func() {
				testURL := "generic://user:pass@host.tld/api/v1/webhook?__title=w&contenttype=a%2Fb&template=f&title=t"

				url, err := url.Parse(testURL)
				Expect(err).NotTo(HaveOccurred(), "parsing")

				config := &Config{}
				pkr := format.NewPropKeyResolver(config)
				Expect(pkr.SetDefaultProps(config)).To(Succeed())
				err = config.SetURL(url)
				Expect(err).NotTo(HaveOccurred(), "verifying")

				outputURL := config.GetURL()
				Expect(outputURL.String()).To(Equal(testURL))

			})
		})
	})

	Describe("building the payload", func() {
		var service Service
		var config Config
		BeforeEach(func() {
			service = Service{}
			config = Config{
				MessageKey: "message",
				TitleKey:   "title",
			}
		})
		When("no template is specified", func() {
			It("should use the message as payload", func() {
				payload, err := service.getPayload(&config, types.Params{"message": "test message"})
				Expect(err).NotTo(HaveOccurred())
				contents, err := ioutil.ReadAll(payload)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(contents)).To(Equal("test message"))
			})
		})
		When("template is specified as `JSON`", func() {
			It("should create a JSON object as the payload", func() {
				config.Template = "JSON"
				params := types.Params{"title": "test title"}
				updateParams(&config, params, "test message")
				payload, err := service.getPayload(&config, params)
				Expect(err).NotTo(HaveOccurred())
				contents, err := ioutil.ReadAll(payload)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(contents)).To(MatchJSON(`{
					"title":   "test title",
					"message": "test message"
				}`))
			})
			When("alternate keys are specified", func() {
				It("should create a JSON object using the specified keys", func() {
					config.Template = "JSON"
					config.MessageKey = "body"
					config.TitleKey = "header"
					params := types.Params{"title": "test title"}
					updateParams(&config, params, "test message")
					payload, err := service.getPayload(&config, params)
					Expect(err).NotTo(HaveOccurred())
					contents, err := ioutil.ReadAll(payload)
					Expect(err).NotTo(HaveOccurred())
					Expect(string(contents)).To(MatchJSON(`{
						"header":   "test title",
						"body": "test message"
					}`))
				})
			})
		})
		When("a valid template is specified", func() {
			It("should apply the template to the message payload", func() {
				err := service.SetTemplateString("news", `{{.title}} ==> {{.message}}`)
				Expect(err).NotTo(HaveOccurred())
				params := types.Params{}
				params.SetTitle("BREAKING NEWS")
				params.SetMessage("it's today!")
				config.Template = "news"
				payload, err := service.getPayload(&config, params)
				Expect(err).NotTo(HaveOccurred())
				contents, err := ioutil.ReadAll(payload)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(contents)).To(Equal("BREAKING NEWS ==> it's today!"))
			})
			When("given nil params", func() {
				It("should apply template with message data", func() {
					err := service.SetTemplateString("arrows", `==> {{.message}} <==`)
					Expect(err).NotTo(HaveOccurred())
					config.Template = "arrows"
					payload, err := service.getPayload(&config, types.Params{"message": "LOOK AT ME"})
					Expect(err).NotTo(HaveOccurred())
					contents, err := ioutil.ReadAll(payload)
					Expect(err).NotTo(HaveOccurred())
					Expect(string(contents)).To(Equal("==> LOOK AT ME <=="))
				})
			})
		})
		When("an unknown template is specified", func() {
			It("should return an error", func() {
				_, err := service.getPayload(&Config{Template: "missing"}, nil)
				Expect(err).To(HaveOccurred())
			})
		})

	})
	Describe("sending the payload", func() {
		var err error
		var service Service
		BeforeEach(func() {
			httpmock.Activate()
		})
		AfterEach(func() {
			httpmock.DeactivateAndReset()
		})

		It("should not report an error if the server accepts the payload", func() {
			serviceURL, _ := url.Parse("generic://host.tld/webhook")
			err = service.Initialize(serviceURL, logger)
			Expect(err).NotTo(HaveOccurred())

			targetURL := "https://host.tld/webhook"
			httpmock.RegisterResponder("POST", targetURL, httpmock.NewStringResponder(200, ""))

			err = service.Send("Message", nil)
			Expect(err).NotTo(HaveOccurred())
		})
		It("should not panic if an error occurs when sending the payload", func() {
			serviceURL, _ := url.Parse("generic://host.tld/webhook")
			err = service.Initialize(serviceURL, logger)
			Expect(err).NotTo(HaveOccurred())

			targetURL := "https://host.tld/webhook"
			httpmock.RegisterResponder("POST", targetURL, httpmock.NewErrorResponder(errors.New("dummy error")))

			err = service.Send("Message", nil)
			Expect(err).To(HaveOccurred())
		})
		It("should not return an error when an unknown param is encountered", func() {
			serviceURL, _ := url.Parse("generic://host.tld/webhook")
			err = service.Initialize(serviceURL, logger)
			Expect(err).NotTo(HaveOccurred())

			targetURL := "https://host.tld/webhook"
			httpmock.RegisterResponder("POST", targetURL, httpmock.NewStringResponder(200, ""))

			err = service.Send("Message", &types.Params{"unknown": "param"})
			Expect(err).NotTo(HaveOccurred())
		})
		It("should use the configured HTTP method", func() {
			serviceURL, _ := url.Parse("generic://host.tld/webhook?method=GET")
			err = service.Initialize(serviceURL, logger)
			Expect(err).NotTo(HaveOccurred())

			targetURL := "https://host.tld/webhook"
			httpmock.RegisterResponder("GET", targetURL, httpmock.NewStringResponder(200, ""))

			err = service.Send("Message", nil)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

func testCustomURL(testURL string) (*Config, *url.URL) {
	customURL, err := url.Parse(testURL)
	Expect(err).NotTo(HaveOccurred())
	config, pkr, err := ConfigFromWebhookURL(*customURL)
	Expect(err).NotTo(HaveOccurred())
	return config, config.getURL(&pkr)
}

func testServiceURL(testURL string) (*Config, *url.URL) {
	serviceURL, err := url.Parse(testURL)
	Expect(err).NotTo(HaveOccurred())
	config, pkr := DefaultConfig()
	err = config.setURL(&pkr, serviceURL)
	Expect(err).NotTo(HaveOccurred())
	return config, config.getURL(&pkr)
}
