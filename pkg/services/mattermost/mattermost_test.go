package mattermost

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/nicholas-fedor/shoutrrr/internal/testutils"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var (
	service          *Service
	envMattermostURL *url.URL
	_                = ginkgo.BeforeSuite(func() {
		service = &Service{}
		envMattermostURL, _ = url.Parse(os.Getenv("SHOUTRRR_MATTERMOST_URL"))
	})
)

func TestMattermost(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Shoutrrr Mattermost Suite")
}

var _ = ginkgo.Describe("the mattermost service", func() {
	ginkgo.When("running integration tests", func() {
		ginkgo.It("should work without errors", func() {
			if envMattermostURL.String() == "" {
				return
			}
			serviceURL, _ := url.Parse(envMattermostURL.String())
			gomega.Expect(service.Initialize(serviceURL, testutils.TestLogger())).To(gomega.Succeed())
			err := service.Send(
				"this is an integration test",
				nil,
			)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
	})
	ginkgo.Describe("the mattermost config", func() {
		ginkgo.When("generating a config object", func() {
			mattermostURL, _ := url.Parse("mattermost://mattermost.my-domain.com/thisshouldbeanapitoken")
			config := &Config{}
			err := config.SetURL(mattermostURL)
			ginkgo.It("should not have caused an error", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
			ginkgo.It("should set host", func() {
				gomega.Expect(config.Host).To(gomega.Equal("mattermost.my-domain.com"))
			})
			ginkgo.It("should set token", func() {
				gomega.Expect(config.Token).To(gomega.Equal("thisshouldbeanapitoken"))
			})
			ginkgo.It("should not set channel or username", func() {
				gomega.Expect(config.Channel).To(gomega.BeEmpty())
				gomega.Expect(config.UserName).To(gomega.BeEmpty())
			})
		})
		ginkgo.When("generating a new config with url, that has no token", func() {
			ginkgo.It("should return an error", func() {
				mattermostURL, _ := url.Parse("mattermost://mattermost.my-domain.com")
				config := &Config{}
				err := config.SetURL(mattermostURL)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
		ginkgo.When("generating a config object with username only", func() {
			mattermostURL, _ := url.Parse("mattermost://testUserName@mattermost.my-domain.com/thisshouldbeanapitoken")
			config := &Config{}
			err := config.SetURL(mattermostURL)
			ginkgo.It("should not have caused an error", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
			ginkgo.It("should set username", func() {
				gomega.Expect(config.UserName).To(gomega.Equal("testUserName"))
			})
			ginkgo.It("should not set channel", func() {
				gomega.Expect(config.Channel).To(gomega.BeEmpty())
			})
		})
		ginkgo.When("generating a config object with channel only", func() {
			mattermostURL, _ := url.Parse("mattermost://mattermost.my-domain.com/thisshouldbeanapitoken/testChannel")
			config := &Config{}
			err := config.SetURL(mattermostURL)
			ginkgo.It("should not hav caused an error", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
			ginkgo.It("should set channel", func() {
				gomega.Expect(config.Channel).To(gomega.Equal("testChannel"))
			})
			ginkgo.It("should not set username", func() {
				gomega.Expect(config.UserName).To(gomega.BeEmpty())
			})
		})
		ginkgo.When("generating a config object with channel an userName", func() {
			mattermostURL, _ := url.Parse("mattermost://testUserName@mattermost.my-domain.com/thisshouldbeanapitoken/testChannel")
			config := &Config{}
			err := config.SetURL(mattermostURL)
			ginkgo.It("should not hav caused an error", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
			ginkgo.It("should set channel", func() {
				gomega.Expect(config.Channel).To(gomega.Equal("testChannel"))
			})
			ginkgo.It("should set username", func() {
				gomega.Expect(config.UserName).To(gomega.Equal("testUserName"))
			})
		})
		ginkgo.When("using DisableTLS and port", func() {
			mattermostURL, _ := url.Parse("mattermost://watchtower@home.lan:8065/token/channel?disabletls=yes")
			config := &Config{}
			gomega.Expect(config.SetURL(mattermostURL)).To(gomega.Succeed())
			ginkgo.It("should preserve host with port", func() {
				gomega.Expect(config.Host).To(gomega.Equal("home.lan:8065"))
			})
			ginkgo.It("should set DisableTLS", func() {
				gomega.Expect(config.DisableTLS).To(gomega.BeTrue())
			})
			ginkgo.It("should generate http URL", func() {
				gomega.Expect(buildURL(config)).To(gomega.Equal("http://home.lan:8065/hooks/token"))
			})
			ginkgo.It("should serialize back correctly", func() {
				gomega.Expect(config.GetURL().String()).To(gomega.Equal("mattermost://watchtower@home.lan:8065/token/channel?disabletls=Yes"))
			})
		})
		ginkgo.Describe("initializing with DisableTLS", func() {
			ginkgo.BeforeEach(func() {
				httpmock.Activate()
			})
			ginkgo.AfterEach(func() {
				httpmock.DeactivateAndReset()
			})
			ginkgo.It("should use plain HTTP transport when DisableTLS is true", func() {
				mattermostURL, _ := url.Parse("mattermost://user@host:8080/token?disabletls=yes")
				service := &Service{}
				err := service.Initialize(mattermostURL, testutils.TestLogger())
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				httpmock.ActivateNonDefault(service.httpClient)
				httpmock.RegisterResponder("POST", "http://host:8080/hooks/token", httpmock.NewStringResponder(200, ""))

				err = service.Send("Test message", nil)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(buildURL(service.Config)).To(gomega.Equal("http://host:8080/hooks/token"))
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
				config := Config{
					Host:  "mattermost.host",
					Token: "token",
				}
				serviceURL := config.GetURL()
				service := Service{}
				err = service.Initialize(serviceURL, nil)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				httpmock.ActivateNonDefault(service.httpClient)
				httpmock.RegisterResponder("POST", "https://mattermost.host/hooks/token", httpmock.NewStringResponder(200, ""))
				err = service.Send("Message", nil)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
			ginkgo.It("should return an error if the server rejects the payload", func() {
				config := Config{
					Host:  "mattermost.host",
					Token: "token",
				}
				serviceURL := config.GetURL()
				service := Service{}
				err = service.Initialize(serviceURL, nil)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				httpmock.ActivateNonDefault(service.httpClient)
				httpmock.RegisterResponder("POST", "https://mattermost.host/hooks/token", httpmock.NewStringResponder(403, "Forbidden"))
				err = service.Send("Message", nil)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(err.Error()).To(gomega.ContainSubstring("failed to send notification to service"))
				resp := httpmock.NewStringResponse(403, "Forbidden")
				resp.Status = "403 Forbidden"
				httpmock.RegisterResponder("POST", "https://mattermost.host/hooks/token", httpmock.ResponderFromResponse(resp))
			})
		})
	})

	ginkgo.When("generating a config object", func() {
		ginkgo.It("should not set icon", func() {
			slackURL, _ := url.Parse("mattermost://AAAAAAAAA/BBBBBBBBB")
			config, configError := CreateConfigFromURL(slackURL)

			gomega.Expect(configError).NotTo(gomega.HaveOccurred())
			gomega.Expect(config.Icon).To(gomega.BeEmpty())
		})
		ginkgo.It("should set icon", func() {
			slackURL, _ := url.Parse("mattermost://AAAAAAAAA/BBBBBBBBB?icon=test")
			config, configError := CreateConfigFromURL(slackURL)

			gomega.Expect(configError).NotTo(gomega.HaveOccurred())
			gomega.Expect(config.Icon).To(gomega.BeIdenticalTo("test"))
		})
	})
	ginkgo.Describe("creating the payload", func() {
		ginkgo.Describe("the icon fields", func() {
			payload := JSON{}
			ginkgo.It("should set IconURL when the configured icon looks like an URL", func() {
				payload.SetIcon("https://example.com/logo.png")
				gomega.Expect(payload.IconURL).To(gomega.Equal("https://example.com/logo.png"))
				gomega.Expect(payload.IconEmoji).To(gomega.BeEmpty())
			})
			ginkgo.It("should set IconEmoji when the configured icon does not look like an URL", func() {
				payload.SetIcon("tanabata_tree")
				gomega.Expect(payload.IconEmoji).To(gomega.Equal("tanabata_tree"))
				gomega.Expect(payload.IconURL).To(gomega.BeEmpty())
			})
			ginkgo.It("should clear both fields when icon is empty", func() {
				payload.SetIcon("")
				gomega.Expect(payload.IconEmoji).To(gomega.BeEmpty())
				gomega.Expect(payload.IconURL).To(gomega.BeEmpty())
			})
		})
	})
	ginkgo.Describe("Sending messages", func() {
		ginkgo.When("sending a message completely without parameters", func() {
			mattermostURL, _ := url.Parse("mattermost://mattermost.my-domain.com/thisshouldbeanapitoken")
			config := &Config{}
			gomega.Expect(config.SetURL(mattermostURL)).To(gomega.Succeed())
			ginkgo.It("should generate the correct url to call", func() {
				generatedURL := buildURL(config)
				gomega.Expect(generatedURL).To(gomega.Equal("https://mattermost.my-domain.com/hooks/thisshouldbeanapitoken"))
			})
			ginkgo.It("should generate the correct JSON body", func() {
				json, err := CreateJSONPayload(config, "this is a message", nil)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(string(json)).To(gomega.Equal("{\"text\":\"this is a message\"}"))
			})
		})
		ginkgo.When("sending a message with pre set username and channel", func() {
			mattermostURL, _ := url.Parse("mattermost://testUserName@mattermost.my-domain.com/thisshouldbeanapitoken/testChannel")
			config := &Config{}
			gomega.Expect(config.SetURL(mattermostURL)).To(gomega.Succeed())
			ginkgo.It("should generate the correct JSON body", func() {
				json, err := CreateJSONPayload(config, "this is a message", nil)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(string(json)).To(gomega.Equal("{\"text\":\"this is a message\",\"username\":\"testUserName\",\"channel\":\"testChannel\"}"))
			})
		})
		ginkgo.When("sending a message with pre set username and channel but overwriting them with parameters", func() {
			mattermostURL, _ := url.Parse("mattermost://testUserName@mattermost.my-domain.com/thisshouldbeanapitoken/testChannel")
			config := &Config{}
			gomega.Expect(config.SetURL(mattermostURL)).To(gomega.Succeed())
			ginkgo.It("should generate the correct JSON body", func() {
				params := (*types.Params)(&map[string]string{"username": "overwriteUserName", "channel": "overwriteChannel"})
				json, err := CreateJSONPayload(config, "this is a message", params)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(string(json)).To(gomega.Equal("{\"text\":\"this is a message\",\"username\":\"overwriteUserName\",\"channel\":\"overwriteChannel\"}"))
			})
		})
	})

	ginkgo.When("parsing the configuration URL", func() {
		ginkgo.It("should be identical after de-/serialization", func() {
			input := "mattermost://bot@mattermost.host/token/channel"

			config := &Config{}
			gomega.Expect(config.SetURL(testutils.URLMust(input))).To(gomega.Succeed())
			gomega.Expect(config.GetURL().String()).To(gomega.Equal(input))
		})
	})

	ginkgo.Describe("creating configurations", func() {
		ginkgo.When("given a url with channel field", func() {
			ginkgo.It("should not throw an error", func() {
				serviceURL := testutils.URLMust(`mattermost://user@mockserver/atoken/achannel`)
				gomega.Expect((&Config{}).SetURL(serviceURL)).To(gomega.Succeed())
			})
		})
		ginkgo.When("given a url with title prop", func() {
			ginkgo.It("should not throw an error", func() {
				serviceURL := testutils.URLMust(`mattermost://user@mockserver/atoken?icon=https%3A%2F%2Fexample%2Fsomething.png`)
				gomega.Expect((&Config{}).SetURL(serviceURL)).To(gomega.Succeed())
			})
		})
		ginkgo.When("given a url with all fields and props", func() {
			ginkgo.It("should not throw an error", func() {
				serviceURL := testutils.URLMust(`mattermost://user@mockserver/atoken/achannel?icon=https%3A%2F%2Fexample%2Fsomething.png`)
				gomega.Expect((&Config{}).SetURL(serviceURL)).To(gomega.Succeed())
			})
		})
		ginkgo.When("given a url with invalid props", func() {
			ginkgo.It("should return an error", func() {
				serviceURL := testutils.URLMust(`matrix://user@mockserver/atoken?foo=bar`)
				gomega.Expect((&Config{}).SetURL(serviceURL)).To(gomega.HaveOccurred())
			})
		})
		ginkgo.When("parsing the configuration URL", func() {
			ginkgo.It("should be identical after de-/serialization", func() {
				testURL := "mattermost://user@mockserver/atoken/achannel?icon=something"

				url, err := url.Parse(testURL)
				gomega.Expect(err).NotTo(gomega.HaveOccurred(), "parsing")

				config := &Config{}
				err = config.SetURL(url)
				gomega.Expect(err).NotTo(gomega.HaveOccurred(), "verifying")

				outputURL := config.GetURL()
				fmt.Println(outputURL.String(), testURL)

				gomega.Expect(outputURL.String()).To(gomega.Equal(testURL))
			})
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
			config := Config{
				Host:  "mattermost.host",
				Token: "token",
			}
			serviceURL := config.GetURL()
			service := Service{}
			err = service.Initialize(serviceURL, nil)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			httpmock.ActivateNonDefault(service.httpClient)

			httpmock.RegisterResponder("POST", "https://mattermost.host/hooks/token", httpmock.NewStringResponder(200, ``))

			err = service.Send("Message", nil)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
	})

	ginkgo.Describe("the basic service API", func() {
		ginkgo.Describe("the service config", func() {
			ginkgo.It("should implement basic service config API methods correctly", func() {
				testutils.TestConfigGetInvalidQueryValue(&Config{})

				testutils.TestConfigSetDefaultValues(&Config{})

				testutils.TestConfigGetEnumsCount(&Config{}, 0)
				testutils.TestConfigGetFieldsCount(&Config{}, 5)
			})
		})
		ginkgo.Describe("the service instance", func() {
			ginkgo.BeforeEach(func() {
				httpmock.Activate()
			})
			ginkgo.AfterEach(func() {
				httpmock.DeactivateAndReset()
			})
			ginkgo.It("should implement basic service API methods correctly", func() {
				serviceURL := testutils.URLMust("mattermost://mockhost/mocktoken")
				gomega.Expect(service.Initialize(serviceURL, testutils.TestLogger())).To(gomega.Succeed())
				testutils.TestServiceSetInvalidParamValue(service, "foo", "bar")
			})
		})
	})

	ginkgo.It("should return the correct service ID", func() {
		service := &Service{}
		gomega.Expect(service.GetID()).To(gomega.Equal("mattermost"))
	})
})
