package ntfy

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/nicholas-fedor/shoutrrr/internal/testutils"
	"github.com/nicholas-fedor/shoutrrr/pkg/format"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	gomegaformat "github.com/onsi/gomega/format"
)

func TestNtfy(t *testing.T) {
	gomegaformat.CharactersAroundMismatchToInclude = 20

	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Shoutrrr Ntfy Suite")
}

var (
	service    *Service = &Service{}
	envBarkURL *url.URL
	logger     *log.Logger = testutils.TestLogger()
	_                      = ginkgo.BeforeSuite(func() {
		envBarkURL, _ = url.Parse(os.Getenv("SHOUTRRR_NTFY_URL"))
	})
)

var _ = ginkgo.Describe("the ntfy service", func() {
	ginkgo.When("running integration tests", func() {
		ginkgo.It("should not error out", func() {
			if envBarkURL.String() == "" {
				ginkgo.Skip("No integration test ENV URL was set")

				return
			}
			configURL := testutils.URLMust(envBarkURL.String())
			gomega.Expect(service.Initialize(configURL, logger)).To(gomega.Succeed())
			gomega.Expect(service.Send("This is an integration test message", nil)).To(gomega.Succeed())
		})
	})

	ginkgo.Describe("the config", func() {
		ginkgo.When("getting a API URL", func() {
			ginkgo.It("should return the expected URL", func() {
				gomega.Expect((&Config{
					Host:   "host:8080",
					Scheme: "http",
					Topic:  "topic",
				}).GetAPIURL()).To(gomega.Equal("http://host:8080/topic"))
			})
		})
		ginkgo.When("only required fields are set", func() {
			ginkgo.It("should set the optional fields to the defaults", func() {
				serviceURL := testutils.URLMust("ntfy://hostname/topic")
				gomega.Expect(service.Initialize(serviceURL, logger)).To(gomega.Succeed())
				gomega.Expect(*service.Config).To(gomega.Equal(Config{
					Host:     "hostname",
					Topic:    "topic",
					Scheme:   "https",
					Tags:     []string{""},
					Actions:  []string{""},
					Priority: 3,
					Firebase: true,
					Cache:    true,
				}))
			})
		})
		ginkgo.When("parsing the configuration URL", func() {
			ginkgo.It("should be identical after de-/serialization", func() {
				testURL := "ntfy://user:pass@example.com:2225/topic?cache=No&click=CLICK&firebase=No&icon=ICON&priority=Max&scheme=http&title=TITLE"
				config := &Config{}
				pkr := format.NewPropKeyResolver(config)
				gomega.Expect(config.setURL(&pkr, testutils.URLMust(testURL))).To(gomega.Succeed(), "verifying")
				gomega.Expect(config.GetURL().String()).To(gomega.Equal(testURL))
			})
		})
	})

	ginkgo.When("sending the push payload", func() {
		ginkgo.BeforeEach(func() {
			httpmock.Activate()
		})
		ginkgo.AfterEach(func() {
			httpmock.DeactivateAndReset()
		})

		ginkgo.It("should not report an error if the server accepts the payload", func() {
			serviceURL := testutils.URLMust("ntfy://:devicekey@hostname/testtopic")
			gomega.Expect(service.Initialize(serviceURL, logger)).To(gomega.Succeed())
			httpmock.RegisterResponder("POST", service.Config.GetAPIURL(), testutils.JSONRespondMust(200, apiResponse{
				Code:    http.StatusOK,
				Message: "OK",
			}))
			gomega.Expect(service.Send("Message", nil)).To(gomega.Succeed())
		})

		ginkgo.It("should not panic if a server error occurs", func() {
			serviceURL := testutils.URLMust("ntfy://:devicekey@hostname/testtopic")
			gomega.Expect(service.Initialize(serviceURL, logger)).To(gomega.Succeed())
			httpmock.RegisterResponder("POST", service.Config.GetAPIURL(), testutils.JSONRespondMust(500, apiResponse{
				Code:    500,
				Message: "someone turned off the internet",
			}))
			gomega.Expect(service.Send("Message", nil)).To(gomega.HaveOccurred())
		})

		ginkgo.It("should not panic if a communication error occurs", func() {
			httpmock.DeactivateAndReset()
			serviceURL := testutils.URLMust("ntfy://:devicekey@nonresolvablehostname/testtopic")
			gomega.Expect(service.Initialize(serviceURL, logger)).To(gomega.Succeed())
			gomega.Expect(service.Send("Message", nil)).To(gomega.HaveOccurred())
		})
	})

	ginkgo.Describe("the basic service API", func() {
		ginkgo.Describe("the service config", func() {
			ginkgo.It("should implement basic service config API methods correctly", func() {
				testutils.TestConfigGetInvalidQueryValue(&Config{})
				testutils.TestConfigSetInvalidQueryValue(&Config{}, "ntfy://host/topic?foo=bar")
				testutils.TestConfigSetDefaultValues(&Config{})
				testutils.TestConfigGetEnumsCount(&Config{}, 1)
				testutils.TestConfigGetFieldsCount(&Config{}, 15)
			})
		})
		ginkgo.Describe("the service instance", func() {
			ginkgo.BeforeEach(func() {
				httpmock.Activate()
			})
			ginkgo.AfterEach(func() {
				httpmock.DeactivateAndReset()
			})
			ginkgo.It("should implement basic service API methods correctly", func() {
				serviceURL := testutils.URLMust("ntfy://:devicekey@hostname/testtopic")
				gomega.Expect(service.Initialize(serviceURL, logger)).To(gomega.Succeed())
				testutils.TestServiceSetInvalidParamValue(service, "foo", "bar")
			})
		})
	})

	ginkgo.It("should return the correct service ID", func() {
		service := &Service{}
		gomega.Expect(service.GetID()).To(gomega.Equal("ntfy"))
	})
})
