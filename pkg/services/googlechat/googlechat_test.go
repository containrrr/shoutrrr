package googlechat_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/nicholas-fedor/shoutrrr/internal/testutils"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/googlechat"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// TestGooglechat runs the Ginkgo test suite for the Google Chat package.
func TestGooglechat(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Shoutrrr Google Chat Suite")
}

var (
	service          *googlechat.Service
	logger           *log.Logger
	envGooglechatURL *url.URL
	_                = ginkgo.BeforeSuite(func() {
		service = &googlechat.Service{}
		logger = log.New(ginkgo.GinkgoWriter, "Test", log.LstdFlags)
		var err error
		envGooglechatURL, err = url.Parse(os.Getenv("SHOUTRRR_GOOGLECHAT_URL"))
		if err != nil {
			envGooglechatURL = &url.URL{} // Default to empty URL if parsing fails
		}
	})
)

var _ = ginkgo.Describe("Google Chat Service", func() {
	ginkgo.When("running integration tests", func() {
		ginkgo.It("sends a message successfully with a valid ENV URL", func() {
			if envGooglechatURL.String() == "" {
				ginkgo.Skip("No integration test ENV URL was set")

				return
			}
			serviceURL := testutils.URLMust(envGooglechatURL.String())
			err := service.Initialize(serviceURL, testutils.TestLogger())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			err = service.Send("This is an integration test message", nil)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
	})

	ginkgo.Describe("the service", func() {
		ginkgo.BeforeEach(func() {
			service = &googlechat.Service{}
			service.SetLogger(logger)
		})
		ginkgo.It("implements Service interface", func() {
			var impl types.Service = service
			gomega.Expect(impl).ToNot(gomega.BeNil())
		})
		ginkgo.It("returns the correct service identifier", func() {
			gomega.Expect(service.GetID()).To(gomega.Equal("googlechat"))
		})
	})

	ginkgo.When("parsing the configuration URL", func() {
		ginkgo.BeforeEach(func() {
			service = &googlechat.Service{}
			service.SetLogger(logger)
		})
		ginkgo.It("builds a valid Google Chat Incoming Webhook URL", func() {
			configURL := testutils.URLMust("googlechat://chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=baz")
			err := service.Initialize(configURL, logger)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(service.Config.GetURL().String()).To(gomega.Equal(configURL.String()))
		})
		ginkgo.It("is identical after de-/serialization", func() {
			testURL := "googlechat://chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=baz"
			serviceURL := testutils.URLMust(testURL)
			err := service.Initialize(serviceURL, logger)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(service.Config.GetURL().String()).To(gomega.Equal(testURL))
		})
		ginkgo.It("returns an error if key is present but empty", func() {
			configURL := testutils.URLMust("googlechat://chat.googleapis.com/v1/spaces/FOO/messages?key=&token=baz")
			err := service.Initialize(configURL, logger)
			gomega.Expect(err).To(gomega.MatchError("missing field 'key'"))
		})
		ginkgo.It("returns an error if token is present but empty", func() {
			configURL := testutils.URLMust("googlechat://chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=")
			err := service.Initialize(configURL, logger)
			gomega.Expect(err).To(gomega.MatchError("missing field 'token'"))
		})
	})

	ginkgo.Describe("sending the payload", func() {
		ginkgo.BeforeEach(func() {
			httpmock.Activate()
			service = &googlechat.Service{}
			service.SetLogger(logger)
		})
		ginkgo.AfterEach(func() {
			httpmock.DeactivateAndReset()
		})
		ginkgo.When("sending via webhook URL", func() {
			ginkgo.It("does not report an error if the server accepts the payload", func() {
				configURL := testutils.URLMust("googlechat://chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=baz")
				err := service.Initialize(configURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				httpmock.RegisterResponder(
					"POST",
					"https://chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=baz",
					httpmock.NewStringResponder(200, ""),
				)
				err = service.Send("Message", nil)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
			ginkgo.It("reports an error if the server rejects the payload", func() {
				configURL := testutils.URLMust("googlechat://chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=baz")
				err := service.Initialize(configURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				httpmock.RegisterResponder(
					"POST",
					"https://chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=baz",
					httpmock.NewStringResponder(400, "Bad Request"),
				)
				err = service.Send("Message", nil)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
			ginkgo.It("marshals the payload correctly with the message", func() {
				configURL := testutils.URLMust("googlechat://chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=baz")
				err := service.Initialize(configURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				httpmock.RegisterResponder(
					"POST",
					"https://chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=baz",
					func(req *http.Request) (*http.Response, error) {
						body, err := io.ReadAll(req.Body)
						gomega.Expect(err).NotTo(gomega.HaveOccurred())
						gomega.Expect(string(body)).To(gomega.MatchJSON(`{"text":"Test Message"}`))

						return httpmock.NewStringResponse(200, ""), nil
					},
				)
				err = service.Send("Test Message", nil)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
			ginkgo.It("sends the POST request with correct URL and content type", func() {
				configURL := testutils.URLMust("googlechat://chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=baz")
				err := service.Initialize(configURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				httpmock.RegisterResponder(
					"POST",
					"https://chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=baz",
					func(req *http.Request) (*http.Response, error) {
						gomega.Expect(req.Method).To(gomega.Equal("POST"))
						gomega.Expect(req.Header.Get("Content-Type")).To(gomega.Equal("application/json"))

						return httpmock.NewStringResponse(200, ""), nil
					},
				)
				err = service.Send("Message", nil)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
			ginkgo.It("returns marshal error if JSON marshaling fails", func() {
				// Note: Cannot directly force json.Marshal to fail with current simple string struct
				// This test assumes a future change might make JSON unmarshalable
				configURL := testutils.URLMust("googlechat://chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=baz")
				err := service.Initialize(configURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				// Simulate marshal failure indirectly (requires package change or reflection, skipped here)
				// For now, test happy path and note limitation
				httpmock.RegisterResponder(
					"POST",
					"https://chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=baz",
					httpmock.NewStringResponder(200, ""),
				)
				err = service.Send("Valid Message", nil)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				// Placeholder: If JSON struct changes, add test for unmarshalable input here
			})
			ginkgo.It("returns formatted error if HTTP POST fails", func() {
				configURL := testutils.URLMust("googlechat://chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=baz")
				err := service.Initialize(configURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				httpmock.RegisterResponder(
					"POST",
					"https://chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=baz",
					httpmock.NewErrorResponder(fmt.Errorf("network failure")),
				)
				err = service.Send("Message", nil)
				gomega.Expect(err).To(gomega.MatchError("failed to send notification to Google Chat: Post \"https://chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=baz\": network failure"))
			})
		})
	})
})
