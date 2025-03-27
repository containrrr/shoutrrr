package pushbullet_test

import (
	"errors"
	"net/url"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/nicholas-fedor/shoutrrr/internal/testutils"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/pushbullet"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestPushbullet(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Shoutrrr Pushbullet Suite")
}

var (
	service          *pushbullet.Service
	envPushbulletURL *url.URL
	_                = ginkgo.BeforeSuite(func() {
		service = &pushbullet.Service{}
		envPushbulletURL, _ = url.Parse(os.Getenv("SHOUTRRR_PUSHBULLET_URL"))
	})
)

var _ = ginkgo.Describe("the pushbullet service", func() {
	ginkgo.When("running integration tests", func() {
		ginkgo.It("should not error out", func() {
			if envPushbulletURL.String() == "" {
				return
			}

			serviceURL, _ := url.Parse(envPushbulletURL.String())
			err := service.Initialize(serviceURL, testutils.TestLogger())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			err = service.Send("This is an integration test message", nil)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
		ginkgo.It("returns the correct service identifier", func() {
			gomega.Expect(service.GetID()).To(gomega.Equal("pushbullet"))
		})
	})

	ginkgo.Describe("the pushbullet config", func() {
		ginkgo.When("generating a config object", func() {
			ginkgo.It("should set token", func() {
				pushbulletURL, _ := url.Parse("pushbullet://tokentokentokentokentokentokentoke")
				config := pushbullet.Config{}
				err := config.SetURL(pushbulletURL)

				gomega.Expect(config.Token).To(gomega.Equal("tokentokentokentokentokentokentoke"))
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})

			ginkgo.It("should set the device from path", func() {
				pushbulletURL, _ := url.Parse("pushbullet://tokentokentokentokentokentokentoke/test")
				config := pushbullet.Config{}
				err := config.SetURL(pushbulletURL)

				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(config.Targets).To(gomega.HaveLen(1))
				gomega.Expect(config.Targets).To(gomega.ContainElements("test"))
			})

			ginkgo.It("should set the channel from path", func() {
				pushbulletURL, _ := url.Parse("pushbullet://tokentokentokentokentokentokentoke/foo#bar")
				config := pushbullet.Config{}
				err := config.SetURL(pushbulletURL)

				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(config.Targets).To(gomega.HaveLen(2))
				gomega.Expect(config.Targets).To(gomega.ContainElements("foo", "#bar"))
			})
		})

		ginkgo.When("parsing the configuration URL", func() {
			ginkgo.It("should be identical after de-/serialization", func() {
				testURL := "pushbullet://tokentokentokentokentokentokentoke/device?title=Great+News"

				config := &pushbullet.Config{}
				err := config.SetURL(testutils.URLMust(testURL))
				gomega.Expect(err).NotTo(gomega.HaveOccurred(), "verifying")

				outputURL := config.GetURL()
				gomega.Expect(outputURL.String()).To(gomega.Equal(testURL))
			})
		})
	})

	ginkgo.Describe("building the payload", func() {
		ginkgo.It("Email target should only populate one the correct field", func() {
			push := pushbullet.PushRequest{}
			push.SetTarget("iam@email.com")
			gomega.Expect(push.Email).To(gomega.Equal("iam@email.com"))
			gomega.Expect(push.DeviceIden).To(gomega.BeEmpty())
			gomega.Expect(push.ChannelTag).To(gomega.BeEmpty())
		})

		ginkgo.It("Device target should only populate one the correct field", func() {
			push := pushbullet.PushRequest{}
			push.SetTarget("device")
			gomega.Expect(push.Email).To(gomega.BeEmpty())
			gomega.Expect(push.DeviceIden).To(gomega.Equal("device"))
			gomega.Expect(push.ChannelTag).To(gomega.BeEmpty())
		})

		ginkgo.It("Channel target should only populate one the correct field", func() {
			push := pushbullet.PushRequest{}
			push.SetTarget("#channel")
			gomega.Expect(push.Email).To(gomega.BeEmpty())
			gomega.Expect(push.DeviceIden).To(gomega.BeEmpty())
			gomega.Expect(push.ChannelTag).To(gomega.Equal("channel"))
		})
	})

	ginkgo.Describe("sending the payload", func() {
		var err error
		targetURL := "https://api.pushbullet.com/v2/pushes"
		ginkgo.BeforeEach(func() {
			httpmock.Activate()
		})

		ginkgo.AfterEach(func() {
			httpmock.DeactivateAndReset()
		})

		ginkgo.It("should not report an error if the server accepts the payload", func() {
			err = initService("pushbullet://tokentokentokentokentokentokentoke/test")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			response := pushbullet.PushResponse{
				Type:   "note",
				Body:   "Message",
				Title:  "Shoutrrr notification", // Matches default
				Active: true,
			}
			responder, _ := httpmock.NewJsonResponder(200, &response)
			httpmock.RegisterResponder("POST", targetURL, responder)

			err = service.Send("Message", nil)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should not panic if an error occurs when sending the payload", func() {
			err = initService("pushbullet://tokentokentokentokentokentokentoke/test")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			httpmock.RegisterResponder("POST", targetURL, httpmock.NewErrorResponder(errors.New("")))

			err = service.Send("Message", nil)
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should return an error if the response type is incorrect", func() {
			err = initService("pushbullet://tokentokentokentokentokentokentoke/test")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			response := pushbullet.PushResponse{
				Type:   "link", // Incorrect type
				Body:   "Message",
				Title:  "Shoutrrr notification",
				Active: true,
			}
			responder, _ := httpmock.NewJsonResponder(200, &response)
			httpmock.RegisterResponder("POST", targetURL, responder)

			err = service.Send("Message", nil)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("unexpected response type"))
		})

		ginkgo.It("should return an error if the response body does not match", func() {
			err = initService("pushbullet://tokentokentokentokentokentokentoke/test")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			response := pushbullet.PushResponse{
				Type:   "note",
				Body:   "Wrong message",
				Title:  "Shoutrrr notification",
				Active: true,
			}
			responder, _ := httpmock.NewJsonResponder(200, &response)
			httpmock.RegisterResponder("POST", targetURL, responder)

			err = service.Send("Message", nil)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("response body mismatch"))
		})

		ginkgo.It("should return an error if the response title does not match", func() {
			err = initService("pushbullet://tokentokentokentokentokentokentoke/test")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			response := pushbullet.PushResponse{
				Type:   "note",
				Body:   "Message",
				Title:  "Wrong Title",
				Active: true,
			}
			responder, _ := httpmock.NewJsonResponder(200, &response)
			httpmock.RegisterResponder("POST", targetURL, responder)

			err = service.Send("Message", nil)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("response title mismatch"))
		})

		ginkgo.It("should return an error if the push is not active", func() {
			err = initService("pushbullet://tokentokentokentokentokentokentoke/test")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			response := pushbullet.PushResponse{
				Type:   "note",
				Body:   "Message",
				Title:  "Shoutrrr notification", // Matches default
				Active: false,
			}
			responder, _ := httpmock.NewJsonResponder(200, &response)
			httpmock.RegisterResponder("POST", targetURL, responder)

			err = service.Send("Message", nil)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("push notification is not active"))
		})
	})
})

func initService(rawURL string) error {
	serviceURL, err := url.Parse(rawURL)
	gomega.ExpectWithOffset(1, err).NotTo(gomega.HaveOccurred())

	return service.Initialize(serviceURL, testutils.TestLogger())
}
