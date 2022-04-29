package pushover_test

import (
	"errors"
	"log"
	"net/url"
	"os"
	"testing"

	"github.com/containrrr/shoutrrr/internal/testutils"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/pushover"
	"github.com/jarcoal/httpmock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const hookURL = "https://api.pushover.net/1/messages.json"

func TestPushover(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pushover Suite")
}

var (
	service        *pushover.Service
	config         *pushover.Config
	keyResolver    format.PropKeyResolver
	envPushoverURL *url.URL
	logger         *log.Logger
)
var _ = Describe("the pushover service", func() {
	BeforeSuite(func() {
		service = &pushover.Service{}
		logger = log.New(GinkgoWriter, "Test", log.LstdFlags)
		envPushoverURL, _ = url.Parse(os.Getenv("SHOUTRRR_PUSHOVER_URL"))
	})
	When("running integration tests", func() {
		It("should work", func() {
			if envPushoverURL.String() == "" {
				return
			}
			serviceURL, _ := url.Parse(envPushoverURL.String())
			var err = service.Initialize(serviceURL, testutils.TestLogger())
			Expect(err).NotTo(HaveOccurred())
			err = service.Send("this is an integration test", nil)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

var _ = Describe("the pushover config", func() {
	BeforeEach(func() {
		config = &pushover.Config{}
		keyResolver = format.NewPropKeyResolver(config)
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
			err := keyResolver.Set("devices", "a,b,c,d")
			Expect(err).NotTo(HaveOccurred())
			Expect(config.Devices).To(Equal([]string{"a", "b", "c", "d"}))
		})
		It("should update priority when a valid number is supplied", func() {
			err := keyResolver.Set("priority", "1")
			Expect(err).NotTo(HaveOccurred())
			Expect(config.Priority).To(Equal(int8(1)))
		})
		It("should update priority when a negative number is supplied", func() {
			Expect(keyResolver.Set("priority", "-1")).To(Succeed())
			Expect(config.Priority).To(BeEquivalentTo(-1))

			Expect(keyResolver.Set("priority", "-2")).To(Succeed())
			Expect(config.Priority).To(BeEquivalentTo(-2))
		})
		It("should update the title when it is supplied", func() {
			err := keyResolver.Set("title", "new title")
			Expect(err).NotTo(HaveOccurred())
			Expect(config.Title).To(Equal("new title"))
		})
		It("should return an error if priority is not a number", func() {
			err := keyResolver.Set("priority", "super-duper")
			Expect(err).To(HaveOccurred())
		})
		It("should return an error if the key is not recognized", func() {
			err := keyResolver.Set("devicey", "a,b,c,d")
			Expect(err).To(HaveOccurred())
		})
	})
	When("getting a config key", func() {
		It("should join it with commas if the key is devices", func() {
			config.Devices = []string{"a", "b", "c"}
			value, err := keyResolver.Get("devices")
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(Equal("a,b,c"))
		})
		It("should return an error if the key is not recognized", func() {
			_, err := keyResolver.Get("devicey")
			Expect(err).To(HaveOccurred())
		})
	})

	When("listing the query fields", func() {
		It("should return the keys \"devices\",\"priority\",\"title\"", func() {
			fields := keyResolver.QueryFields()
			Expect(fields).To(Equal([]string{"devices", "priority", "title"}))
		})
	})

	Describe("sending the payload", func() {
		BeforeEach(func() {
			httpmock.Activate()
		})
		AfterEach(func() {
			httpmock.DeactivateAndReset()
		})
		It("should not report an error if the server accepts the payload", func() {
			serviceURL, err := url.Parse("pushover://:apptoken@usertoken")
			Expect(err).NotTo(HaveOccurred())

			err = service.Initialize(serviceURL, logger)
			Expect(err).NotTo(HaveOccurred())

			httpmock.RegisterResponder("POST", hookURL, httpmock.NewStringResponder(200, ""))

			err = service.Send("Message", nil)
			Expect(err).NotTo(HaveOccurred())
		})
		It("should not panic if an error occurs when sending the payload", func() {
			serviceURL, err := url.Parse("pushover://:apptoken@usertoken")
			Expect(err).NotTo(HaveOccurred())

			err = service.Initialize(serviceURL, logger)
			Expect(err).NotTo(HaveOccurred())

			httpmock.RegisterResponder("POST", hookURL, httpmock.NewErrorResponder(errors.New("dummy error")))

			err = service.Send("Message", nil)
			Expect(err).To(HaveOccurred())
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
	Expect(err).To(HaveOccurred())
	Expect(err.Error()).To(Equal(string(msg)))
}
