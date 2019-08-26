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
	config         *pushover.Config
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

var _ = Describe("the pushover config", func() {
	BeforeEach(func() {
		config = &pushover.Config{}
	})
	When("updating it using an url", func() {
		It("should update the username using the host part of the url", func() {
			url := createURL("simme", "dummy")
			err := config.SetURL(url)
			Expect(err).NotTo(HaveOccurred())
			Expect(config.User).To(Equal("simme"))
		})
		It("should update the token using the password part of the url", func() {
			url := createURL("dummy", "TestToken")
			err := config.SetURL(url)
			Expect(err).NotTo(HaveOccurred())
			Expect(config.Token).To(Equal("TestToken"))
		})
		It("should error if supplied with an empty username", func() {
			url := createURL("", "token")
			expectErrorMessageGivenURL(pushover.UserMissing, url)
		})
		It("should error if supplied with an empty token", func() {
			url := createURL("user", "")
			expectErrorMessageGivenURL(pushover.TokenMissing, url)
		})
	})
	When("getting the current config", func() {
		It("should return the config that is currently set as an url", func() {
			config.User = "simme"
			config.Token = "test-token"

			url := config.GetURL()
			password, _ := url.User.Password()
			Expect(url.Host).To(Equal(config.User))
			Expect(password).To(Equal(config.Token))
			Expect(url.Scheme).To(Equal("pushover"))
		})
	})
	When("setting a config key", func() {
		It("should split it by commas if the key is devices", func() {
			err := config.Set("devices", "a,b,c,d")
			Expect(err).NotTo(HaveOccurred())
			Expect(config.Devices).To(Equal([]string { "a", "b", "c", "d"}))
		})
		It("should return an error if the key is not devices", func() {
			err := config.Set("devicey", "a,b,c,d")
			Expect(err).To(HaveOccurred())
		})
	})
	When("getting a config key", func() {
		It("should join it with commas if the key is devices", func() {
			config.Devices = []string { "a", "b", "c" }
			value, err := config.Get("devices")
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(Equal("a,b,c"))
		})
		It("should return an error if the key is not devices", func() {
			_, err := config.Get("devicey")
			Expect(err).To(HaveOccurred())
		})
	})

	When("listing the query fields", func() {
		It("should return the key \"devices\"", func() {
			fields := config.QueryFields()
			Expect(fields).To(Equal([]string { "devices" }))
		})
	})
})

func createURL(username string, token string) *url.URL {
	return &url.URL {
		User: url.UserPassword("Token", token),
		Host: username,
	}
}

func expectErrorMessageGivenURL(msg pushover.ErrorMessage, url *url.URL) {
	err := config.SetURL(url)
	Expect(err).To(HaveOccurred())
	Expect(err.Error()).To(Equal(string(msg)))
}