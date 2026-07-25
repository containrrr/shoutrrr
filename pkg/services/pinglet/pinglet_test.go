package pinglet

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/containrrr/shoutrrr/internal/testutils"
	"github.com/containrrr/shoutrrr/pkg/format"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gomegaformat "github.com/onsi/gomega/format"
)

func TestPinglet(t *testing.T) {
	gomegaformat.CharactersAroundMismatchToInclude = 20
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Pinglet Suite")
}

var (
	service       *Service = &Service{}
	envPingletURL *url.URL
	logger        *log.Logger = testutils.TestLogger()
	_                         = BeforeSuite(func() {
		envPingletURL, _ = url.Parse(os.Getenv("SHOUTRRR_PINGLET_URL"))
	})
)

var _ = Describe("the pinglet service", func() {

	When("running integration tests", func() {
		It("should not error out", func() {
			if envPingletURL.String() == "" {
				Skip("No integration test ENV URL was set")
				return
			}

			configURL := testutils.URLMust(envPingletURL.String())
			Expect(service.Initialize(configURL, logger)).To(Succeed())
			Expect(service.Send("This is an integration test message", nil)).To(Succeed())
		})
	})

	Describe("the config", func() {
		When("getting an API URL", func() {
			It("should return the expected URL", func() {
				Expect((&Config{
					Host:      "host:8080",
					Scheme:    "http",
					Namespace: "acme",
					Topic:     "deploys",
				}).GetAPIURL()).To(Equal("http://host:8080/acme/deploys"))
			})
		})
		When("only required fields are set", func() {
			It("should set the optional fields to the defaults", func() {
				serviceURL := testutils.URLMust("pinglet://token@hostname/acme/deploys")
				Expect(service.Initialize(serviceURL, logger)).To(Succeed())

				Expect(*service.config).To(Equal(Config{
					Token:     "token",
					Host:      "hostname",
					Namespace: "acme",
					Topic:     "deploys",
					Scheme:    "https",
					Priority:  Priority.Normal,
				}))
			})
		})
		When("parsing the configuration URL", func() {
			It("should be identical after de-/serialization", func() {
				testURL := "pinglet://token@example.com:2225/acme/deploys?badges=Host%3Aweb-1&data=region%3Aeu-west&priority=urgent&scheme=http&title=TITLE"
				config := &Config{}
				pkr := format.NewPropKeyResolver(config)
				Expect(config.setURL(&pkr, testutils.URLMust(testURL))).To(Succeed(), "verifying")
				Expect(config.GetURL().String()).To(Equal(testURL))
			})
		})
		When("no API key is supplied", func() {
			It("should return an error", func() {
				serviceURL := testutils.URLMust("pinglet://hostname/acme/deploys")
				Expect(service.Initialize(serviceURL, logger)).To(HaveOccurred())
			})
		})
		When("the namespace and/or topic are missing", func() {
			It("should return an error", func() {
				Expect(service.Initialize(testutils.URLMust("pinglet://token@hostname"), logger)).To(HaveOccurred())
				Expect(service.Initialize(testutils.URLMust("pinglet://token@hostname/acme"), logger)).To(HaveOccurred())
			})
		})
		When("a path prefix is supplied", func() {
			It("should treat the last two segments as namespace/topic and keep the prefix", func() {
				serviceURL := testutils.URLMust("pinglet://token@hostname/prefix/deep/acme/deploys?scheme=http")
				Expect(service.Initialize(serviceURL, logger)).To(Succeed())

				Expect(service.config.Path).To(Equal("prefix/deep"))
				Expect(service.config.Namespace).To(Equal("acme"))
				Expect(service.config.Topic).To(Equal("deploys"))
				Expect(service.config.GetAPIURL()).To(Equal("http://hostname/prefix/deep/acme/deploys"))
			})
			It("should round-trip a prefixed URL", func() {
				testURL := "pinglet://token@hostname/prefix/acme/deploys?priority=urgent"
				Expect(service.Initialize(testutils.URLMust(testURL), logger)).To(Succeed())
				Expect(service.config.GetURL().String()).To(Equal(testURL))
			})
		})
	})

	Describe("the priority", func() {
		When("parsing priority values", func() {
			It("should map strings and prefix aliases to the correct value", func() {
				for input, expected := range map[string]priority{
					"silent": Priority.Silent,
					"s":      Priority.Silent,
					"NORMAL": Priority.Normal,
					"urgent": Priority.Urgent,
					"u":      Priority.Urgent,
				} {
					serviceURL := testutils.URLMust("pinglet://token@hostname/acme/deploys?priority=" + input)
					Expect(service.Initialize(serviceURL, logger)).To(Succeed(), input)
					Expect(service.config.Priority).To(Equal(expected), input)
				}
			})
			It("should error on an unknown priority", func() {
				serviceURL := testutils.URLMust("pinglet://token@hostname/acme/deploys?priority=invalid")
				Expect(service.Initialize(serviceURL, logger)).To(HaveOccurred())
			})
		})
	})

	When("sending the push payload", func() {
		BeforeEach(func() {
			httpmock.Activate()
		})
		AfterEach(func() {
			httpmock.DeactivateAndReset()
		})

		It("should not report an error if the server accepts the payload", func() {
			serviceURL := testutils.URLMust("pinglet://token@hostname/acme/deploys")
			Expect(service.Initialize(serviceURL, logger)).To(Succeed())

			httpmock.RegisterResponder("POST", service.config.GetAPIURL(), testutils.JSONRespondMust(200, apiResponse{Code: 200, Message: "OK"}))

			Expect(service.Send("Message", nil)).To(Succeed())
		})
		It("should send the expected payload", func() {
			serviceURL := testutils.URLMust("pinglet://mykey@app.pinglet.co.uk/acme/deploys?priority=urgent&badges=CPU:95%25,Host:web-1&data=region:eu-west&title=Hello")
			Expect(service.Initialize(serviceURL, logger)).To(Succeed())

			var capturedBody pushPayload
			var capturedAuth string
			httpmock.RegisterResponder("POST", "https://app.pinglet.co.uk/acme/deploys", func(req *http.Request) (*http.Response, error) {
				capturedAuth = req.Header.Get("Authorization")
				Expect(json.NewDecoder(req.Body).Decode(&capturedBody)).To(Succeed())
				return httpmock.NewStringResponse(200, "{}"), nil
			})

			Expect(service.Send("It works", nil)).To(Succeed())

			Expect(capturedAuth).To(Equal("Bearer mykey"))
			Expect(capturedBody.Message).To(Equal("It works"))
			Expect(capturedBody.Title).To(Equal("Hello"))
			Expect(capturedBody.Priority).To(Equal("urgent"))
			Expect(capturedBody.Badges).To(Equal(map[string]string{"CPU": "95%", "Host": "web-1"}))
			Expect(capturedBody.Data).To(Equal(map[string]string{"region": "eu-west"}))
		})
		It("should omit the title, badges and metadata when none are set", func() {
			serviceURL := testutils.URLMust("pinglet://mykey@hostname/acme/deploys?scheme=http")
			Expect(service.Initialize(serviceURL, logger)).To(Succeed())

			var capturedBody map[string]interface{}
			httpmock.RegisterResponder("POST", "http://hostname/acme/deploys", func(req *http.Request) (*http.Response, error) {
				Expect(json.NewDecoder(req.Body).Decode(&capturedBody)).To(Succeed())
				return httpmock.NewStringResponse(200, "{}"), nil
			})

			Expect(service.Send("It works", nil)).To(Succeed())

			Expect(capturedBody).ToNot(HaveKey("title"))
			Expect(capturedBody).ToNot(HaveKey("badges"))
			Expect(capturedBody).ToNot(HaveKey("data"))
			Expect(capturedBody["priority"]).To(Equal("normal"))
		})
		It("should not panic if a server error occurs", func() {
			serviceURL := testutils.URLMust("pinglet://token@hostname/acme/deploys")
			Expect(service.Initialize(serviceURL, logger)).To(Succeed())

			httpmock.RegisterResponder("POST", service.config.GetAPIURL(), testutils.JSONRespondMust(500, apiResponse{Code: 500, Message: "someone turned off the internet"}))

			Expect(service.Send("Message", nil)).To(HaveOccurred())
		})
		It("should not panic if a communication error occurs", func() {
			httpmock.DeactivateAndReset()
			serviceURL := testutils.URLMust("pinglet://token@nonresolvablehostname/acme/deploys")
			Expect(service.Initialize(serviceURL, logger)).To(Succeed())
			Expect(service.Send("Message", nil)).To(HaveOccurred())
		})
	})

	Describe("badge and metadata limits", func() {
		It("should cap badges at three entries", func() {
			serviceURL := testutils.URLMust("pinglet://token@hostname/acme/deploys?badges=a:1,b:2,c:3,d:4,e:5")
			Expect(service.Initialize(serviceURL, logger)).To(Succeed())

			Expect(service.truncatedBadges(service.config)).To(HaveLen(maxBadgeCount))
		})
		It("should truncate over-length badge keys/values", func() {
			serviceURL := testutils.URLMust("pinglet://token@hostname/acme/deploys?badges=" +
				strings.Repeat("k", 30) + ":" + strings.Repeat("v", 40))
			Expect(service.Initialize(serviceURL, logger)).To(Succeed())

			Expect(service.truncatedBadges(service.config)).To(Equal(map[string]string{
				strings.Repeat("k", maxBadgeKeyLen): strings.Repeat("v", maxBadgeValueLen),
			}))
		})
		It("should truncate over-length metadata keys/values", func() {
			serviceURL := testutils.URLMust("pinglet://token@hostname/acme/deploys?data=" +
				strings.Repeat("k", 70) + ":" + strings.Repeat("v", 300))
			Expect(service.Initialize(serviceURL, logger)).To(Succeed())

			Expect(service.truncatedData(service.config)).To(Equal(map[string]string{
				strings.Repeat("k", maxDataKeyLen): strings.Repeat("v", maxDataValueLen),
			}))
		})
	})

	Describe("the basic service API", func() {
		Describe("the service config", func() {
			It("should implement basic service config API methods correctly", func() {
				testutils.TestConfigGetInvalidQueryValue(&Config{})
				testutils.TestConfigSetInvalidQueryValue(&Config{}, "pinglet://token@host/acme/deploys?foo=bar")

				testutils.TestConfigGetEnumsCount(&Config{}, 1)
				testutils.TestConfigGetFieldsCount(&Config{}, 5)
			})
		})
		Describe("the service instance", func() {
			BeforeEach(func() {
				httpmock.Activate()
			})
			AfterEach(func() {
				httpmock.DeactivateAndReset()
			})
			It("should implement basic service API methods correctly", func() {
				serviceURL := testutils.URLMust("pinglet://token@hostname/acme/deploys")
				Expect(service.Initialize(serviceURL, logger)).To(Succeed())
				testutils.TestServiceSetInvalidParamValue(service, "foo", "bar")
			})
		})
	})
})
