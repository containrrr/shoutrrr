package smtp_test

import (
	. "github.com/containrrr/shoutrrr/pkg/plugins/smtp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
)

var _ = Describe("the SMTP plugin", func() {
	BeforeSuite(func() {
		plugin = &Plugin{}
		envSMTPURL = os.Getenv("SHOUTRRR_SMTP_URL")
	})
	When("running integration tests", func() {
		It("should work without errors", func() {
			if envSMTPURL == "" {
				return
			}
			err := plugin.Send(envSMTPURL, "this is an integration test")
			Expect(err).NotTo(HaveOccurred())
		})
	})
	When("parsing the configuration URL", func() {
		It("should be identical after de-/serialization", func() {
			testURL := "smtp://user:password@example.com:2225/?fromAddress=sender@example.com&fromName=Sender&toAddresses=rec1@example.com,rec2@example.com"

			config, err := plugin.CreateConfigFromURL(testURL)
			Expect(err).NotTo(HaveOccurred())

			outputURL := CreateAPIURLFromConfig(config)

			Expect(outputURL).To(Equal(testURL))

		})
		It("should return an error", func() {
			When("fromAddress is missing", func() {
				testURL := "smtp://user:password@example.com:2225/?toAddresses=rec1@example.com,rec2@example.com"
				_, err := plugin.CreateConfigFromURL(testURL)
				Expect(err).To(HaveOccurred())
			})
			When("toAddresses are missing", func() {
				testURL := "smtp://user:password@example.com:2225/?fromAddress=sender@example.com"
				_, err := plugin.CreateConfigFromURL(testURL)
				Expect(err).To(HaveOccurred())
			})

		})
	})
})