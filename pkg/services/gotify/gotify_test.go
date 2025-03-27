package gotify_test

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/nicholas-fedor/shoutrrr/internal/testutils"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/gotify"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// Test constants.
const (
	TargetURL = "https://my.gotify.tld/message?token=Aaa.bbb.ccc.ddd"
)

// TestGotify runs the Ginkgo test suite for the Gotify package.
func TestGotify(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Shoutrrr Gotify Suite")
}

var (
	service      *gotify.Service
	logger       *log.Logger
	envGotifyURL *url.URL
	_            = ginkgo.BeforeSuite(func() {
		service = &gotify.Service{}
		logger = log.New(ginkgo.GinkgoWriter, "Test", log.LstdFlags)
		var err error
		envGotifyURL, err = url.Parse(os.Getenv("SHOUTRRR_GOTIFY_URL"))
		if err != nil {
			envGotifyURL = &url.URL{} // Default to empty URL if parsing fails
		}
	})
)

var _ = ginkgo.Describe("the Gotify service", func() {
	ginkgo.When("running integration tests", func() {
		ginkgo.It("sends a message successfully with a valid ENV URL", func() {
			if envGotifyURL.String() == "" {
				ginkgo.Skip("No integration test ENV URL was set")

				return
			}
			serviceURL := testutils.URLMust(envGotifyURL.String())
			err := service.Initialize(serviceURL, testutils.TestLogger())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			err = service.Send("This is an integration test message", nil)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
	})

	ginkgo.Describe("the service", func() {
		ginkgo.BeforeEach(func() {
			service = &gotify.Service{}
			service.SetLogger(logger)
		})
		ginkgo.It("returns the correct service identifier", func() {
			gomega.Expect(service.GetID()).To(gomega.Equal("gotify"))
		})
	})

	ginkgo.When("parsing the configuration URL", func() {
		ginkgo.BeforeEach(func() {
			service = &gotify.Service{}
			service.SetLogger(logger)
		})
		ginkgo.It("builds a valid Gotify URL without path", func() {
			configURL := testutils.URLMust("gotify://my.gotify.tld/Aaa.bbb.ccc.ddd")
			err := service.Initialize(configURL, logger)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(service.Config.GetURL().String()).To(gomega.Equal(configURL.String()))
		})
		ginkgo.When("TLS is disabled", func() {
			ginkgo.It("uses http scheme", func() {
				configURL := testutils.URLMust("gotify://my.gotify.tld/Aaa.bbb.ccc.ddd?disabletls=yes")
				err := service.Initialize(configURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(service.Config.DisableTLS).To(gomega.BeTrue())
			})
		})
		ginkgo.When("a custom path is provided", func() {
			ginkgo.It("includes the path in the URL", func() {
				configURL := testutils.URLMust("gotify://my.gotify.tld/gotify/Aaa.bbb.ccc.ddd")
				err := service.Initialize(configURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(service.Config.GetURL().String()).To(gomega.Equal(configURL.String()))
			})
		})
		ginkgo.When("the token has an invalid length", func() {
			ginkgo.It("reports an error during send", func() {
				configURL := testutils.URLMust("gotify://my.gotify.tld/short") // Length < 15
				err := service.Initialize(configURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				err = service.Send("Message", nil)
				gomega.Expect(err).To(gomega.MatchError("invalid gotify token \"short\""))
			})
		})
		ginkgo.When("the token has an invalid prefix", func() {
			ginkgo.It("reports an error during send", func() {
				configURL := testutils.URLMust("gotify://my.gotify.tld/Chwbsdyhwwgarxd") // Starts with 'C', not 'A'
				err := service.Initialize(configURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				err = service.Send("Message", nil)
				gomega.Expect(err).To(gomega.MatchError("invalid gotify token \"Chwbsdyhwwgarxd\""))
			})
		})
		ginkgo.It("is identical after de-/serialization with path", func() {
			testURL := "gotify://my.gotify.tld/gotify/Aaa.bbb.ccc.ddd?title=Test+title"
			serviceURL := testutils.URLMust(testURL)
			err := service.Initialize(serviceURL, logger)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(service.Config.GetURL().String()).To(gomega.Equal(testURL))
		})
		ginkgo.It("is identical after de-/serialization without path", func() {
			testURL := "gotify://my.gotify.tld/Aaa.bbb.ccc.ddd?disabletls=Yes&priority=1&title=Test+title"
			serviceURL := testutils.URLMust(testURL)
			err := service.Initialize(serviceURL, logger)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(service.Config.GetURL().String()).To(gomega.Equal(testURL))
		})
		ginkgo.It("allows slash at the end of the token", func() {
			configURL := testutils.URLMust("gotify://my.gotify.tld/Aaa.bbb.ccc.ddd/")
			err := service.Initialize(configURL, logger)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(service.Config.Token).To(gomega.Equal("Aaa.bbb.ccc.ddd"))
		})
		ginkgo.It("allows slash at the end of the token with additional path", func() {
			configURL := testutils.URLMust("gotify://my.gotify.tld/path/to/gotify/Aaa.bbb.ccc.ddd/")
			err := service.Initialize(configURL, logger)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(service.Config.Token).To(gomega.Equal("Aaa.bbb.ccc.ddd"))
		})
		ginkgo.It("does not crash on empty token or path slash", func() {
			configURL := testutils.URLMust("gotify://my.gotify.tld//")
			err := service.Initialize(configURL, logger)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(service.Config.Token).To(gomega.Equal(""))
		})
	})

	ginkgo.When("the token contains invalid characters", func() {
		ginkgo.It("reports an error during send", func() {
			configURL := testutils.URLMust("gotify://my.gotify.tld/Aaa.bbb.ccc.dd!")
			err := service.Initialize(configURL, logger)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			err = service.Send("Message", nil)
			gomega.Expect(err).To(gomega.MatchError("invalid gotify token \"Aaa.bbb.ccc.dd!\""))
		})
	})

	ginkgo.Describe("sending the payload", func() {
		ginkgo.BeforeEach(func() {
			service = &gotify.Service{}
			service.SetLogger(logger)
			configURL := testutils.URLMust("gotify://my.gotify.tld/Aaa.bbb.ccc.ddd")
			err := service.Initialize(configURL, logger)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			httpmock.ActivateNonDefault(service.GetHTTPClient())
			httpmock.Activate()
		})
		ginkgo.AfterEach(func() {
			httpmock.DeactivateAndReset()
		})
		ginkgo.When("sending via webhook URL", func() {
			ginkgo.It("does not report an error if the server accepts the payload", func() {
				httpmock.RegisterResponder(
					"POST",
					TargetURL,
					testutils.JSONRespondMust(200, map[string]interface{}{
						"id":       float64(1),
						"appid":    float64(1),
						"message":  "Message",
						"title":    "Shoutrrr notification",
						"priority": float64(0),
						"date":     "2023-01-01T00:00:00Z",
					}),
				)
				err := service.Send("Message", nil)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
			ginkgo.It("reports an error if the server rejects the payload with an error response", func() {
				httpmock.RegisterResponder(
					"POST",
					TargetURL,
					testutils.JSONRespondMust(401, map[string]interface{}{
						"error":            "Unauthorized",
						"errorCode":        float64(401),
						"errorDescription": "you need to provide a valid access token or user credentials to access this api",
					}),
				)
				err := service.Send("Message", nil)
				gomega.Expect(err).To(gomega.MatchError("server respondend with Unauthorized (401): you need to provide a valid access token or user credentials to access this api"))
			})
			ginkgo.It("reports an error if sending fails with a network error", func() {
				httpmock.RegisterResponder(
					"POST",
					TargetURL,
					httpmock.NewErrorResponder(fmt.Errorf("network failure")),
				)
				err := service.Send("Message", nil)
				gomega.Expect(err).To(gomega.MatchError("failed to send notification to Gotify: error sending payload: Post \"https://my.gotify.tld/message?token=Aaa.bbb.ccc.ddd\": network failure"))
			})
			ginkgo.It("logs an error if params update fails", func() {
				var logBuffer bytes.Buffer
				service.SetLogger(log.New(&logBuffer, "Test", log.LstdFlags))
				httpmock.RegisterResponder(
					"POST",
					TargetURL,
					testutils.JSONRespondMust(200, map[string]interface{}{
						"id":       float64(1),
						"appid":    float64(1),
						"message":  "Message",
						"title":    "Shoutrrr notification",
						"priority": float64(0),
						"date":     "2023-01-01T00:00:00Z",
					}),
				)
				params := types.Params{"priority": "invalid"}
				err := service.Send("Message", &params)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(logBuffer.String()).To(gomega.ContainSubstring("Failed to update params"))
			})
		})
	})
})
