package opsgenie

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/containrrr/shoutrrr/pkg/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	mockAPIKey = "eb243592-faa2-4ba2-a551q-1afdf565c889"
)

func TestOpsGenie(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr OpsGenie Suite")
}

var _ = Describe("the OpsGenie service", func() {
	var (
		mockServer   *httptest.Server
		service      *Service
		checkRequest func(body string, header http.Header)
	)

	BeforeEach(func() {
		// Initialize a mock http server
		httpHandler := func(w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(r.Body)
			Expect(err).To(BeNil())
			defer r.Body.Close()

			checkRequest(string(body), r.Header)
		}
		mockServer = httptest.NewTLSServer(http.HandlerFunc(httpHandler))

		// Our mock server doesn't have a valid cert
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		// Building a mock URL.
		// It'll look something like: opsgenie://127.0.0.1:63457/eb243592-faa2-4ba2-a551q-1afdf565c889
		mockServerURL, err := url.Parse(mockServer.URL)
		Expect(err).To(BeNil())
		mockURL, err := url.Parse(fmt.Sprintf("opsgenie://%s/%s", mockServerURL.Host, mockAPIKey))
		Expect(err).To(BeNil())

		// Initialize a mock logger
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)

		// Initialize the OpsGenie service
		service = &Service{}
		err = service.Initialize(mockURL, logger)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		mockServer.Close()
	})

	It("should send an alert", func() {
		checkRequest = func(body string, header http.Header) {
			Expect(header["Authorization"][0]).To(Equal("GenieKey " + mockAPIKey))
			Expect(header["Content-Type"][0]).To(Equal("application/json"))
			Expect(body).To(Equal(`{"message":"hello world"}`))
		}

		err := service.Send("hello world", &types.Params{})
		Expect(err).To(BeNil())
	})

	When("provided nil params", func() {
		It("should send an alert without additional fields", func() {
			checkRequest = func(body string, header http.Header) {
				Expect(header["Authorization"][0]).To(Equal("GenieKey " + mockAPIKey))
				Expect(header["Content-Type"][0]).To(Equal("application/json"))
				Expect(body).To(Equal(`{"message":"hello world"}`))
			}

			err := service.Send("hello world", nil)
			Expect(err).To(BeNil())
		})
	})

	When("provided params", func() {
		It("should send an alert with all field", func() {
			checkRequest = func(body string, header http.Header) {
				Expect(header["Authorization"][0]).To(Equal("GenieKey " + mockAPIKey))
				Expect(header["Content-Type"][0]).To(Equal("application/json"))
				Expect(body).To(Equal(`{"` +
					`message":"An example alert message",` +
					`"alias":"Life is too short for no alias",` +
					`"description":"Every alert needs a description",` +
					`"responders":[{"id":"4513b7ea-3b91-438f-b7e4-e3e54af9147c","type":"team"},{"name":"NOC","type":"team"}],` +
					`"visibleTo":[{"id":"4513b7ea-3b91-438f-b7e4-e3e54af9147c","type":"team"},{"name":"rocket_team","type":"team"}],` +
					`"actions":"An action",` +
					`"tags":"tag1 tag2",` +
					`"details":{"key1":"value1","key2":"value2"},` +
					`"entity":"An example entity",` +
					`"source":"The source",` +
					`"priority":"P1",` +
					`"user":"Dracula",` +
					`"note":"Here is a note"` +
					`}`))
			}

			err := service.Send("An example alert message", &types.Params{
				"alias":       "Life is too short for no alias",
				"description": "Every alert needs a description",
				"responders":  `[{"id":"4513b7ea-3b91-438f-b7e4-e3e54af9147c","type":"team"},{"name":"NOC","type":"team"}]`,
				"visibleTo":   `[{"id":"4513b7ea-3b91-438f-b7e4-e3e54af9147c","type":"team"},{"name":"rocket_team","type":"team"}]`,
				"actions":     "An action",
				"tags":        "tag1 tag2",
				"details":     `{"key1": "value1", "key2": "value2"}`,
				"entity":      "An example entity",
				"source":      "The source",
				"priority":    "P1",
				"user":        "Dracula",
				"note":        "Here is a note",
			})
			Expect(err).To(BeNil())
		})
	})
})
