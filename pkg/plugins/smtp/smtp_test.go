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
	plugin        *Plugin
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
})