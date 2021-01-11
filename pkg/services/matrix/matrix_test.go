package matrix_test

import (
	"github.com/containrrr/shoutrrr/internal/testutils"
	. "github.com/containrrr/shoutrrr/pkg/services/matrix"
	"github.com/jarcoal/httpmock"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	"os"
	"testing"
)

func TestMatrix(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Matrix Suite")
}

var _ = Describe("the matrix service", func() {
	var service *Service
	logger := log.New(GinkgoWriter, "Test", log.LstdFlags)
	envMatrixURL := os.Getenv("SHOUTRRR_MATRIX_URL")

	BeforeSuite(func() {

	})

	BeforeEach(func() {
		service = &Service{}
	})

	When("running integration tests", func() {
		It("should not error out", func() {
			if envMatrixURL == "" {
				return
			}
			serviceURL, err := url.Parse(envMatrixURL)
			Expect(err).NotTo(HaveOccurred())
			err = service.Initialize(serviceURL, logger)
			Expect(err).NotTo(HaveOccurred())
			err = service.Send("This is an integration test message", nil)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("creating configurations", func() {
		When("given an url", func() {

		})
	})

	Describe("the matrix client", func() {

		BeforeEach(func() {
			httpmock.Activate()
		})

		When("sending a message", func() {
			It("should not report any errors", func() {
				setupMockResponders()
				serviceURL, _ := url.Parse("matrix://user:pass@mockserver")
				err := service.Initialize(serviceURL, logger)
				Expect(err).NotTo(HaveOccurred())

				err = service.Send("Test message", nil)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		AfterEach(func() {
			httpmock.DeactivateAndReset()
		})
	})

	It("should implement basic service API methods correctly", func() {
		testutils.TestConfigGetInvalidQueryValue(&Config{})
		testutils.TestConfigSetInvalidQueryValue(&Config{}, "matrix://user:pass@host/?foo=bar")

		testutils.TestConfigGetEnumsCount(&Config{}, 0)
		testutils.TestConfigGetFieldsCount(&Config{}, 2)
	})
})

func setupMockResponders() {
	const mockServer = "https://mockserver"

	httpmock.RegisterResponder(
		"GET",
		mockServer+APILogin,
		httpmock.NewStringResponder(200, `{"flows": [ { "type": "m.login.password" } ] }`))

	httpmock.RegisterResponder(
		"POST",
		mockServer+APILogin,
		httpmock.NewStringResponder(200, `{ "access_token": "TOKEN", "home_server": "mockserver", "user_id": "test:mockerserver" }`))

	httpmock.RegisterResponder(
		"GET",
		mockServer+APISync,
		httpmock.NewStringResponder(200, `{}`))

	httpmock.RegisterResponder(
		"POST",
		mockServer+APISendMessage,
		httpmock.NewStringResponder(200, `{ "event_id": "7" }`))

}

func expectErrorOnInit(service *Service, rawURL string, logger *log.Logger) {
	serviceURL, _ := url.Parse(rawURL)
	err := service.Initialize(serviceURL, logger)
	Expect(err).To(HaveOccurred())
}
