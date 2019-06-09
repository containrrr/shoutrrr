package smtp_test

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/services"
	. "github.com/containrrr/shoutrrr/pkg/services/smtp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	"net/url"
	"os"
	"testing"
)

func TestSMTP(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr SMTP Suite")
}

var (
	service    *Service
	envSMTPURL string
	opts       *services.ServiceOpts
	config     *Config
)

var _ = Describe("the SMTP service", func() {

	BeforeSuite(func() {
		service = &Service{}
		envSMTPURL = os.Getenv("SHOUTRRR_SMTP_URL")
		opts = services.CreateServiceOpts(
				log.New(GinkgoWriter, "Test", log.LstdFlags),
				true,
				map[string]string {},
			)
	})
	BeforeEach(func() {
		config = &Config{}
	})
	When("running integration tests", func() {
		It("should work without errors", func() {
			if envSMTPURL == "" {
				return
			}
			serviceURL, err := url.Parse(envSMTPURL)
			Expect(err).NotTo(HaveOccurred())

			err = service.Send(serviceURL, "this is an integration test", opts)
			Expect(err).NotTo(HaveOccurred())
		})
	})
	When("parsing the configuration URL", func() {
		It("should be identical after de-/serialization", func() {
			testURL := "smtp://user:password@example.com:2225/?fromAddress=sender@example.com&fromName=Sender&toAddresses=rec1@example.com,rec2@example.com&auth=None&subject=Subject&startTls=No&useHTML=No"

			url, err := url.Parse(testURL)
			Expect(err).NotTo(HaveOccurred(),"parsing")

			err = config.SetURL(url)
			Expect(err).NotTo(HaveOccurred(),"verifying")

			outputURL := config.GetURL()

			fmt.Println(outputURL.String())

			Expect(outputURL.String()).To(Equal(testURL))

		})
		When("fromAddress is missing", func() {
			It("should return an error", func() {
				testURL := "smtp://user:password@example.com:2225/?toAddresses=rec1@example.com,rec2@example.com"

				url, err := url.Parse(testURL)
				Expect(err).NotTo(HaveOccurred(), "parsing")

				err = config.SetURL(url)
				Expect(err).To(HaveOccurred(), "verifying")
			})
		})
		When("toAddresses are missing", func(){
			It("should return an error", func() {
				testURL := "smtp://user:password@example.com:2225/?fromAddress=sender@example.com"

				url, err := url.Parse(testURL)
				Expect(err).NotTo(HaveOccurred(), "parsing")


				err = config.SetURL(url)
				Expect(err).To(HaveOccurred(), "verifying")
			})

		})
	})
})