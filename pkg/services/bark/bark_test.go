package bark_test

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/nicholas-fedor/shoutrrr/internal/testutils"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/bark"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
)

// TestBark runs the Ginkgo test suite for the bark package.
func TestBark(t *testing.T) {
	format.CharactersAroundMismatchToInclude = 20 // Show more context in failure output

	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Shoutrrr Bark Suite")
}

var (
	service    *bark.Service = &bark.Service{}        // Bark service instance for testing
	envBarkURL *url.URL                               // Environment-provided URL for integration tests
	logger     *log.Logger   = testutils.TestLogger() // Shared logger for tests
	_                        = ginkgo.BeforeSuite(func() {
		// Load the integration test URL from environment, if available
		var err error
		envBarkURL, err = url.Parse(os.Getenv("SHOUTRRR_BARK_URL"))
		if err != nil {
			envBarkURL = &url.URL{} // Default to empty URL if parsing fails
		}
	})
)

var _ = ginkgo.Describe("the bark service", func() {
	ginkgo.When("running integration tests", func() {
		ginkgo.It("sends a message successfully with a valid ENV URL", func() {
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
		ginkgo.When("getting an API URL", func() {
			ginkgo.It("constructs the expected URL for various path formats", func() {
				gomega.Expect(getAPIForPath("path")).To(gomega.Equal("https://host/path/endpoint"))
				gomega.Expect(getAPIForPath("/path")).To(gomega.Equal("https://host/path/endpoint"))
				gomega.Expect(getAPIForPath("/path/")).To(gomega.Equal("https://host/path/endpoint"))
				gomega.Expect(getAPIForPath("path/")).To(gomega.Equal("https://host/path/endpoint"))
				gomega.Expect(getAPIForPath("/")).To(gomega.Equal("https://host/endpoint"))
				gomega.Expect(getAPIForPath("")).To(gomega.Equal("https://host/endpoint"))
			})
		})

		ginkgo.When("only required fields are set", func() {
			ginkgo.It("applies default values to optional fields", func() {
				serviceURL := testutils.URLMust("bark://:devicekey@hostname")
				gomega.Expect(service.Initialize(serviceURL, logger)).To(gomega.Succeed())
				gomega.Expect(*service.Config).To(gomega.Equal(bark.Config{
					Host:      "hostname",
					DeviceKey: "devicekey",
					Scheme:    "https",
				}))
			})
		})

		ginkgo.When("parsing the configuration URL", func() {
			ginkgo.It("preserves all fields after de-/serialization", func() {
				testURL := "bark://:device-key@example.com:2225/?badge=5&category=CAT&group=GROUP&scheme=http&title=TITLE&url=URL"
				config := &bark.Config{}
				gomega.Expect(config.SetURL(testutils.URLMust(testURL))).To(gomega.Succeed(), "verifying")
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

		ginkgo.It("sends successfully when the server accepts the payload", func() {
			serviceURL := testutils.URLMust("bark://:devicekey@hostname")
			gomega.Expect(service.Initialize(serviceURL, logger)).To(gomega.Succeed())
			httpmock.RegisterResponder("POST", service.Config.GetAPIURL("push"),
				testutils.JSONRespondMust(200, bark.APIResponse{
					Code:    http.StatusOK,
					Message: "OK",
				}))
			gomega.Expect(service.Send("Message", nil)).To(gomega.Succeed())
		})

		ginkgo.It("reports an error for a server error response", func() {
			serviceURL := testutils.URLMust("bark://:devicekey@hostname")
			gomega.Expect(service.Initialize(serviceURL, logger)).To(gomega.Succeed())
			httpmock.RegisterResponder("POST", service.Config.GetAPIURL("push"),
				testutils.JSONRespondMust(500, bark.APIResponse{
					Code:    500,
					Message: "someone turned off the internet",
				}))
			gomega.Expect(service.Send("Message", nil)).To(gomega.HaveOccurred())
		})

		ginkgo.It("handles an unexpected server response gracefully", func() {
			serviceURL := testutils.URLMust("bark://:devicekey@hostname")
			gomega.Expect(service.Initialize(serviceURL, logger)).To(gomega.Succeed())
			httpmock.RegisterResponder("POST", service.Config.GetAPIURL("push"),
				testutils.JSONRespondMust(200, bark.APIResponse{
					Code:    500,
					Message: "For some reason, the response code and HTTP code is different?",
				}))
			gomega.Expect(service.Send("Message", nil)).To(gomega.HaveOccurred())
		})

		ginkgo.It("handles communication errors without panicking", func() {
			httpmock.DeactivateAndReset() // Ensure no mocks interfere
			serviceURL := testutils.URLMust("bark://:devicekey@nonresolvablehostname")
			gomega.Expect(service.Initialize(serviceURL, logger)).To(gomega.Succeed())
			gomega.Expect(service.Send("Message", nil)).To(gomega.HaveOccurred())
		})
	})

	ginkgo.Describe("the basic service API", func() {
		ginkgo.Describe("the service config", func() {
			ginkgo.It("implements basic service config API methods correctly", func() {
				testutils.TestConfigGetInvalidQueryValue(&bark.Config{})
				testutils.TestConfigSetInvalidQueryValue(&bark.Config{}, "bark://:mock-device@host/?foo=bar")
				testutils.TestConfigSetDefaultValues(&bark.Config{})
				testutils.TestConfigGetEnumsCount(&bark.Config{}, 0)
				testutils.TestConfigGetFieldsCount(&bark.Config{}, 9)
			})
		})

		ginkgo.Describe("the service instance", func() {
			ginkgo.BeforeEach(func() {
				httpmock.Activate()
			})
			ginkgo.AfterEach(func() {
				httpmock.DeactivateAndReset()
			})
			ginkgo.It("implements basic service API methods correctly", func() {
				serviceURL := testutils.URLMust("bark://:devicekey@hostname")
				gomega.Expect(service.Initialize(serviceURL, logger)).To(gomega.Succeed())
				testutils.TestServiceSetInvalidParamValue(service, "foo", "bar")
			})
			ginkgo.It("returns the correct service identifier", func() {
				// No initialization needed since GetID is static
				gomega.Expect(service.GetID()).To(gomega.Equal("bark"))
			})
		})
	})
})

// getAPIForPath is a helper to construct an API URL for testing.
func getAPIForPath(path string) string {
	c := bark.Config{Host: "host", Path: path, Scheme: "https"}

	return c.GetAPIURL("endpoint")
}
