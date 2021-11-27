package discourse

import (
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/jarcoal/httpmock"
	"log"
	"net/url"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDiscourse(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Discourse Suite")
}

var (
	service *Service
	logger  *log.Logger
)

var _ = Describe("the discourse service", func() {
	BeforeSuite(func() {
		service = &Service{}
		logger = log.New(GinkgoWriter, "Test", log.LstdFlags)
	})
	Describe("the service", func() {
		It("should implement Service interface", func() {
			var impl types.Service = service
			Expect(impl).ToNot(BeNil())
		})
	})
	Describe("the config", func() {
		When("parsing the configuration URL", func() {
			It("should be identical after de-/serialization", func() {
				testURL := "discourse://user:ap1k3y@discohost.com:1443/regular?category=1&title=Test+Title&topic=4"

				url, err := url.Parse(testURL)
				Expect(err).NotTo(HaveOccurred(), "parsing")

				config := &Config{}
				err = config.SetURL(url)
				Expect(err).NotTo(HaveOccurred(), "verifying")

				outputURL := config.GetURL()

				Expect(outputURL.String()).To(Equal(testURL))

			})
		})
	})
	Describe("sending the payload", func() {
		var err error
		BeforeEach(func() {
			httpmock.Activate()
		})
		AfterEach(func() {
			httpmock.DeactivateAndReset()
		})
		It("should not report an error if the server accepts the payload", func() {
			config := Config{
				APIKey:   "tok3ntok3ntok3n",
				Username: "bot",
				Host:     "mockserver",
				Type:     PostTypes.Post,
				Title:    "A sufficiently long title",
			}
			serviceURL := config.GetURL()
			service := Service{}
			err = service.Initialize(serviceURL, logger)
			Expect(err).NotTo(HaveOccurred())

			httpmock.RegisterResponder(
				"POST",
				`https://mockserver/posts.json`,
				httpmock.NewStringResponder(200, `{}`))

			err = service.Send("Message", nil)
			Expect(err).NotTo(HaveOccurred())
		})
		It("should report the first error returned by the server", func() {
			config := Config{
				APIKey:   "tok3ntok3ntok3n",
				Username: "bot",
				Host:     "mockserver",
				Type:     PostTypes.Post,
				Title:    "",
			}
			serviceURL := config.GetURL()
			service := Service{}
			err = service.Initialize(serviceURL, logger)
			Expect(err).NotTo(HaveOccurred())

			httpmock.RegisterResponder(
				"POST",
				`https://mockserver/posts.json`,
				httpmock.NewStringResponder(422, `{"errors": [ "The first error", "The second error" ] }`))

			err = service.Send("Message", nil)
			Expect(err).To(MatchError(`discourse API: "The first error"`))
		})
	})
})
