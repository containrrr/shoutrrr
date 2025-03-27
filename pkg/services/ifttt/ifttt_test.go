package ifttt_test

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/nicholas-fedor/shoutrrr/internal/testutils"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/ifttt"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// TestIFTTT runs the Ginkgo test suite for the IFTTT package.
func TestIFTTT(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Shoutrrr IFTTT Suite")
}

var (
	service    *ifttt.Service
	logger     *log.Logger
	envTestURL string
	_          = ginkgo.BeforeSuite(func() {
		service = &ifttt.Service{}
		logger = testutils.TestLogger()
		envTestURL = os.Getenv("SHOUTRRR_IFTTT_URL")
	})
)

var _ = ginkgo.Describe("the IFTTT service", func() {
	ginkgo.When("running integration tests", func() {
		ginkgo.It("sends a message successfully with a valid ENV URL", func() {
			if envTestURL == "" {
				ginkgo.Skip("No integration test ENV URL was set")

				return
			}
			serviceURL := testutils.URLMust(envTestURL)
			err := service.Initialize(serviceURL, logger)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			err = service.Send("This is an integration test", nil)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
	})

	ginkgo.Describe("the service", func() {
		ginkgo.BeforeEach(func() {
			service = &ifttt.Service{}
			service.SetLogger(logger)
		})
		ginkgo.It("returns the correct service identifier", func() {
			gomega.Expect(service.GetID()).To(gomega.Equal("ifttt"))
		})
	})

	ginkgo.When("parsing the configuration URL", func() {
		ginkgo.BeforeEach(func() {
			service = &ifttt.Service{}
			service.SetLogger(logger)
		})
		ginkgo.It("returns an error if no arguments are supplied", func() {
			serviceURL := testutils.URLMust("ifttt://")
			err := service.Initialize(serviceURL, logger)
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
		ginkgo.It("returns an error if no webhook ID is given", func() {
			serviceURL := testutils.URLMust("ifttt:///?events=event1")
			err := service.Initialize(serviceURL, logger)
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
		ginkgo.It("returns an error if no events are given", func() {
			serviceURL := testutils.URLMust("ifttt://dummyID")
			err := service.Initialize(serviceURL, logger)
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
		ginkgo.It("returns an error when an invalid query key is given", func() { // Line 54
			serviceURL := testutils.URLMust("ifttt://dummyID/?events=event1&badquery=foo")
			err := service.Initialize(serviceURL, logger)
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
		ginkgo.It("returns an error if message value is above 3", func() {
			serviceURL := testutils.URLMust("ifttt://dummyID/?events=event1&messagevalue=8")
			err := service.Initialize(serviceURL, logger)
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
		ginkgo.It("returns an error if message value is below 1", func() { // Line 60
			serviceURL := testutils.URLMust("ifttt://dummyID/?events=event1&messagevalue=0")
			err := service.Initialize(serviceURL, logger)
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
		ginkgo.It("does not return an error if webhook ID and at least one event are given", func() {
			serviceURL := testutils.URLMust("ifttt://dummyID/?events=event1")
			err := service.Initialize(serviceURL, logger)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
		ginkgo.It("returns an error if titlevalue is invalid", func() { // Line 78
			serviceURL := testutils.URLMust("ifttt://dummyID/?events=event1&titlevalue=4")
			err := service.Initialize(serviceURL, logger)
			gomega.Expect(err).To(gomega.MatchError("invalid value for titlevalue: only values 1-3 or 0 (for disabling) are supported"))
		})
		ginkgo.It("returns an error if titlevalue equals messagevalue", func() { // Line 82
			serviceURL := testutils.URLMust("ifttt://dummyID/?events=event1&messagevalue=2&titlevalue=2")
			err := service.Initialize(serviceURL, logger)
			gomega.Expect(err).To(gomega.MatchError("titlevalue cannot use the same number as messagevalue"))
		})
	})

	ginkgo.When("serializing a config to URL", func() {
		ginkgo.BeforeEach(func() {
			service = &ifttt.Service{}
			service.SetLogger(logger)
		})
		ginkgo.When("given multiple events", func() {
			ginkgo.It("returns an URL with all events comma-separated", func() {
				configURL := testutils.URLMust("ifttt://dummyID/?events=foo%2Cbar%2Cbaz")
				err := service.Initialize(configURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				resultURL := service.Config.GetURL().String()
				gomega.Expect(resultURL).To(gomega.Equal(configURL.String()))
			})
		})
		ginkgo.When("given values", func() {
			ginkgo.It("returns an URL with all values", func() {
				configURL := testutils.URLMust("ifttt://dummyID/?events=event1&value1=v1&value2=v2&value3=v3")
				err := service.Initialize(configURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				resultURL := service.Config.GetURL().String()
				gomega.Expect(resultURL).To(gomega.Equal(configURL.String()))
			})
		})
	})

	ginkgo.When("sending a message", func() {
		ginkgo.BeforeEach(func() {
			httpmock.Activate()
			service = &ifttt.Service{}
			service.SetLogger(logger)
		})
		ginkgo.AfterEach(func() {
			httpmock.DeactivateAndReset()
		})
		ginkgo.It("errors if the response code is not 200-299", func() {
			configURL := testutils.URLMust("ifttt://dummy/?events=foo")
			err := service.Initialize(configURL, logger)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			httpmock.RegisterResponder(
				"POST",
				"https://maker.ifttt.com/trigger/foo/with/key/dummy",
				httpmock.NewStringResponder(404, ""),
			)
			err = service.Send("hello", nil)
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
		ginkgo.It("does not error if the response code is 200", func() {
			configURL := testutils.URLMust("ifttt://dummy/?events=foo")
			err := service.Initialize(configURL, logger)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			httpmock.RegisterResponder(
				"POST",
				"https://maker.ifttt.com/trigger/foo/with/key/dummy",
				httpmock.NewStringResponder(200, ""),
			)
			err = service.Send("hello", nil)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
		ginkgo.It("returns an error if params update fails", func() { // Line 55
			configURL := testutils.URLMust("ifttt://dummy/?events=event1")
			err := service.Initialize(configURL, logger)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			params := types.Params{"messagevalue": "invalid"}
			err = service.Send("hello", &params)
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
		ginkgo.DescribeTable("sets message to correct value field based on messagevalue",
			func(messageValue int, expectedField string) { // Lines 30, 32, 34
				configURL := testutils.URLMust(fmt.Sprintf("ifttt://dummy/?events=event1&messagevalue=%d", messageValue))
				err := service.Initialize(configURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				httpmock.RegisterResponder(
					"POST",
					"https://maker.ifttt.com/trigger/event1/with/key/dummy",
					func(req *http.Request) (*http.Response, error) {
						body, err := io.ReadAll(req.Body)
						gomega.Expect(err).NotTo(gomega.HaveOccurred())
						var payload jsonPayload
						err = json.Unmarshal(body, &payload)
						gomega.Expect(err).NotTo(gomega.HaveOccurred())
						switch expectedField {
						case "Value1":
							gomega.Expect(payload.Value1).To(gomega.Equal("hello"))
							gomega.Expect(payload.Value2).To(gomega.Equal(""))
							gomega.Expect(payload.Value3).To(gomega.Equal(""))
						case "Value2":
							gomega.Expect(payload.Value1).To(gomega.Equal(""))
							gomega.Expect(payload.Value2).To(gomega.Equal("hello"))
							gomega.Expect(payload.Value3).To(gomega.Equal(""))
						case "Value3":
							gomega.Expect(payload.Value1).To(gomega.Equal(""))
							gomega.Expect(payload.Value2).To(gomega.Equal(""))
							gomega.Expect(payload.Value3).To(gomega.Equal("hello"))
						}

						return httpmock.NewStringResponse(200, ""), nil
					},
				)
				err = service.Send("hello", nil)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			},
			ginkgo.Entry("messagevalue=1 sets Value1", 1, "Value1"),
			ginkgo.Entry("messagevalue=2 sets Value2", 2, "Value2"),
			ginkgo.Entry("messagevalue=3 sets Value3", 3, "Value3"),
		)
		ginkgo.It("overrides Value2 with params when messagevalue is 1", func() { // Line 36
			configURL := testutils.URLMust("ifttt://dummy/?events=event1&messagevalue=1")
			err := service.Initialize(configURL, logger)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			httpmock.RegisterResponder(
				"POST",
				"https://maker.ifttt.com/trigger/event1/with/key/dummy",
				func(req *http.Request) (*http.Response, error) {
					body, err := io.ReadAll(req.Body)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					var payload jsonPayload
					err = json.Unmarshal(body, &payload)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(payload.Value1).To(gomega.Equal("hello"))
					gomega.Expect(payload.Value2).To(gomega.Equal("y"))
					gomega.Expect(payload.Value3).To(gomega.Equal(""))

					return httpmock.NewStringResponse(200, ""), nil
				},
			)
			params := types.Params{
				"value2": "y",
			}
			err = service.Send("hello", &params)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
		ginkgo.It("overrides payload values with params", func() { // Lines 17, 21, 25
			configURL := testutils.URLMust("ifttt://dummy/?events=event1&value1=a&value2=b&value3=c&messagevalue=2")
			err := service.Initialize(configURL, logger)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			httpmock.RegisterResponder(
				"POST",
				"https://maker.ifttt.com/trigger/event1/with/key/dummy",
				func(req *http.Request) (*http.Response, error) {
					body, err := io.ReadAll(req.Body)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					var payload jsonPayload
					err = json.Unmarshal(body, &payload)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(payload.Value1).To(gomega.Equal("x"))
					gomega.Expect(payload.Value2).To(gomega.Equal("hello"))
					gomega.Expect(payload.Value3).To(gomega.Equal("z"))

					return httpmock.NewStringResponse(200, ""), nil
				},
			)
			params := types.Params{
				"value1": "x",
				// "value2": "y", // Omitted to let message override
				"value3": "z",
			}
			err = service.Send("hello", &params)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
		ginkgo.It("fails with multiple events when one errors", func() { // Line 82
			configURL := testutils.URLMust("ifttt://dummy/?events=event1%2Cevent2")
			err := service.Initialize(configURL, logger)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			httpmock.RegisterResponder(
				"POST",
				"https://maker.ifttt.com/trigger/event1/with/key/dummy",
				httpmock.NewStringResponder(200, ""),
			)
			httpmock.RegisterResponder(
				"POST",
				"https://maker.ifttt.com/trigger/event2/with/key/dummy",
				httpmock.NewStringResponder(404, ""),
			)
			err = service.Send("hello", nil)
			gomega.Expect(err).To(gomega.MatchError("failed to send IFTTT event \"event2\": got response status code 404"))
		})
		ginkgo.It("fails with network error", func() { // Line 95
			configURL := testutils.URLMust("ifttt://dummy/?events=event1")
			err := service.Initialize(configURL, logger)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			httpmock.RegisterResponder(
				"POST",
				"https://maker.ifttt.com/trigger/event1/with/key/dummy",
				httpmock.NewErrorResponder(fmt.Errorf("network failure")),
			)
			err = service.Send("hello", nil)
			gomega.Expect(err).To(gomega.MatchError("failed to send IFTTT event \"event1\": Post \"https://maker.ifttt.com/trigger/event1/with/key/dummy\": network failure"))
		})
	})
})

type jsonPayload struct {
	Value1 string `json:"value1"`
	Value2 string `json:"value2"`
	Value3 string `json:"value3"`
}
