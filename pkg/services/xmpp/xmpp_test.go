//+build xmpp

package xmpp

import (
	"log"
	"net/url"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/containrrr/shoutrrr/internal/testutils"
	"github.com/containrrr/shoutrrr/pkg/util"
)

func TestTeams(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr XMPP Suite")
}

var (
	logger *log.Logger
)

var _ = Describe("the XMPP service", func() {
	BeforeSuite(func() {
		logger = util.TestLogger()
	})
	When("initialized with all the requisite params", func() {
		It("should initialize without any errors", func() {
			serviceURL, _ := url.Parse("xmpp://user:password@example.com/?toAddress=r@example.com")
			service := Service{}
			err := service.Initialize(serviceURL, logger)

			Expect(err).NotTo(HaveOccurred())

		})
	})
	It("should implement basic service API methods correctly", func() {
		testutils.TestConfigGetInvalidQueryValue(&Config{})
		testutils.TestConfigSetInvalidQueryValue(&Config{}, "xmpp://example.com/?toAddress=r@example.com&foo=bar")

		testutils.TestConfigGetEnumsCount(&Config{}, 0)
		testutils.TestConfigGetFieldsCount(&Config{}, 3)
	})

})
