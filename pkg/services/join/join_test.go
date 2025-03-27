package join_test

import (
	"net/url"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/nicholas-fedor/shoutrrr/internal/testutils"
	"github.com/nicholas-fedor/shoutrrr/pkg/format"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/join"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestJoin(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Join Suite")
}

var (
	service    *join.Service
	config     *join.Config
	pkr        format.PropKeyResolver
	envJoinURL *url.URL
	_          = ginkgo.BeforeSuite(func() {
		service = &join.Service{}
		envJoinURL, _ = url.Parse(os.Getenv("SHOUTRRR_JOIN_URL"))
	})
)

var _ = ginkgo.Describe("the join service", func() {
	ginkgo.When("running integration tests", func() {
		ginkgo.It("should work", func() {
			if envJoinURL.String() == "" {
				return
			}
			serviceURL, _ := url.Parse(envJoinURL.String())
			err := service.Initialize(serviceURL, testutils.TestLogger())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			err = service.Send("this is an integration test", nil)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
		ginkgo.It("returns the correct service identifier", func() {
			gomega.Expect(service.GetID()).To(gomega.Equal("join"))
		})
	})
})

var _ = ginkgo.Describe("the join config", func() {
	ginkgo.BeforeEach(func() {
		config = &join.Config{}
		pkr = format.NewPropKeyResolver(config)
	})
	ginkgo.When("updating it using an url", func() {
		ginkgo.It("should update the API key using the password part of the url", func() {
			url := createURL("dummy", "TestToken", "testDevice")
			err := config.SetURL(url)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(config.APIKey).To(gomega.Equal("TestToken"))
		})
		ginkgo.It("should error if supplied with an empty token", func() {
			url := createURL("user", "", "testDevice")
			expectErrorMessageGivenURL(join.APIKeyMissing, url)
		})
	})
	ginkgo.When("getting the current config", func() {
		ginkgo.It("should return the config that is currently set as an url", func() {
			config.APIKey = "test-token"

			url := config.GetURL()
			password, _ := url.User.Password()
			gomega.Expect(password).To(gomega.Equal(config.APIKey))
			gomega.Expect(url.Scheme).To(gomega.Equal("join"))
		})
	})
	ginkgo.When("setting a config key", func() {
		ginkgo.It("should split it by commas if the key is devices", func() {
			err := pkr.Set("devices", "a,b,c,d")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(config.Devices).To(gomega.Equal([]string{"a", "b", "c", "d"}))
		})
		ginkgo.It("should update icon when an icon is supplied", func() {
			err := pkr.Set("icon", "https://example.com/icon.png")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(config.Icon).To(gomega.Equal("https://example.com/icon.png"))
		})
		ginkgo.It("should update the title when it is supplied", func() {
			err := pkr.Set("title", "new title")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(config.Title).To(gomega.Equal("new title"))
		})
		ginkgo.It("should return an error if the key is not recognized", func() {
			err := pkr.Set("devicey", "a,b,c,d")
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
	})
	ginkgo.When("getting a config key", func() {
		ginkgo.It("should join it with commas if the key is devices", func() {
			config.Devices = []string{"a", "b", "c"}
			value, err := pkr.Get("devices")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(value).To(gomega.Equal("a,b,c"))
		})
		ginkgo.It("should return an error if the key is not recognized", func() {
			_, err := pkr.Get("devicey")
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
	})

	ginkgo.When("listing the query fields", func() {
		ginkgo.It("should return the keys \"devices\", \"icon\", \"title\" in alphabetical order", func() {
			fields := pkr.QueryFields()
			gomega.Expect(fields).To(gomega.Equal([]string{"devices", "icon", "title"}))
		})
	})

	ginkgo.When("parsing the configuration URL", func() {
		ginkgo.It("should be identical after de-/serialization", func() {
			input := "join://Token:apikey@join?devices=dev1%2Cdev2&icon=warning&title=hey"
			config := &join.Config{}
			gomega.Expect(config.SetURL(testutils.URLMust(input))).To(gomega.Succeed())
			gomega.Expect(config.GetURL().String()).To(gomega.Equal(input))
		})
	})

	ginkgo.Describe("sending the payload", func() {
		var err error
		ginkgo.BeforeEach(func() {
			httpmock.Activate()
		})
		ginkgo.AfterEach(func() {
			httpmock.DeactivateAndReset()
		})
		ginkgo.It("should not report an error if the server accepts the payload", func() {
			config := join.Config{
				APIKey:  "apikey",
				Devices: []string{"dev1"},
			}
			serviceURL := config.GetURL()
			service := join.Service{}
			err = service.Initialize(serviceURL, nil)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			httpmock.RegisterResponder(
				"POST",
				"https://joinjoaomgcd.appspot.com/_ah/api/messaging/v1/sendPush",
				httpmock.NewStringResponder(200, ``),
			)

			err = service.Send("Message", nil)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
	})
})

func createURL(username string, token string, devices string) *url.URL {
	return &url.URL{
		User:     url.UserPassword("Token", token),
		Host:     username,
		RawQuery: "devices=" + devices,
	}
}

func expectErrorMessageGivenURL(msg join.ErrorMessage, url *url.URL) {
	err := config.SetURL(url)
	gomega.Expect(err).To(gomega.HaveOccurred())
	gomega.Expect(err.Error()).To(gomega.Equal(string(msg)))
}
