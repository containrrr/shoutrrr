package pushbullet_test

import (
	. "github.com/containrrr/shoutrrr/pkg/services/pushbullet"
	"github.com/containrrr/shoutrrr/pkg/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/url"
	"os"
	"testing"
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
			service.Initialize(serviceURL, util.TestLogger())
			err := service.Send("This is an integration test message", nil)
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
				Expect(config.Targets[0]).To(Equal("test"))
			})
			It("should set the channel from path", func() {
				pushbulletURL, _ := url.Parse("pushbullet://tokentokentokentokentokentokentoke/#test")
				config := Config{}
				err := config.SetURL(pushbulletURL)

				Expect(err).NotTo(HaveOccurred())
				Expect(config.Targets[0]).To(Equal("#test"))
			})
		})
	})
})

func expectErrorMessageGivenURL(msg ErrorMessage, pushbulletURL *url.URL) {
	err := service.Initialize(pushbulletURL, util.TestLogger())
	Expect(err).To(HaveOccurred())
	Expect(err.Error()).To(Equal(string(msg)))
}
