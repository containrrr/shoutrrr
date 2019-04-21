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
	plugin *pushover.PushoverPlugin
	envPushoverUrl string
)
var _ = Describe("the pushover plugin", func() {
	BeforeSuite(func() {
		plugin = &pushover.PushoverPlugin{}
		envPushoverUrl = os.Getenv("SHOUTRRR_PUSHOVER_URL")
	})
	When("running integration tests", func() {
		It("should work", func() {
			if envPushoverUrl == "" {
				return
			}
			err := plugin.Send(envPushoverUrl, "this is an integration test")
			Expect(err).NotTo(HaveOccurred())
		})
	})
})