package mattermost

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/containrrr/shoutrrr/internal/testutils"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	service          *Service
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
			Expect(service.Initialize(serviceURL, testutils.TestLogger())).To(Succeed())
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
		When("generating a config object with username only", func() {
			mattermostURL, _ := url.Parse("mattermost://testUserName@mattermost.my-domain.com/thisshouldbeanapitoken")
			config := &Config{}
			err := config.SetURL(mattermostURL)
			It("should not have caused an error", func() {
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
			mattermostURL, _ := url.Parse("mattermost://mattermost.my-domain.com/thisshouldbeanapitoken/testChannel")
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
			mattermostURL, _ := url.Parse("mattermost://testUserName@mattermost.my-domain.com/thisshouldbeanapitoken/testChannel")
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
	When("generating a config object", func() {
		It("should not set icon", func() {
			slackURL, _ := url.Parse("mattermost://AAAAAAAAA/BBBBBBBBB")
			config, configError := CreateConfigFromURL(slackURL)

			Expect(configError).NotTo(HaveOccurred())
			Expect(config.Icon).To(BeEmpty())
		})
		It("should set icon", func() {
			slackURL, _ := url.Parse("mattermost://AAAAAAAAA/BBBBBBBBB?icon=test")
			config, configError := CreateConfigFromURL(slackURL)

			Expect(configError).NotTo(HaveOccurred())
			Expect(config.Icon).To(BeIdenticalTo("test"))
		})
	})
	Describe("creating the payload", func() {
		Describe("the icon fields", func() {
			payload := JSON{}
			It("should set IconURL when the configured icon looks like an URL", func() {
				payload.SetIcon("https://example.com/logo.png")
				Expect(payload.IconURL).To(Equal("https://example.com/logo.png"))
				Expect(payload.IconEmoji).To(BeEmpty())
			})
			It("should set IconEmoji when the configured icon does not look like an URL", func() {
				payload.SetIcon("tanabata_tree")
				Expect(payload.IconEmoji).To(Equal("tanabata_tree"))
				Expect(payload.IconURL).To(BeEmpty())
			})
			It("should clear both fields when icon is empty", func() {
				payload.SetIcon("")
				Expect(payload.IconEmoji).To(BeEmpty())
				Expect(payload.IconURL).To(BeEmpty())
			})
		})
	})
	Describe("Sending messages", func() {
		When("sending a message completely without parameters", func() {
			mattermostURL, _ := url.Parse("mattermost://mattermost.my-domain.com/thisshouldbeanapitoken")
			config := &Config{}
			Expect(config.SetURL(mattermostURL)).To(Succeed())
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
			mattermostURL, _ := url.Parse("mattermost://testUserName@mattermost.my-domain.com/thisshouldbeanapitoken/testChannel")
			config := &Config{}
			Expect(config.SetURL(mattermostURL)).To(Succeed())
			It("should generate the correct JSON body", func() {
				json, err := CreateJSONPayload(config, "this is a message", nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(json)).To(Equal("{\"text\":\"this is a message\",\"username\":\"testUserName\",\"channel\":\"testChannel\"}"))
			})
		})
		When("sending a message with pre set username and channel but overwriting them with parameters", func() {
			mattermostURL, _ := url.Parse("mattermost://testUserName@mattermost.my-domain.com/thisshouldbeanapitoken/testChannel")
			config := &Config{}
			Expect(config.SetURL(mattermostURL)).To(Succeed())
			It("should generate the correct JSON body", func() {
				params := (*types.Params)(&map[string]string{"username": "overwriteUserName", "channel": "overwriteChannel"})
				json, err := CreateJSONPayload(config, "this is a message", params)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(json)).To(Equal("{\"text\":\"this is a message\",\"username\":\"overwriteUserName\",\"channel\":\"overwriteChannel\"}"))
			})
		})
	})

	When("parsing the configuration URL", func() {
		It("should be identical after de-/serialization", func() {
			input := "mattermost://bot@mattermost.host/token/channel"

			config := &Config{}
			Expect(config.SetURL(testutils.URLMust(input))).To(Succeed())
			Expect(config.GetURL().String()).To(Equal(input))
		})
	})

	Describe("creating configurations", func() {
		When("given a url with channel field", func() {
			It("should not throw an error", func() {
				serviceURL := testutils.URLMust(`mattermost://user@mockserver/atoken/achannel`)
				Expect((&Config{}).SetURL(serviceURL)).To(Succeed())
			})
		})
		When("given a url with title prop", func() {
			It("should not throw an error", func() {
				serviceURL := testutils.URLMust(`mattermost://user@mockserver/atoken?icon=https%3A%2F%2Fexample%2Fsomething.png`)
				Expect((&Config{}).SetURL(serviceURL)).To(Succeed())
			})
		})
		When("given a url with all fields and props", func() {
			It("should not throw an error", func() {
				serviceURL := testutils.URLMust(`mattermost://user@mockserver/atoken/achannel?icon=https%3A%2F%2Fexample%2Fsomething.png`)
				Expect((&Config{}).SetURL(serviceURL)).To(Succeed())
			})
		})
		When("given a url with invalid props", func() {
			It("should return an error", func() {
				serviceURL := testutils.URLMust(`matrix://user@mockserver/atoken?foo=bar`)
				Expect((&Config{}).SetURL(serviceURL)).To(HaveOccurred())
			})
		})
		When("parsing the configuration URL", func() {
			It("should be identical after de-/serialization", func() {
				testURL := "mattermost://user@mockserver/atoken/achannel?icon=something"

				url, err := url.Parse(testURL)
				Expect(err).NotTo(HaveOccurred(), "parsing")

				config := &Config{}
				err = config.SetURL(url)
				Expect(err).NotTo(HaveOccurred(), "verifying")

				outputURL := config.GetURL()
				fmt.Println(outputURL.String(), testURL)

				Expect(outputURL.String()).To(Equal(testURL))

			})
		})
	})

	Describe("sending the payload", func() {
		var err error
		BeforeEach(func() {
			httpmock.Activate()
		})
		AfterEach(func() {
			httpmock.DeactivateAndReset()
		})
		It("should not report an error if the server accepts the payload", func() {
			config := Config{
				Host:  "mattermost.host",
				Token: "token",
			}
			serviceURL := config.GetURL()
			service := Service{}
			err = service.Initialize(serviceURL, nil)
			Expect(err).NotTo(HaveOccurred())

			httpmock.RegisterResponder("POST", "https://mattermost.host/hooks/token", httpmock.NewStringResponder(200, ``))

			err = service.Send("Message", nil)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("the basic service API", func() {
		Describe("the service config", func() {
			It("should implement basic service config API methods correctly", func() {
				testutils.TestConfigGetInvalidQueryValue(&Config{})

				testutils.TestConfigSetDefaultValues(&Config{})

				testutils.TestConfigGetEnumsCount(&Config{}, 0)
				testutils.TestConfigGetFieldsCount(&Config{}, 4)
			})
		})
		Describe("the service instance", func() {
			BeforeEach(func() {
				httpmock.Activate()
			})
			AfterEach(func() {
				httpmock.DeactivateAndReset()
			})
			It("should implement basic service API methods correctly", func() {
				serviceURL := testutils.URLMust("bark://mockhost/mocktoken")
				Expect(service.Initialize(serviceURL, testutils.TestLogger())).To(Succeed())
				testutils.TestServiceSetInvalidParamValue(service, "foo", "bar")
			})
		})
	})
})
