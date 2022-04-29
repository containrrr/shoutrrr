package bark

import (
	"log"

	"github.com/containrrr/shoutrrr/pkg/util"
	"github.com/jarcoal/httpmock"

	"net/http"
	"net/url"
	"os"
	"testing"

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
		logger = util.TestLogger()
		envBarkURL, _ = url.Parse(os.Getenv("SHOUTRRR_BARK_URL"))
	})

	When("running integration tests", func() {
		It("should not error out", func() {
			if envBarkURL.String() == "" {
				Skip("No integration test ENV URL was set")
				return
			}

			configURL := util.URLMust(envBarkURL.String())
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
				serviceURL := util.URLMust("bark://:devicekey@hostname")
				Expect(service.Initialize(serviceURL, logger)).To(Succeed())

				Expect(*service.config).To(Equal(Config{
					Host:      "hostname",
					DeviceKey: "devicekey",
					Scheme:    "https",
				}))
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
			serviceURL := util.URLMust("bark://:devicekey@hostname")
			Expect(service.Initialize(serviceURL, logger)).To(Succeed())

			httpmock.RegisterResponder("POST", service.config.GetAPIURL("push"), util.JSONRespondMust(200, apiResponse{
				Code:    http.StatusOK,
				Message: "OK",
			}))

			Expect(service.Send("Message", nil)).To(Succeed())
		})
		It("should not panic if an error occurs when sending the payload", func() {
			serviceURL := util.URLMust("bark://:devicekey@hostname")
			Expect(service.Initialize(serviceURL, logger)).To(Succeed())

			httpmock.RegisterResponder("POST", service.config.GetAPIURL("push"), util.JSONRespondMust(500, apiResponse{
				Code:    500,
				Message: "someone turned off the internet",
			}))

			Expect(service.Send("Message", nil)).To(HaveOccurred())
		})
	})
})

func getAPIForPath(path string) string {
	c := Config{Host: "host", Path: path, Scheme: "https"}
	return c.GetAPIURL("endpoint")
}
