package ntfy

import (
	"github.com/containrrr/shoutrrr/internal/testutils"
	"github.com/containrrr/shoutrrr/pkg/format"

	"log"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gomegaformat "github.com/onsi/gomega/format"
)

func TestNtfy(t *testing.T) {
	gomegaformat.CharactersAroundMismatchToInclude = 20
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Ntfy Suite")
}

var (
	service    *Service = &Service{}
	envBarkURL *url.URL
	logger     *log.Logger = testutils.TestLogger()
	_                      = BeforeSuite(func() {
		envBarkURL, _ = url.Parse(os.Getenv("SHOUTRRR_NTFY_URL"))
	})
)

var _ = Describe("the ntfy service", func() {

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

				Expect((&Config{
					Host:   "host:8080",
					Scheme: "http",
					Topic:  "topic",
				}).GetAPIURL()).To(Equal("http://host:8080/topic"))
			})
		})
		When("only required fields are set", func() {
			It("should set the optional fields to the defaults", func() {
				serviceURL := testutils.URLMust("ntfy://hostname/topic")
				Expect(service.Initialize(serviceURL, logger)).To(Succeed())

				Expect(*service.config).To(Equal(Config{
					Host:     "hostname",
					Topic:    "topic",
					Scheme:   "https",
					Tags:     []string{""},
					Actions:  []string{""},
					Priority: 3,
					Firebase: true,
					Cache:    true,
					Template: false,
					Message:  "",
				}))
			})
		})

		When("validate that template behaves correctly", func() {
			BeforeEach(func() {
				httpmock.Activate()
			})
			AfterEach(func() {
				httpmock.DeactivateAndReset()
			})

			It("should not allow templated messages if message and title are not set", func() {
				serviceURL := testutils.URLMust("ntfy://hostname/topic?template=true&message=&title=hello")
				Expect(service.Initialize(serviceURL, logger)).ShouldNot(HaveOccurred())
				Expect(service.Send("", nil)).Should(HaveOccurred())
			})
			It("should allow templated messages if title is set", func() {
				serviceURL := testutils.URLMust("ntfy://:devicekey@hostname/topic?template=true&message=&title=hello")
				Expect(service.Initialize(serviceURL, logger)).ShouldNot(HaveOccurred())

				httpmock.RegisterResponder("POST", service.config.GetAPIURL(), testutils.JSONRespondMust(200, apiResponse{
					Code:    http.StatusOK,
					Message: "OK",
				}))
				Expect(service.Send("{}", nil)).To(Succeed())
			})
			It("should allow templated messages if message is set", func() {
				serviceURL := testutils.URLMust("ntfy://:devicekey@hostname/topic?template=true&message=hello&title=")
				Expect(service.Initialize(serviceURL, logger)).ShouldNot(HaveOccurred())

				httpmock.RegisterResponder("POST", service.config.GetAPIURL(), testutils.JSONRespondMust(200, apiResponse{
					Code:    http.StatusOK,
					Message: "OK",
				}))
				Expect(service.Send("{}", nil)).To(Succeed())
			})
			It("should not allow templated messages if body is empty", func() {
				serviceURL := testutils.URLMust("ntfy://hostname/topic?template=true&message=hello&title=hello")
				Expect(service.Initialize(serviceURL, logger)).ShouldNot(HaveOccurred())
				Expect(service.Send("", nil)).Should(HaveOccurred())
			})
			It("should validate that message is empty if template is disabled", func() {
				serviceURL := testutils.URLMust("ntfy://hostname/topic?template=false&message=test")
				Expect(service.Initialize(serviceURL, logger)).ShouldNot(HaveOccurred())
				Expect(service.Send("test", nil)).Should(HaveOccurred())
			})
		})

		When("parsing the configuration URL", func() {
			It("should be identical after de-/serialization", func() {
				testURL := "ntfy://user:pass@example.com:2225/topic?cache=No&click=CLICK&firebase=No&icon=ICON&priority=Max&scheme=http&template=Yes&title=TITLE"
				config := &Config{}
				pkr := format.NewPropKeyResolver(config)
				Expect(config.setURL(&pkr, testutils.URLMust(testURL))).To(Succeed(), "verifying")
				Expect(config.GetURL().String()).To(Equal(testURL))
				Expect(config.Template).To(Equal(true))
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
			serviceURL := testutils.URLMust("ntfy://:devicekey@hostname")
			Expect(service.Initialize(serviceURL, logger)).To(Succeed())

			httpmock.RegisterResponder("POST", service.config.GetAPIURL(), testutils.JSONRespondMust(200, apiResponse{
				Code:    http.StatusOK,
				Message: "OK",
			}))

			Expect(service.Send("Message", nil)).To(Succeed())
		})
		It("should not panic if a server error occurs", func() {
			serviceURL := testutils.URLMust("ntfy://:devicekey@hostname")
			Expect(service.Initialize(serviceURL, logger)).To(Succeed())

			httpmock.RegisterResponder("POST", service.config.GetAPIURL(), testutils.JSONRespondMust(500, apiResponse{
				Code:    500,
				Message: "someone turned off the internet",
			}))

			Expect(service.Send("Message", nil)).To(HaveOccurred())
		})
		It("should not panic if a communication error occurs", func() {
			httpmock.DeactivateAndReset()
			serviceURL := testutils.URLMust("ntfy://:devicekey@nonresolvablehostname")
			Expect(service.Initialize(serviceURL, logger)).To(Succeed())
			Expect(service.Send("Message", nil)).To(HaveOccurred())
		})
	})

	Describe("the basic service API", func() {
		Describe("the service config", func() {
			It("should implement basic service config API methods correctly", func() {
				testutils.TestConfigGetInvalidQueryValue(&Config{})
				testutils.TestConfigSetInvalidQueryValue(&Config{}, "ntfy://host/topic?foo=bar")

				testutils.TestConfigSetDefaultValues(&Config{})

				testutils.TestConfigGetEnumsCount(&Config{}, 1)
				testutils.TestConfigGetFieldsCount(&Config{}, 17)
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
				serviceURL := testutils.URLMust("ntfy://:devicekey@hostname")
				Expect(service.Initialize(serviceURL, logger)).To(Succeed())
				testutils.TestServiceSetInvalidParamValue(service, "foo", "bar")
			})
		})
	})
})
