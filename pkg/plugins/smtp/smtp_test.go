package smtp_test

import (
	"fmt"
	plugin2 "github.com/containrrr/shoutrrr/pkg/plugin"
	. "github.com/containrrr/shoutrrr/pkg/plugins/smtp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	url2 "net/url"
	"os"
	"testing"
)

func TestSMTP(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr SMTP Suite")
}

var (
	plugin     *Plugin
	envSMTPURL string
	opts       plugin2.PluginOpts
	config     *Config
)

var _ = Describe("the SMTP plugin", func() {

	BeforeSuite(func() {
		plugin = &Plugin{}
		envSMTPURL = os.Getenv("SHOUTRRR_SMTP_URL")
		opts = plugin2.PluginOpts{
			Verbose: true,
			Logger: log.New(GinkgoWriter, "Test", log.LstdFlags),
		}
	})
	BeforeEach(func() {
		config = &Config{}
	})
	When("running integration tests", func() {
		It("should work without errors", func() {
			if envSMTPURL == "" {
				return
			}
			url, err := url2.Parse(envSMTPURL)
			Expect(err).NotTo(HaveOccurred())

			err = plugin.Send(*url, "this is an integration test", opts)
			Expect(err).NotTo(HaveOccurred())
		})
	})
	When("parsing the configuration URL", func() {
		It("should be identical after de-/serialization", func() {
			testURL := "smtp://user:password@example.com:2225/?fromAddress=sender@example.com&fromName=Sender&toAddresses=rec1@example.com,rec2@example.com&auth=None&subject=Subject&startTls=No&useHTML=No"

			url, err := url2.Parse(testURL)
			Expect(err).NotTo(HaveOccurred(),"parsing")

			err = config.SetURL(*url)
			Expect(err).NotTo(HaveOccurred(),"verifying")

			outputURL := config.GetURL()

			fmt.Println(outputURL.String())

			Expect(outputURL.String()).To(Equal(testURL))

		})
		When("fromAddress is missing", func() {
			It("should return an error", func() {
				testURL := "smtp://user:password@example.com:2225/?toAddresses=rec1@example.com,rec2@example.com"

				url, err := url2.Parse(testURL)
				Expect(err).NotTo(HaveOccurred(), "parsing")

				err = config.SetURL(*url)
				Expect(err).To(HaveOccurred(), "verifying")
			})
		})
		When("toAddresses are missing", func(){
			It("should return an error", func() {
				testURL := "smtp://user:password@example.com:2225/?fromAddress=sender@example.com"

				url, err := url2.Parse(testURL)
				Expect(err).NotTo(HaveOccurred(), "parsing")


				err = config.SetURL(*url)
				Expect(err).To(HaveOccurred(), "verifying")
			})

		})
	})
})