package xmpp

import (
	"log"
	"net/url"
	"os"
	"testing"

	"github.com/containrrr/shoutrrr/internal/testutils"
	"github.com/containrrr/shoutrrr/pkg/format"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gomegaformat "github.com/onsi/gomega/format"
)

func TestXMPP(t *testing.T) {
	gomegaformat.CharactersAroundMismatchToInclude = 20
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr XMPP Suite")
}

var (
	service    *Service  = &Service{}
	envXMPPURL *url.URL
	logger     *log.Logger = testutils.TestLogger()
	_                      = BeforeSuite(func() {
		envXMPPURL, _ = url.Parse(os.Getenv("SHOUTRRR_XMPP_URL"))
	})
)

var _ = Describe("the XMPP service", func() {

	When("running integration tests", func() {
		It("should not error out", func() {
			if envXMPPURL.String() == "" {
				Skip("No integration test ENV URL was set")
				return
			}

			configURL := testutils.URLMust(envXMPPURL.String())
			Expect(service.Initialize(configURL, logger)).To(Succeed())
			Expect(service.Send("This is an integration test message", nil)).To(Succeed())
		})
	})

	Describe("the config", func() {
		When("only required fields are set", func() {
			It("should set the optional fields to the defaults", func() {
				serviceURL := testutils.URLMust("xmpp://user:pass@hostname?receiver=target@example.com")
				Expect(service.Initialize(serviceURL, logger)).To(Succeed())

				Expect(service.config.Host).To(Equal("hostname"))
				Expect(service.config.Username).To(Equal("user"))
				Expect(service.config.Password).To(Equal("pass"))
				Expect(service.config.Port).To(Equal(uint16(5222)))
				Expect(service.config.Receiver).To(ConsistOf("target@example.com"))
				Expect(service.config.StartTLS).To(BeTrue())
				Expect(service.config.TLS).To(BeFalse())
			})
		})

		When("receiver is missing from the URL", func() {
			It("should return an error", func() {
				serviceURL := testutils.URLMust("xmpp://user:pass@hostname")
				Expect(service.Initialize(serviceURL, logger)).To(MatchError(ContainSubstring("receiver missing")))
			})
		})

		When("parsing the configuration URL", func() {
			It("should be identical after de-/serialization", func() {
				testURL := "xmpp://user:pass@example.com:5222?receiver=target%40example.com&starttls=No&tls=Yes"
				config := &Config{}
				pkr := format.NewPropKeyResolver(config)
				Expect(config.setURL(&pkr, testutils.URLMust(testURL))).To(Succeed())
				Expect(config.GetURL().String()).To(Equal(testURL))
			})
		})
	})

	Describe("the basic service API", func() {
		Describe("the service config", func() {
			It("should implement basic service config API methods correctly", func() {
				testutils.TestConfigGetInvalidQueryValue(&Config{})
				testutils.TestConfigSetInvalidQueryValue(&Config{}, "xmpp://user:pass@host?receiver=t@e.com&foo=bar")

				testutils.TestConfigSetDefaultValues(&Config{})

				testutils.TestConfigGetEnumsCount(&Config{}, 0)
				testutils.TestConfigGetFieldsCount(&Config{}, 5)
			})
		})
	})
})
