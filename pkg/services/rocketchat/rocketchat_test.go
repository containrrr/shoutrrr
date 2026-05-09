package rocketchat

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/containrrr/shoutrrr/internal/testutils"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	service          *Service
	envRocketchatURL *url.URL
	_                = BeforeSuite(func() {
		service = &Service{}
		envRocketchatURL, _ = url.Parse(os.Getenv("SHOUTRRR_ROCKETCHAT_URL"))
	})
)

func TestRocketchat(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Rocketchat Suite")
}

var _ = Describe("the rocketchat service", func() {

	When("running integration tests", func() {
		It("should work without errors", func() {
			if envRocketchatURL.String() == "" {
				return
			}
			serviceURL, _ := url.Parse(envRocketchatURL.String())
			service.Initialize(serviceURL, testutils.TestLogger())
			err := service.Send(
				"this is an integration test",
				nil,
			)
			Expect(err).NotTo(HaveOccurred())
		})
	})
	Describe("the rocketchat config", func() {
		When("generating a config object", func() {
			rocketchatURL, _ := url.Parse("rocketchat://rocketchat.my-domain.com/tokenA/tokenB")
			config := &Config{}
			err := config.SetURL(rocketchatURL)
			It("should not have caused an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should set host", func() {
				Expect(config.Host).To(Equal("rocketchat.my-domain.com"))
			})
			It("should set token A", func() {
				Expect(config.TokenA).To(Equal("tokenA"))
			})
			It("should set token B", func() {
				Expect(config.TokenB).To(Equal("tokenB"))
			})
			It("should not set channel or username", func() {
				Expect(config.Channel).To(BeEmpty())
				Expect(config.UserName).To(BeEmpty())
			})
			It("should generate a URL without an empty port", func() {
				Expect(config.GetURL().String()).To(Equal("rocketchat://rocketchat.my-domain.com/tokenA/tokenB"))
			})
		})
		When("generating a config object with a port", func() {
			rocketchatURL, _ := url.Parse("rocketchat://rocketchat.my-domain.com:5055/tokenA/tokenB")
			config := &Config{}
			err := config.SetURL(rocketchatURL)
			It("should not have caused an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should preserve the port in the generated URL", func() {
				Expect(config.GetURL().String()).To(Equal("rocketchat://rocketchat.my-domain.com:5055/tokenA/tokenB"))
			})
		})
		When("generating a new config with url, that has no token", func() {
			rocketchatURL, _ := url.Parse("rocketchat://rocketchat.my-domain.com")
			config := &Config{}
			err := config.SetURL(rocketchatURL)
			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
		When("generating a config object with username only", func() {
			rocketchatURL, _ := url.Parse("rocketchat://testUserName@rocketchat.my-domain.com/tokenA/tokenB")
			config := &Config{}
			err := config.SetURL(rocketchatURL)
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
			rocketchatURL, _ := url.Parse("rocketchat://rocketchat.my-domain.com/tokenA/tokenB/testChannel")
			config := &Config{}
			err := config.SetURL(rocketchatURL)
			It("should not hav caused an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should set channel", func() {
				Expect(config.Channel).To(Equal("#testChannel"))
			})
			It("should not set username", func() {
				Expect(config.UserName).To(BeEmpty())
			})
		})
		When("generating a config object with channel and userName", func() {
			rocketchatURL, _ := url.Parse("rocketchat://testUserName@rocketchat.my-domain.com/tokenA/tokenB/testChannel")
			config := &Config{}
			err := config.SetURL(rocketchatURL)
			It("should not hav caused an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should set channel", func() {
				Expect(config.Channel).To(Equal("#testChannel"))
			})
			It("should set username", func() {
				Expect(config.UserName).To(Equal("testUserName"))
			})
		})
		When("generating a config object with user and userName", func() {
			rocketchatURL, _ := url.Parse("rocketchat://testUserName@rocketchat.my-domain.com/tokenA/tokenB/@user")
			config := &Config{}
			err := config.SetURL(rocketchatURL)
			It("should not hav caused an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should set channel", func() {
				Expect(config.Channel).To(Equal("@user"))
			})
			It("should set username", func() {
				Expect(config.UserName).To(Equal("testUserName"))
			})
		})
	})
	Describe("Sending messages", func() {
		BeforeEach(func() {
			httpmock.Activate()
		})

		AfterEach(func() {
			httpmock.DeactivateAndReset()
		})

		When("sending a message completely without parameters", func() {
			rocketchatURL, _ := url.Parse("rocketchat://rocketchat.my-domain.com/tokenA/tokenB")
			config := &Config{}
			config.SetURL(rocketchatURL)
			It("should generate the correct url to call", func() {
				generatedURL := buildURL(config)
				Expect(generatedURL).To(Equal("https://rocketchat.my-domain.com/hooks/tokenA/tokenB"))
			})
			It("should generate the correct JSON body", func() {
				json, err := CreateJSONPayload(config, "this is a message", nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(json)).To(Equal("{\"text\":\"this is a message\"}"))
			})
		})
		When("sending a message with pre set username and channel", func() {
			rocketchatURL, _ := url.Parse("rocketchat://testUserName@rocketchat.my-domain.com/tokenA/tokenB/testChannel")
			config := &Config{}
			config.SetURL(rocketchatURL)
			It("should generate the correct JSON body", func() {
				json, err := CreateJSONPayload(config, "this is a message", nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(json)).To(Equal("{\"text\":\"this is a message\",\"username\":\"testUserName\",\"channel\":\"#testChannel\"}"))
			})
		})
		When("sending a message with pre set username and channel but overwriting them with parameters", func() {
			rocketchatURL, _ := url.Parse("rocketchat://testUserName@rocketchat.my-domain.com/tokenA/tokenB/testChannel")
			config := &Config{}
			config.SetURL(rocketchatURL)
			It("should generate the correct JSON body", func() {
				params := (*types.Params)(&map[string]string{"username": "overwriteUserName", "channel": "overwriteChannel"})
				json, err := CreateJSONPayload(config, "this is a message", params)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(json)).To(Equal("{\"text\":\"this is a message\",\"username\":\"overwriteUserName\",\"channel\":\"overwriteChannel\"}"))
			})
		})
		When("sending to an URL which contains HOST:PORT", func() {
			rocketchatURL, _ := url.Parse("rocketchat://testUserName@rocketchat.my-domain.com:5055/tokenA/tokenB/testChannel")
			config := &Config{}
			config.SetURL(rocketchatURL)
			It("should generate a correct hook URL https://HOST:PORT", func() {
				hookURL := buildURL(config)
				Expect(hookURL).To(ContainSubstring("my-domain.com:5055"))
			})
		})
		When("sending to an URL with badly syntaxed #channel name", func() {
			It("should properly parse the Channel", func() {
				rocketchatURL, _ := url.Parse("rocketchat://testUserName@rocketchat.my-domain.com:5055/tokenA/tokenB/###########################testChannel")
				config := &Config{}
				config.SetURL(rocketchatURL)
				Expect(config.Channel).To(ContainSubstring("###########################testChannel"))
			})
			It("should properly parse the Channel", func() {
				rocketchatURL, _ := url.Parse("rocketchat://testUserName@rocketchat.my-domain.com:5055/tokenA/tokenB/#testChannel")
				config := &Config{}
				config.SetURL(rocketchatURL)
				Expect(config.Channel).To(ContainSubstring("#testChannel"))
			})
		})
		When("the webhook accepts the notification", func() {
			It("should post the expected JSON payload", func() {
				serviceURL, _ := url.Parse("rocketchat://testUserName@rocketchat.my-domain.com/tokenA/tokenB/testChannel")
				hookURL := "https://rocketchat.my-domain.com/hooks/tokenA/tokenB"

				service := &Service{}
				err := service.Initialize(serviceURL, testutils.TestLogger())
				Expect(err).NotTo(HaveOccurred())

				httpmock.RegisterResponder(http.MethodPost, hookURL, func(req *http.Request) (*http.Response, error) {
					body, err := io.ReadAll(req.Body)
					Expect(err).NotTo(HaveOccurred())
					Expect(req.Header.Get("Content-Type")).To(Equal("application/json"))
					Expect(string(body)).To(MatchJSON(`{
						"text": "this is a message",
						"username": "testUserName",
						"channel": "#testChannel"
					}`))

					return httpmock.NewStringResponse(http.StatusOK, ""), nil
				})

				err = service.Send("this is a message", nil)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		When("the webhook rejects the notification", func() {
			It("should include the response body in the error", func() {
				serviceURL, _ := url.Parse("rocketchat://rocketchat.my-domain.com/tokenA/tokenB")
				hookURL := "https://rocketchat.my-domain.com/hooks/tokenA/tokenB"

				service := &Service{}
				err := service.Initialize(serviceURL, testutils.TestLogger())
				Expect(err).NotTo(HaveOccurred())

				httpmock.RegisterResponder(http.MethodPost, hookURL, httpmock.NewStringResponder(http.StatusBadRequest, "bad payload"))

				err = service.Send("this is a message", nil)
				Expect(err).To(MatchError("notification failed: 400 bad payload"))
			})
		})
		When("the webhook cannot be reached", func() {
			It("should report the transport error with host and port context", func() {
				serviceURL, _ := url.Parse("rocketchat://rocketchat.my-domain.com:5055/tokenA/tokenB")
				hookURL := "https://rocketchat.my-domain.com:5055/hooks/tokenA/tokenB"

				service := &Service{}
				err := service.Initialize(serviceURL, testutils.TestLogger())
				Expect(err).NotTo(HaveOccurred())

				httpmock.RegisterResponder(http.MethodPost, hookURL, httpmock.NewErrorResponder(errors.New("network down")))

				err = service.Send("this is a message", nil)
				Expect(err).To(MatchError(ContainSubstring("network down")))
				Expect(err).To(MatchError(ContainSubstring("HOST: rocketchat.my-domain.com")))
				Expect(err).To(MatchError(ContainSubstring("PORT: 5055")))
			})
		})
	})
})
