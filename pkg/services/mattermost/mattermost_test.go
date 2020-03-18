package mattermost

import (
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/url"
	"os"
	"testing"
)

var (
	service     *Service
	envMattermostURL *url.URL
)

func TestMattermost(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Mattermost Suite")
}

var _ = Describe("the mattermost service", func() {
	BeforeSuite(func() {
		service = &Service{}
		envMattermostURL, _ = url.Parse(os.Getenv("SHOUTRRR_MATTERMOST_URL"))
	})
	When("running integration tests", func() {
		It("should work without errors", func() {
			if envMattermostURL.String() == "" {
				return
			}
			serviceURL, _ := url.Parse(envMattermostURL.String())
			service.Initialize(serviceURL, util.TestLogger())
			err := service.Send(
				"this is an integration test",
				nil,
			)
			Expect(err).NotTo(HaveOccurred())
		})
	})
	Describe("the mattermost config", func() {
		When("generating a config object", func() {
			mattermostURL, _ := url.Parse("mattermost://mattermost.my-domain.com/thisshouldbeanapitoken")
			config := &Config{}
			err := config.SetURL(mattermostURL)
			It("should not have caused an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should set host", func() {
				Expect(config.Host).To(Equal("mattermost.my-domain.com"))
			})
			It("should set token", func() {
				Expect(config.Token).To(Equal("thisshouldbeanapitoken"))
			})
			It("should not set channel or username", func() {
				Expect(config.Channel).To(BeEmpty())
				Expect(config.UserName).To(BeEmpty())
			})
		})
		When("generating a new config with url, that has no token", func() {
			mattermostURL, _ := url.Parse("mattermost://mattermost.my-domain.com")
			config := &Config{}
			err := config.SetURL(mattermostURL)
			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
		When("generating a new config with url, that has a token but only one more argument", func() {
			mattermostURL, _ := url.Parse("mattermost://mattermost.my-domain.com/thisshouldbeanapitoken/testUser")
			config := &Config{}
			err := config.SetURL(mattermostURL)
			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
		When("generating a config object with username only", func() {
			mattermostURL, _ := url.Parse("mattermost://mattermost.my-domain.com/thisshouldbeanapitoken/testUserName/")
			config := &Config{}
			err := config.SetURL(mattermostURL)
			It("should not hav caused an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should set username", func() {
				Expect(config.UserName).To(Equal("testUserName"))
			})
			It("should not set channel", func() {
				Expect(config.Channel).To(BeEmpty())
			})
		})
		When("generating a config object with channel only", func() {
			mattermostURL, _ := url.Parse("mattermost://mattermost.my-domain.com/thisshouldbeanapitoken//testChannel")
			config := &Config{}
			err := config.SetURL(mattermostURL)
			It("should not hav caused an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should set channel", func() {
				Expect(config.Channel).To(Equal("testChannel"))
			})
			It("should not set channel", func() {
				Expect(config.UserName).To(BeEmpty())
			})
		})
		When("generating a config object with channel an userName", func() {
			mattermostURL, _ := url.Parse("mattermost://mattermost.my-domain.com/thisshouldbeanapitoken/testUserName/testChannel")
			config := &Config{}
			err := config.SetURL(mattermostURL)
			It("should not hav caused an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should set channel", func() {
				Expect(config.Channel).To(Equal("testChannel"))
			})
			It("should set username", func() {
				Expect(config.UserName).To(Equal("testUserName"))
			})
		})
	})
	Describe("Sending messages", func() {
		When("sending a message completely without parameters", func() {
			mattermostURL, _ := url.Parse("mattermost://mattermost.my-domain.com/thisshouldbeanapitoken")
			config := &Config{}
			config.SetURL(mattermostURL)
			It("should generate the correct url to call", func() {
				generatedURL := buildURL(config)
				Expect(generatedURL).To(Equal("https://mattermost.my-domain.com/hooks/thisshouldbeanapitoken"))
			})
			It("should generate the correct JSON body", func() {
				json, err := CreateJSONPayload(config, "this is a message", nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(json)).To(Equal("{\"text\":\"this is a message\"}"))
			})
		})
		When("sending a message with pre set username and channel", func() {
			mattermostURL, _ := url.Parse("mattermost://mattermost.my-domain.com/thisshouldbeanapitoken/testUserName/testChannel")
			config := &Config{}
			config.SetURL(mattermostURL)
			It("should generate the correct JSON body", func() {
				json, err := CreateJSONPayload(config, "this is a message", nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(json)).To(Equal("{\"text\":\"this is a message\",\"username\":\"testUserName\",\"channel\":\"testChannel\"}"))
			})
		})
		When("sending a message with pre set username and channel but overwriting them with parameters", func() {
			mattermostURL, _ := url.Parse("mattermost://mattermost.my-domain.com/thisshouldbeanapitoken/testUserName/testChannel")
			config := &Config{}
			config.SetURL(mattermostURL)
			It("should generate the correct JSON body", func() {
				params := (*types.Params)(&map[string]string{"username": "overwriteUserName","channel": "overwriteChannel",})
				json, err := CreateJSONPayload(config, "this is a message", params)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(json)).To(Equal("{\"text\":\"this is a message\",\"username\":\"overwriteUserName\",\"channel\":\"overwriteChannel\"}"))
			})
		})
	})
})
