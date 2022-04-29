package bark

import (
	"github.com/containrrr/shoutrrr/internal/testutils"
	"github.com/containrrr/shoutrrr/pkg/format"

	"log"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gomegaformat "github.com/onsi/gomega/format"
)

func TestBark(t *testing.T) {
	gomegaformat.CharactersAroundMismatchToInclude = 20
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Bark Suite")
}

var (
	service    *Service
	envBarkURL *url.URL
	logger     *log.Logger
)

var _ = Describe("the bark service", func() {

	BeforeSuite(func() {
		service = &Service{}
		logger = testutils.TestLogger()
		envBarkURL, _ = url.Parse(os.Getenv("SHOUTRRR_BARK_URL"))
	})

	When("running integration tests", func() {
		It("should not error out", func() {
			if envBarkURL.String() == "" {
				Skip("No integration test ENV URL was set")
				return
			}

			configURL := testutils.URLMust(envBarkURL.String())
			Expect(service.Initialize(configURL, logger)).To(Succeed())
			Expect(service.Send("This is an integration test message", nil)).To(Succeed())
		})
	})

	Describe("the config", func() {
		When("getting a API URL", func() {
			It("should return the expected URL", func() {
				Expect(getAPIForPath("path")).To(Equal("https://host/path/endpoint"))
				Expect(getAPIForPath("/path")).To(Equal("https://host/path/endpoint"))
				Expect(getAPIForPath("/path/")).To(Equal("https://host/path/endpoint"))
				Expect(getAPIForPath("path/")).To(Equal("https://host/path/endpoint"))
				Expect(getAPIForPath("/")).To(Equal("https://host/endpoint"))
				Expect(getAPIForPath("")).To(Equal("https://host/endpoint"))
			})
		})
		When("only required fields are set", func() {
			It("should set the optional fields to the defaults", func() {
				serviceURL := testutils.URLMust("bark://:devicekey@hostname")
				Expect(service.Initialize(serviceURL, logger)).To(Succeed())

				Expect(*service.config).To(Equal(Config{
					Host:      "hostname",
					DeviceKey: "devicekey",
					Scheme:    "https",
				}))
			})
		})
		When("parsing the configuration URL", func() {
			It("should be identical after de-/serialization", func() {
				testURL := "bark://:device-key@example.com:2225/?badge=5&category=CAT&group=GROUP&scheme=http&title=TITLE&url=URL"
				config := &Config{}
				pkr := format.NewPropKeyResolver(config)
				Expect(config.setURL(&pkr, testutils.URLMust(testURL))).To(Succeed(), "verifying")
				Expect(config.GetURL().String()).To(Equal(testURL))
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
			serviceURL := testutils.URLMust("bark://:devicekey@hostname")
			Expect(service.Initialize(serviceURL, logger)).To(Succeed())

			httpmock.RegisterResponder("POST", service.config.GetAPIURL("push"), testutils.JSONRespondMust(200, apiResponse{
				Code:    http.StatusOK,
				Message: "OK",
			}))

			Expect(service.Send("Message", nil)).To(Succeed())
		})
		It("should not panic if a server error occurs", func() {
			serviceURL := testutils.URLMust("bark://:devicekey@hostname")
			Expect(service.Initialize(serviceURL, logger)).To(Succeed())

			httpmock.RegisterResponder("POST", service.config.GetAPIURL("push"), testutils.JSONRespondMust(500, apiResponse{
				Code:    500,
				Message: "someone turned off the internet",
			}))

			Expect(service.Send("Message", nil)).To(HaveOccurred())
		})
		It("should not panic if a server responds with an unkown message", func() {
			serviceURL := testutils.URLMust("bark://:devicekey@hostname")
			Expect(service.Initialize(serviceURL, logger)).To(Succeed())

			httpmock.RegisterResponder("POST", service.config.GetAPIURL("push"), testutils.JSONRespondMust(200, apiResponse{
				Code:    500,
				Message: "For some reason, the response code and HTTP code is different?",
			}))

			Expect(service.Send("Message", nil)).To(HaveOccurred())
		})
		It("should not panic if a communication error occurs", func() {
			httpmock.DeactivateAndReset()
			serviceURL := testutils.URLMust("bark://:devicekey@nonresolvablehostname")
			Expect(service.Initialize(serviceURL, logger)).To(Succeed())
			Expect(service.Send("Message", nil)).To(HaveOccurred())
		})
	})

	Describe("the basic service API", func() {
		Describe("the service config", func() {
			It("should implement basic service config API methods correctly", func() {
				testutils.TestConfigGetInvalidQueryValue(&Config{})
				testutils.TestConfigSetInvalidQueryValue(&Config{}, "bark://:mock-device@host/?foo=bar")

				testutils.TestConfigSetInvalidParamValue(&Config{}, "foo", "bar")
				testutils.TestConfigSetDefaultValues(&Config{})

				testutils.TestConfigGetEnumsCount(&Config{}, 0)
				testutils.TestConfigGetFieldsCount(&Config{}, 9)
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
				serviceURL := testutils.URLMust("bark://:devicekey@hostname")
				Expect(service.Initialize(serviceURL, logger)).To(Succeed())
				testutils.TestServiceSetInvalidParamValue(service, "foo", "bar")
			})
		})
	})
})

func getAPIForPath(path string) string {
	c := Config{Host: "host", Path: path, Scheme: "https"}
	return c.GetAPIURL("endpoint")
}
