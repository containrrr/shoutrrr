package pushover_test

import (
	"errors"
	"log"
	"net/url"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/nicholas-fedor/shoutrrr/internal/testutils"
	"github.com/nicholas-fedor/shoutrrr/pkg/format"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/pushover"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

const hookURL = "https://api.pushover.net/1/messages.json"

func TestPushover(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Pushover Suite")
}

var (
	service        *pushover.Service
	config         *pushover.Config
	keyResolver    format.PropKeyResolver
	envPushoverURL *url.URL
	logger         *log.Logger
	_              = ginkgo.BeforeSuite(func() {
		service = &pushover.Service{}
		logger = log.New(ginkgo.GinkgoWriter, "Test", log.LstdFlags)
		envPushoverURL, _ = url.Parse(os.Getenv("SHOUTRRR_PUSHOVER_URL"))
	})
)

var _ = ginkgo.Describe("the pushover service", func() {
	ginkgo.When("running integration tests", func() {
		ginkgo.It("should work", func() {
			if envPushoverURL.String() == "" {
				return
			}
			serviceURL, _ := url.Parse(envPushoverURL.String())
			err := service.Initialize(serviceURL, testutils.TestLogger())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			err = service.Send("this is an integration test", nil)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
		ginkgo.It("returns the correct service identifier", func() {
			gomega.Expect(service.GetID()).To(gomega.Equal("pushover"))
		})
	})
})

var _ = ginkgo.Describe("the pushover config", func() {
	ginkgo.BeforeEach(func() {
		config = &pushover.Config{}
		keyResolver = format.NewPropKeyResolver(config)
	})
	ginkgo.When("updating it using an url", func() {
		ginkgo.It("should update the username using the host part of the url", func() {
			url := createURL("simme", "dummy")
			err := config.SetURL(url)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(config.User).To(gomega.Equal("simme"))
		})
		ginkgo.It("should update the token using the password part of the url", func() {
			url := createURL("dummy", "TestToken")
			err := config.SetURL(url)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(config.Token).To(gomega.Equal("TestToken"))
		})
		ginkgo.It("should error if supplied with an empty username", func() {
			url := createURL("", "token")
			expectErrorMessageGivenURL(pushover.UserMissing, url)
		})
		ginkgo.It("should error if supplied with an empty token", func() {
			url := createURL("user", "")
			expectErrorMessageGivenURL(pushover.TokenMissing, url)
		})
	})
	ginkgo.When("getting the current config", func() {
		ginkgo.It("should return the config that is currently set as an url", func() {
			config.User = "simme"
			config.Token = "test-token"

			url := config.GetURL()
			password, _ := url.User.Password()
			gomega.Expect(url.Host).To(gomega.Equal(config.User))
			gomega.Expect(password).To(gomega.Equal(config.Token))
			gomega.Expect(url.Scheme).To(gomega.Equal("pushover"))
		})
	})
	ginkgo.When("setting a config key", func() {
		ginkgo.It("should split it by commas if the key is devices", func() {
			err := keyResolver.Set("devices", "a,b,c,d")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(config.Devices).To(gomega.Equal([]string{"a", "b", "c", "d"}))
		})
		ginkgo.It("should update priority when a valid number is supplied", func() {
			err := keyResolver.Set("priority", "1")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(config.Priority).To(gomega.Equal(int8(1)))
		})
		ginkgo.It("should update priority when a negative number is supplied", func() {
			gomega.Expect(keyResolver.Set("priority", "-1")).To(gomega.Succeed())
			gomega.Expect(config.Priority).To(gomega.BeEquivalentTo(-1))

			gomega.Expect(keyResolver.Set("priority", "-2")).To(gomega.Succeed())
			gomega.Expect(config.Priority).To(gomega.BeEquivalentTo(-2))
		})
		ginkgo.It("should update the title when it is supplied", func() {
			err := keyResolver.Set("title", "new title")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(config.Title).To(gomega.Equal("new title"))
		})
		ginkgo.It("should return an error if priority is not a number", func() {
			err := keyResolver.Set("priority", "super-duper")
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
		ginkgo.It("should return an error if the key is not recognized", func() {
			err := keyResolver.Set("devicey", "a,b,c,d")
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
	})
	ginkgo.When("getting a config key", func() {
		ginkgo.It("should join it with commas if the key is devices", func() {
			config.Devices = []string{"a", "b", "c"}
			value, err := keyResolver.Get("devices")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(value).To(gomega.Equal("a,b,c"))
		})
		ginkgo.It("should return an error if the key is not recognized", func() {
			_, err := keyResolver.Get("devicey")
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
	})

	ginkgo.When("listing the query fields", func() {
		ginkgo.It("should return the keys \"devices\",\"priority\",\"title\"", func() {
			fields := keyResolver.QueryFields()
			gomega.Expect(fields).To(gomega.Equal([]string{"devices", "priority", "title"}))
		})
	})

	ginkgo.Describe("sending the payload", func() {
		ginkgo.BeforeEach(func() {
			httpmock.Activate()
		})
		ginkgo.AfterEach(func() {
			httpmock.DeactivateAndReset()
		})
		ginkgo.It("should not report an error if the server accepts the payload", func() {
			serviceURL, err := url.Parse("pushover://:apptoken@usertoken")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			err = service.Initialize(serviceURL, logger)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			httpmock.RegisterResponder("POST", hookURL, httpmock.NewStringResponder(200, ""))

			err = service.Send("Message", nil)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
		ginkgo.It("should not panic if an error occurs when sending the payload", func() {
			serviceURL, err := url.Parse("pushover://:apptoken@usertoken")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			err = service.Initialize(serviceURL, logger)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			httpmock.RegisterResponder("POST", hookURL, httpmock.NewErrorResponder(errors.New("dummy error")))

			err = service.Send("Message", nil)
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
	})
})

func createURL(username string, token string) *url.URL {
	return &url.URL{
		User: url.UserPassword("Token", token),
		Host: username,
	}
}

func expectErrorMessageGivenURL(msg pushover.ErrorMessage, url *url.URL) {
	err := config.SetURL(url)
	gomega.Expect(err).To(gomega.HaveOccurred())
	gomega.Expect(err.Error()).To(gomega.Equal(string(msg)))
}
