package pushover_test

import (
	"github.com/containrrr/shoutrrr/pkg/services/pushover"
	"github.com/containrrr/shoutrrr/pkg/util"

	"net/url"
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
	service        *pushover.Service
	envPushoverURL *url.URL
)
var _ = Describe("the pushover service", func() {
	BeforeSuite(func() {
		service = &pushover.Service{}
		envPushoverURL, _ = url.Parse(os.Getenv("SHOUTRRR_PUSHOVER_URL"))
	})
	When("running integration tests", func() {
		It("should work", func() {
			if envPushoverURL.String() == "" {
				return
			}
			serviceURL, _ := url.Parse(envPushoverURL.String())
			service.Initialize(serviceURL, util.TestLogger())
			err := service.Send("this is an integration test", nil)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})