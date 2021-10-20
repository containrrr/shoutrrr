package pushbullet_test

import (
	"errors"

	. "github.com/containrrr/shoutrrr/pkg/services/pushbullet"
	"github.com/containrrr/shoutrrr/pkg/util/test"
	"github.com/jarcoal/httpmock"

	"net/url"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPushbullet(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Pushbullet Suite")
}

var (
	service          *Service
	envPushbulletURL *url.URL
)

var _ = Describe("the pushbullet service", func() {

	BeforeSuite(func() {
		service = &Service{}
		envPushbulletURL, _ = url.Parse(os.Getenv("SHOUTRRR_PUSHBULLET_URL"))

	})

	When("running integration tests", func() {
		It("should not error out", func() {
			if envPushbulletURL.String() == "" {
				return
			}

			serviceURL, _ := url.Parse(envPushbulletURL.String())
			err := service.Initialize(serviceURL, test.TestLogger())
			Expect(err).NotTo(HaveOccurred())
			err = service.Send("This is an integration test message", nil)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("the pushbullet config", func() {
		When("generating a config object", func() {
			It("should set token", func() {
				pushbulletURL, _ := url.Parse("pushbullet://tokentokentokentokentokentokentoke")
				config := Config{}
				err := config.SetURL(pushbulletURL)

				Expect(config.Token).To(Equal("tokentokentokentokentokentokentoke"))
				Expect(err).NotTo(HaveOccurred())
			})
			It("should set the device from path", func() {
				pushbulletURL, _ := url.Parse("pushbullet://tokentokentokentokentokentokentoke/test")
				config := Config{}
				err := config.SetURL(pushbulletURL)

				Expect(err).NotTo(HaveOccurred())
				Expect(config.Targets).To(HaveLen(1))
				Expect(config.Targets).To(ContainElements("test"))
			})
			It("should set the channel from path", func() {
				pushbulletURL, _ := url.Parse("pushbullet://tokentokentokentokentokentokentoke/foo#bar")
				config := Config{}
				err := config.SetURL(pushbulletURL)

				Expect(err).NotTo(HaveOccurred())
				Expect(config.Targets).To(HaveLen(2))
				Expect(config.Targets).To(ContainElements("foo", "#bar"))
			})
		})

		When("parsing the configuration URL", func() {
			It("should be identical after de-/serialization", func() {
				testURL := "pushbullet://tokentokentokentokentokentokentoke/device?title=Great+News"

				config := &Config{}
				err := config.SetURL(test.URLMust(testURL))
				Expect(err).NotTo(HaveOccurred(), "verifying")

				outputURL := config.GetURL()
				Expect(outputURL.String()).To(Equal(testURL))

			})
		})
	})

	Describe("building the payload", func() {
		It("Email target should only populate one the correct field", func() {
			push := PushRequest{}
			push.SetTarget("iam@email.com")
			Expect(push.Email).To(Equal("iam@email.com"))
			Expect(push.DeviceIden).To(BeEmpty())
			Expect(push.ChannelTag).To(BeEmpty())
		})
		It("Device target should only populate one the correct field", func() {
			push := PushRequest{}
			push.SetTarget("device")
			Expect(push.Email).To(BeEmpty())
			Expect(push.DeviceIden).To(Equal("device"))
			Expect(push.ChannelTag).To(BeEmpty())
		})
		It("Channel target should only populate one the correct field", func() {
			push := PushRequest{}
			push.SetTarget("#channel")
			Expect(push.Email).To(BeEmpty())
			Expect(push.DeviceIden).To(BeEmpty())
			Expect(push.ChannelTag).To(Equal("channel"))
		})
	})

	Describe("sending the payload", func() {
		var err error
		targetURL := "https://api.pushbullet.com/v2/pushes"
		BeforeEach(func() {
			httpmock.Activate()
		})
		AfterEach(func() {
			httpmock.DeactivateAndReset()
		})
		It("should not report an error if the server accepts the payload", func() {
			err = initService("pushbullet://tokentokentokentokentokentokentoke/test")
			Expect(err).NotTo(HaveOccurred())

			responder, _ := httpmock.NewJsonResponder(200, &PushResponse{})
			httpmock.RegisterResponder("POST", targetURL, responder)

			err = service.Send("Message", nil)
			Expect(err).NotTo(HaveOccurred())
		})
		It("should not panic if an error occurs when sending the payload", func() {
			err = initService("pushbullet://tokentokentokentokentokentokentoke/test")
			Expect(err).NotTo(HaveOccurred())

			httpmock.RegisterResponder("POST", targetURL, httpmock.NewErrorResponder(errors.New("")))

			err = service.Send("Message", nil)
			Expect(err).To(HaveOccurred())
		})
	})
})

func initService(rawURL string) error {
	serviceURL, err := url.Parse(rawURL)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return service.Initialize(serviceURL, test.TestLogger())
}
