package pushover_test

import (
	"github.com/containrrr/shoutrrr/pkg/plugins/pushover"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPushover(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pushover Suite")
}

var (
	plugin         *pushover.PushoverPlugin
	envPushoverURL string
)
var _ = Describe("the pushover plugin", func() {
	BeforeSuite(func() {
		plugin = &pushover.PushoverPlugin{}
		envPushoverURL = os.Getenv("SHOUTRRR_PUSHOVER_URL")
	})
	When("running integration tests", func() {
		It("should work", func() {
			if envPushoverURL == "" {
				return
			}
			err := plugin.Send(envPushoverURL, "this is an integration test")
			Expect(err).NotTo(HaveOccurred())
		})
	})
})