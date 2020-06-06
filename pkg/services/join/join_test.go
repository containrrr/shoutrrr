package join_test

import (
	"github.com/containrrr/shoutrrr/pkg/services/join"
	"github.com/containrrr/shoutrrr/pkg/util"

	"net/url"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestJoin(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Join Suite")
}

var (
	service    *join.Service
	config     *join.Config
	envJoinURL *url.URL
)
var _ = Describe("the join service", func() {
	BeforeSuite(func() {
		service = &join.Service{}
		envJoinURL, _ = url.Parse(os.Getenv("SHOUTRRR_JOIN_URL"))
	})
	When("running integration tests", func() {
		It("should work", func() {
			if envJoinURL.String() == "" {
				return
			}
			serviceURL, _ := url.Parse(envJoinURL.String())
			var err = service.Initialize(serviceURL, util.TestLogger())
			Expect(err).NotTo(HaveOccurred())
			err = service.Send("this is an integration test", nil)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

var _ = Describe("the join config", func() {
	BeforeEach(func() {
		config = &join.Config{}
	})
	When("updating it using an url", func() {
		It("should update the API key using the password part of the url", func() {
			url := createURL("dummy", "TestToken", "testDevice")
			err := config.SetURL(url)
			Expect(err).NotTo(HaveOccurred())
			Expect(config.APIKey).To(Equal("TestToken"))
		})
		It("should error if supplied with an empty token", func() {
			url := createURL("user", "", "testDevice")
			expectErrorMessageGivenURL(join.APIKeyMissing, url)
		})
	})
	When("getting the current config", func() {
		It("should return the config that is currently set as an url", func() {
			config.APIKey = "test-token"

			url := config.GetURL()
			password, _ := url.User.Password()
			Expect(password).To(Equal(config.APIKey))
			Expect(url.Scheme).To(Equal("join"))
		})
	})
	When("setting a config key", func() {
		It("should split it by commas if the key is devices", func() {
			err := config.Set("devices", "a,b,c,d")
			Expect(err).NotTo(HaveOccurred())
			Expect(config.Devices).To(Equal([]string{"a", "b", "c", "d"}))
		})
		It("should update icon when an icon is supplied", func() {
			err := config.Set("icon", "https://example.com/icon.png")
			Expect(err).NotTo(HaveOccurred())
			Expect(config.Icon).To(Equal("https://example.com/icon.png"))
		})
		It("should update the title when it is supplied", func() {
			err := config.Set("title", "new title")
			Expect(err).NotTo(HaveOccurred())
			Expect(config.Title).To(Equal("new title"))
		})
		It("should return an error if the key is not recognized", func() {
			err := config.Set("devicey", "a,b,c,d")
			Expect(err).To(HaveOccurred())
		})
	})
	When("getting a config key", func() {
		It("should join it with commas if the key is devices", func() {
			config.Devices = []string{"a", "b", "c"}
			value, err := config.Get("devices")
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(Equal("a,b,c"))
		})
		It("should return an error if the key is not recognized", func() {
			_, err := config.Get("devicey")
			Expect(err).To(HaveOccurred())
		})
	})

	When("listing the query fields", func() {
		It("should return the keys \"devices\", \"title\", \"icon\"", func() {
			fields := config.QueryFields()
			Expect(fields).To(Equal([]string{"devices", "title", "icon"}))
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
	Expect(err).To(HaveOccurred())
	Expect(err.Error()).To(Equal(string(msg)))
}
