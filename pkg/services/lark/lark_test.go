package lark

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/containrrr/shoutrrr/internal/testutils"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var tt *testing.T

func TestLark(t *testing.T) {
	RegisterFailHandler(Fail)
	tt = t
	RunSpecs(t, "Shoutrrr Lark Suite")
}

var (
	service *Service
	logger  *log.Logger
	_       = BeforeSuite(func() {
		logger = log.New(GinkgoWriter, "Test", log.LstdFlags)
	})
)

const fullURL = "lark://open.larksuite.com/token?secret=sss"

var _ = Describe("Lark Test", func() {
	BeforeEach(func() {
		service = &Service{}
	})
	When("parsing the configuration URL", func() {
		It("should be identical after de-/serialization", func() {
			url := testutils.URLMust(fullURL)
			config := &Config{}
			pkr := format.NewPropKeyResolver(config)
			err := config.setURL(&pkr, url)
			Expect(err).ShouldNot(HaveOccurred())
			outputURL := config.GetURL()
			GinkgoT().Logf("\n\n%s\n%s\n\n-", outputURL, fullURL)
			Expect(outputURL.String()).To(Equal(fullURL))
		})
	})
	Context("basic service API methods", func() {
		var config *Config
		BeforeEach(func() {
			config = &Config{}
		})
		It("should not allow getting invalid query values", func() {
			testutils.TestConfigGetInvalidQueryValue(config)
		})
		It("should not allow setting invalid query values", func() {
			testutils.TestConfigSetInvalidQueryValue(config, "lark://endpoint/token?secret=sss&foo=bar")
		})
		It("should have the expected number of fields and enums", func() {
			testutils.TestConfigGetEnumsCount(config, 0)
			testutils.TestConfigGetFieldsCount(config, 2)
		})
	})
	When("sending a message", func() {
		When("the message is too large", func() {
			It("should return large messes error", func() {
				data := make([]string, 410)
				for i := range data {
					data[i] = "0123456789"
				}
				message := strings.Join(data, "")
				service := Service{config: &Config{}}
				Expect(service.Send(message, nil)).To(MatchError(ErrLargeMessage))
			})
		})
		When("the service is not configured correctly", func() {
			It("should fail when to send message", func() {
				service := Service{config: &Config{}}
				Expect(service.Send("test message", nil)).To(MatchError(ErrInvalidHost))
				service.config.Host = larkHost
				Expect(service.Send("test message", nil)).To(MatchError(ErrNoPath))
			})
		})
		When("an invalid param is passed", func() {
			It("should fail to send messages", func() {
				service := Service{config: &Config{}}
				Expect(service.Send("test message", &types.Params{"invalid": "value"})).To(MatchError("invalid is not a valid config key []"))
			})
		})
		Context("sending message by http", func() {
			BeforeEach(func() {
				httpmock.ActivateNonDefault(httpClient)
			})
			AfterEach(func() {
				httpmock.DeactivateAndReset()
			})
			It("should send text message success", func() {
				httpmock.RegisterResponder(
					http.MethodPost,
					"/open-apis/bot/v2/hook/token",
					httpmock.NewJsonResponderOrPanic(http.StatusOK, map[string]any{"code": 0, "msg": "success"}),
				)
				service := &Service{}
				err := service.Initialize(testutils.URLMust(fullURL), logger)
				Expect(err).ShouldNot(HaveOccurred())
				err = service.Send("message", nil)
				Expect(err).ShouldNot(HaveOccurred())
			})
			It("should send post message success", func() {
				httpmock.RegisterResponder(
					http.MethodPost,
					"/open-apis/bot/v2/hook/token",
					httpmock.NewJsonResponderOrPanic(http.StatusOK, map[string]any{"code": 0, "msg": "success"}),
				)
				service := &Service{}
				err := service.Initialize(testutils.URLMust(fullURL), logger)
				Expect(err).ShouldNot(HaveOccurred())
				err = service.Send("message", &types.Params{"title": "title"})
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("should return error while request error", func() {
				httpmock.RegisterResponder(
					http.MethodPost,
					"/open-apis/bot/v2/hook/token",
					httpmock.NewErrorResponder(errors.New("network error")),
				)
				service := &Service{}
				err := service.Initialize(testutils.URLMust(fullURL), logger)
				Expect(err).ShouldNot(HaveOccurred())
				err = service.Send("message", nil)
				Expect(err).Should(MatchError(ContainSubstring("network error")))
			})
			It("should return error while response not json", func() {
				httpmock.RegisterResponder(
					http.MethodPost,
					"/open-apis/bot/v2/hook/token",
					httpmock.NewStringResponder(http.StatusOK, "some response"),
				)
				service := &Service{}
				err := service.Initialize(testutils.URLMust(fullURL), logger)
				Expect(err).ShouldNot(HaveOccurred())
				err = service.Send("message", nil)
				Expect(err).Should(MatchError(ContainSubstring("invalid character")))
			})
			It("should return error while response code not zero", func() {
				httpmock.RegisterResponder(
					http.MethodPost,
					"/open-apis/bot/v2/hook/token",
					httpmock.NewJsonResponderOrPanic(http.StatusOK, map[string]any{"code": 1, "msg": "some error"}),
				)
				service := &Service{}
				err := service.Initialize(testutils.URLMust(fullURL), logger)
				Expect(err).ShouldNot(HaveOccurred())
				err = service.Send("message", nil)
				Expect(err).Should(MatchError(ContainSubstring("some error")))
			})
		})
	})
})
