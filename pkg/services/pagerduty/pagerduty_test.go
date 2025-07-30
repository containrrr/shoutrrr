package pagerduty

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/types"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	mockIntegrationKey = "eb243592-faa2-4ba2-a551q-1afdf565c889"
	mockHost           = "events.pagerduty.com"
)

func TestPagerDuty(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Pagerduty Suite")
}

var _ = Describe("the PagerDuty service", func() {
	var (
		// a simulated http server to mock out PagerDuty itself
		mockServer *httptest.Server
		// the host of our mock server
		mockHost string
		// function to check if the http request received by the mock server is as expected
		checkRequest func(body string, header http.Header)
		// the shoutrrr PagerDuty service
		service *Service
		// just a mock logger
		mockLogger *log.Logger
	)

	BeforeEach(func() {
		// Initialize a mock http server
		httpHandler := func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			Expect(err).To(BeNil())
			defer r.Body.Close()

			checkRequest(string(body), r.Header)
		}
		mockServer = httptest.NewTLSServer(http.HandlerFunc(httpHandler))

		// Our mock server doesn't have a valid cert
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		// Determine the host of our mock http server
		mockServerURL, err := url.Parse(mockServer.URL)
		Expect(err).To(BeNil())
		mockHost = mockServerURL.Host

		// Initialize a mock logger
		var buf bytes.Buffer
		mockLogger = log.New(&buf, "", 0)
	})

	AfterEach(func() {
		mockServer.Close()
	})

	Context("without query parameters", func() {
		BeforeEach(func() {
			// Initialize service
			serviceURL, err := url.Parse(fmt.Sprintf("pagerduty://%s/%s", mockHost, mockIntegrationKey))
			Expect(err).To(BeNil())

			service = &Service{}
			err = service.Initialize(serviceURL, mockLogger)
			Expect(err).To(BeNil())
		})

		When("sending a simple alert", func() {
			It("should send a request to our mock PagerDuty server with default values", func() {
				checkRequest = func(body string, header http.Header) {
					Expect(header["Content-Type"][0]).To(Equal("application/json"))
					Expect(body).To(Equal(`{` +
						`"payload":{` +
						`"summary":"hello world",` +
						`"severity":"error",` +
						`"source":"default"` +
						`},` +
						`"routing_key":"eb243592-faa2-4ba2-a551q-1afdf565c889",` +
						`"event_action":"trigger"` +
						`}`))
				}

				err := service.Send("hello world", &types.Params{})
				Expect(err).To(BeNil())
			})
		})

	})

	Context("with query parameters", func() {
		BeforeEach(func() {
			// Initialize service
			serviceURL, err := url.Parse(fmt.Sprintf(`pagerduty://%s/%s?severity=critical&source=beszel&action=resolve`, mockHost, mockIntegrationKey))
			Expect(err).To(BeNil())

			service = &Service{}
			err = service.Initialize(serviceURL, mockLogger)
			Expect(err).To(BeNil())
		})

		When("sending a simple alert", func() {
			It("should send a request to our mock PagerDuty server with all fields populated from query parameters", func() {
				checkRequest = func(body string, header http.Header) {
					Expect(header["Content-Type"][0]).To(Equal("application/json"))
					Expect(body).To(Equal(`{` +
						`"payload":{` +
						`"summary":"An example alert message",` +
						`"severity":"critical",` +
						`"source":"beszel"` +
						`},` +
						`"routing_key":"eb243592-faa2-4ba2-a551q-1afdf565c889",` +
						`"event_action":"resolve"` +
						`}`))
				}

				err := service.Send("An example alert message", &types.Params{})
				Expect(err).To(BeNil())
			})
		})
	})
})

var _ = Describe("the PagerDuty Config struct", func() {
	When("generating a config from a simple URL", func() {
		It("should populate the config with host and integration key", func() {
			url, err := url.Parse(fmt.Sprintf("pagerduty://%s/%s", mockHost, mockIntegrationKey))
			Expect(err).To(BeNil())

			config := Config{}
			err = config.SetURL(url)
			Expect(err).To(BeNil())

			Expect(config.IntegrationKey).To(Equal(mockIntegrationKey))
			Expect(config.Host).To(Equal(mockHost))
		})
	})

	When("generating a config from a url with port", func() {
		It("should populate the port field", func() {
			url, err := url.Parse(fmt.Sprintf("pagerduty://%s:12345/%s", mockHost, mockIntegrationKey))
			Expect(err).To(BeNil())

			config := Config{}
			err = config.SetURL(url)
			Expect(err).To(BeNil())

			Expect(config.Port).To(Equal(uint16(12345)))
		})
	})

	When("generating a config from a url with query parameters", func() {
		It("should populate the config fields with the query parameter values", func() {
			queryParams := `severity=critical&source=my-app&action=trigger`
			url, err := url.Parse(fmt.Sprintf("pagerduty://%s:12345/%s?%s", mockHost, mockIntegrationKey, queryParams))
			Expect(err).To(BeNil())

			config := Config{}
			err = config.SetURL(url)
			Expect(err).To(BeNil())

			Expect(config.Severity).To(Equal("critical"))
			Expect(config.Source).To(Equal("my-app"))
			Expect(config.Action).To(Equal("trigger"))
		})
	})

	When("generating a config from a url with differently escaped spaces", func() {
		It("should parse the escaped spaces correctly", func() {
			// Use: '%20', '+' and a normal space
			queryParams := `source=my app`
			url, err := url.Parse(fmt.Sprintf("pagerduty://%s:12345/%s?%s", mockHost, mockIntegrationKey, queryParams))
			Expect(err).To(BeNil())

			config := Config{}
			err = config.SetURL(url)
			Expect(err).To(BeNil())

			Expect(config.Source).To(Equal("my app"))

		})
	})

})

var _ = Describe("the PagerDuty Config default url values", func() {
	When("called", func() {
		It("should extract the default values for host and port", func() {

			config := Config{}
			cfg := reflect.TypeOf(config)
			defaultUrlValues := getDefaultUrlValues(cfg)

			Expect(defaultUrlValues["Host"]).To(Equal("events.pagerduty.com"))
			Expect(defaultUrlValues["Port"]).To(Equal("443"))
		})
	})

})
