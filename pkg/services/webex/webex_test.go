package webex_test

import (
	"fmt"
	"log"

	"github.com/containrrr/shoutrrr/internal/testutils"
	. "github.com/containrrr/shoutrrr/pkg/services/webex"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/jarcoal/httpmock"

	"net/url"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestWebex(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Webex Suite")
}

var (
	service     *Service
	envWebexURL *url.URL
	logger      *log.Logger
	_           = BeforeSuite(func() {
		service = &Service{}
		envWebexURL, _ = url.Parse(os.Getenv("SHOUTRRR_WEBEX_URL"))
		logger = log.New(GinkgoWriter, "Test", log.LstdFlags)
	})
)

var _ = Describe("the webex service", func() {

	When("running integration tests", func() {
		It("should work without errors", func() {
			if envWebexURL.String() == "" {
				return
			}

			serviceURL, _ := url.Parse(envWebexURL.String())
			err := service.Initialize(serviceURL, testutils.TestLogger())
			Expect(err).NotTo(HaveOccurred())

			err = service.Send(
				"this is an integration test",
				nil,
			)
			Expect(err).NotTo(HaveOccurred())
		})
	})
	Describe("the service", func() {
		It("should implement Service interface", func() {
			var impl types.Service = service
			Expect(impl).ToNot(BeNil())
		})
	})
	Describe("creating a config", func() {
		When("given an url and a message", func() {
			It("should return an error if no arguments where supplied", func() {
				serviceURL, _ := url.Parse("webex://")
				err := service.Initialize(serviceURL, nil)
				Expect(err).To(HaveOccurred())
			})
			It("should not return an error if exactly two arguments are given", func() {
				serviceURL, _ := url.Parse("webex://dummyToken@dummyRoom")
				err := service.Initialize(serviceURL, nil)
				Expect(err).NotTo(HaveOccurred())
			})
			It("should return an error if more than two arguments are given", func() {
				serviceURL, _ := url.Parse("webex://dummyToken@dummyRoom/illegal-argument")
				err := service.Initialize(serviceURL, nil)
				Expect(err).To(HaveOccurred())
			})
		})
		When("parsing the configuration URL", func() {
			It("should be identical after de-/serialization", func() {
				testURL := "webex://token@webex?rooms=room"

				url, err := url.Parse(testURL)
				Expect(err).NotTo(HaveOccurred(), "parsing")

				config := &Config{}
				err = config.SetURL(url)
				Expect(err).NotTo(HaveOccurred(), "verifying")

				outputURL := config.GetURL()

				Expect(outputURL.String()).To(Equal(testURL))

			})
		})
	})

	Describe("sending the payload", func() {
		var dummyConfig = Config{
			BotToken: "dummyToken",
			Rooms:    []string{"1", "2"},
		}

		var service Service
		BeforeEach(func() {
			httpmock.Activate()
			service = Service{}
			if err := service.Initialize(dummyConfig.GetURL(), logger); err != nil {
				panic(fmt.Errorf("service initialization failed: %w", err))
			}
		})

		AfterEach(func() {
			httpmock.DeactivateAndReset()
		})

		It("should not report an error if the server accepts the payload", func() {
			setupResponder(&dummyConfig, 200, "")

			Expect(service.Send("Message", nil)).To(Succeed())
		})

		It("should report an error if the server response is not OK", func() {
			setupResponder(&dummyConfig, 400, "")
			Expect(service.Initialize(dummyConfig.GetURL(), logger)).To(Succeed())
			Expect(service.Send("Message", nil)).NotTo(Succeed())
		})

		It("should report an error if the message is empty", func() {
			setupResponder(&dummyConfig, 400, "")
			Expect(service.Initialize(dummyConfig.GetURL(), logger)).To(Succeed())
			Expect(service.Send("", nil)).NotTo(Succeed())
		})
	})

	Describe("doing request", func() {
		dummyConfig := &Config{
			BotToken: "dummyToken",
			Rooms:    []string{"1"},
		}

		It("should add authorization header", func() {
			request, err := BuildRequestFromPayloadAndConfig("", dummyConfig.Rooms[0], dummyConfig)

			Expect(err).To(BeNil())
			Expect(request.Header.Get("Authorization")).To(Equal("Bearer dummyToken"))
		})

		// webex API rejects messages which do not define Content-Type
		It("should add content type header", func() {
			request, err := BuildRequestFromPayloadAndConfig("", dummyConfig.Rooms[0], dummyConfig)

			Expect(err).To(BeNil())
			Expect(request.Header.Get("Content-Type")).To(Equal("application/json"))
		})
	})
})

func setupResponder(config *Config, code int, body string) {
	httpmock.RegisterResponder("POST", MessagesEndpoint, httpmock.NewStringResponder(code, body))
}
