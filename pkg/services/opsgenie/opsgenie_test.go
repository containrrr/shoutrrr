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
	mockHost   = "api.opsgenie.com"
)

func TestOpsGenie(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr OpsGenie Suite")
}

var _ = Describe("the OpsGenie service", func() {
	var (
		// a simulated http server to mock out OpsGenie itself
		mockServer *httptest.Server
		// the host of our mock server
		mockHost string
		// function to check if the http request received by the mock server is as expected
		checkRequest func(body string, header http.Header)
		// the shoutrrr OpsGenie service
		service *Service
		// just a mock logger
		mockLogger *log.Logger
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
			serviceURL, err := url.Parse(fmt.Sprintf("opsgenie://%s/%s", mockHost, mockAPIKey))
			Expect(err).To(BeNil())

			service = &Service{}
			err = service.Initialize(serviceURL, mockLogger)
			Expect(err).To(BeNil())
		})

		When("sending a simple alert", func() {
			It("should send a request to our mock OpsGenie server", func() {
				checkRequest = func(body string, header http.Header) {
					Expect(header["Authorization"][0]).To(Equal("GenieKey " + mockAPIKey))
					Expect(header["Content-Type"][0]).To(Equal("application/json"))
					Expect(body).To(Equal(`{"message":"hello world"}`))
				}

				err := service.Send("hello world", &types.Params{})
				Expect(err).To(BeNil())
			})
		})

		When("sending an alert with runtime parameters", func() {
			It("should send a request to our mock OpsGenie server with all fields populated from runtime parameters", func() {
				checkRequest = func(body string, header http.Header) {
					Expect(header["Authorization"][0]).To(Equal("GenieKey " + mockAPIKey))
					Expect(header["Content-Type"][0]).To(Equal("application/json"))
					Expect(body).To(Equal(`{"` +
						`message":"An example alert message",` +
						`"alias":"Life is too short for no alias",` +
						`"description":"Every alert needs a description",` +
						`"responders":[{"type":"team","id":"4513b7ea-3b91-438f-b7e4-e3e54af9147c"},{"type":"team","name":"NOC"},{"type":"user","username":"Donald"},{"type":"user","id":"696f0759-3b0f-4a15-b8c8-19d3dfca33f2"}],` +
						`"visibleTo":[{"type":"team","name":"rocket"}],` +
						`"actions":["action1","action2"],` +
						`"tags":["tag1","tag2"],` +
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
					"responders":  "team:4513b7ea-3b91-438f-b7e4-e3e54af9147c,team:NOC,user:Donald,user:696f0759-3b0f-4a15-b8c8-19d3dfca33f2",
					"visibleTo":   "team:rocket",
					"actions":     "action1,action2",
					"tags":        "tag1,tag2",
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

	Context("with query parameters", func() {
		BeforeEach(func() {
			// Initialize service
			serviceURL, err := url.Parse(fmt.Sprintf(`opsgenie://%s/%s?alias=query-alias&description=query-description&responders=team:query_team&visibleTo=user:query_user&actions=queryAction1,queryAction2&tags=queryTag1,queryTag2&details={"queryKey1": "queryValue1", "queryKey2": "queryValue2"}&entity=query-entity&source=query-source&priority=P2&user=query-user&note=query-note`, mockHost, mockAPIKey))
			Expect(err).To(BeNil())

			service = &Service{}
			err = service.Initialize(serviceURL, mockLogger)
			Expect(err).To(BeNil())
		})

		When("sending a simple alert", func() {
			It("should send a request to our mock OpsGenie server with all fields populated from query parameters", func() {
				checkRequest = func(body string, header http.Header) {
					Expect(header["Authorization"][0]).To(Equal("GenieKey " + mockAPIKey))
					Expect(header["Content-Type"][0]).To(Equal("application/json"))
					Expect(body).To(Equal(`{` +
						`"message":"An example alert message",` +
						`"alias":"query-alias",` +
						`"description":"query-description",` +
						`"responders":[{"type":"team","name":"query_team"}],` +
						`"visibleTo":[{"type":"user","username":"query_user"}],` +
						`"actions":["queryAction1","queryAction2"],` +
						`"tags":["queryTag1","queryTag2"],` +
						`"details":{"queryKey1":"queryValue1","queryKey2":"queryValue2"},` +
						`"entity":"query-entity",` +
						`"source":"query-source",` +
						`"priority":"P2",` +
						`"user":"query-user",` +
						`"note":"query-note"` +
						`}`))
				}

				err := service.Send("An example alert message", &types.Params{})
				Expect(err).To(BeNil())
			})
		})

		When("sending an alert with runtime parameters", func() {
			It("should send a request to our mock OpsGenie server with all fields populated from runtime parameters, overwriting the query parameters", func() {
				checkRequest = func(body string, header http.Header) {
					Expect(header["Authorization"][0]).To(Equal("GenieKey " + mockAPIKey))
					Expect(header["Content-Type"][0]).To(Equal("application/json"))
					Expect(body).To(Equal(`{"` +
						`message":"An example alert message",` +
						`"alias":"Life is too short for no alias",` +
						`"description":"Every alert needs a description",` +
						`"responders":[{"type":"team","id":"4513b7ea-3b91-438f-b7e4-e3e54af9147c"},{"type":"team","name":"NOC"},{"type":"user","username":"Donald"},{"type":"user","id":"696f0759-3b0f-4a15-b8c8-19d3dfca33f2"}],` +
						`"visibleTo":[{"type":"team","name":"rocket"}],` +
						`"actions":["action1","action2"],` +
						`"tags":["tag1","tag2"],` +
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
					"responders":  "team:4513b7ea-3b91-438f-b7e4-e3e54af9147c,team:NOC,user:Donald,user:696f0759-3b0f-4a15-b8c8-19d3dfca33f2",
					"visibleTo":   "team:rocket",
					"actions":     "action1,action2",
					"tags":        "tag1,tag2",
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
})

var _ = Describe("the OpsGenie Config struct", func() {
	When("generating a config from a simple URL", func() {
		It("should populate the config with host and apikey", func() {
			url, err := url.Parse(fmt.Sprintf("opsgenie://%s/%s", mockHost, mockAPIKey))
			Expect(err).To(BeNil())

			config := Config{}
			err = config.SetURL(url)
			Expect(err).To(BeNil())

			Expect(config.ApiKey).To(Equal(mockAPIKey))
			Expect(config.Host).To(Equal(mockHost))
			Expect(config.Port).To(Equal(uint16(0)))
		})
	})

	When("generating a config from a url with port", func() {
		It("should populate the port field", func() {
			url, err := url.Parse(fmt.Sprintf("opsgenie://%s:12345/%s", mockHost, mockAPIKey))
			Expect(err).To(BeNil())

			config := Config{}
			err = config.SetURL(url)
			Expect(err).To(BeNil())

			Expect(config.Port).To(Equal(uint16(12345)))
		})
	})

	When("generating a config from a url with query parameters", func() {
		It("should populate the relevant fields with the query parameter values", func() {
			queryParams := `alias=Life+is+too+short+for+no+alias&description=Every+alert+needs+a+description&actions=An+action&tags=tag1,tag2&details=these+are+details&entity=An+example+entity&source=The+source&priority=P1&user=Dracula&note=Here+is+a+note&responders=user:Test,team:NOC&visibleTo=user:A+User`
			url, err := url.Parse(fmt.Sprintf("opsgenie://%s:12345/%s?%s", mockHost, mockAPIKey, queryParams))
			Expect(err).To(BeNil())

			config := Config{}
			err = config.SetURL(url)
			Expect(err).To(BeNil())

			Expect(config.Alias).To(Equal("Life is too short for no alias"))
			Expect(config.Description).To(Equal("Every alert needs a description"))
			Expect(config.Responders).To(Equal([]Entity{
				{Type: "user", Username: "Test"},
				{Type: "team", Name: "NOC"},
			}))
			Expect(config.VisibleTo).To(Equal([]Entity{
				{Type: "user", Username: "A User"},
			}))
			Expect(config.Actions).To(Equal([]string{"An action"}))
			Expect(config.Tags).To(Equal([]string{"tag1", "tag2"}))
			Expect(config.Details).To(Equal("these are details"))
			Expect(config.Entity).To(Equal("An example entity"))
			Expect(config.Source).To(Equal("The source"))
			Expect(config.Priority).To(Equal("P1"))
			Expect(config.User).To(Equal("Dracula"))
			Expect(config.Note).To(Equal("Here is a note"))

		})
	})

	When("generating a url from a simple config", func() {
		It("should generate a url", func() {
			config := Config{
				Host:   "api.opsgenie.com",
				ApiKey: "eb243592-faa2-4ba2-a551q-1afdf565c889",
			}

			url := config.GetURL()

			Expect(url.String()).To(Equal("opsgenie://api.opsgenie.com/eb243592-faa2-4ba2-a551q-1afdf565c889"))
		})
	})

	When("generating a url from a config with a port", func() {
		It("should generate a url with port", func() {
			config := Config{
				Host:   "api.opsgenie.com",
				ApiKey: "eb243592-faa2-4ba2-a551q-1afdf565c889",
				Port:   12345,
			}

			url := config.GetURL()

			Expect(url.String()).To(Equal("opsgenie://api.opsgenie.com:12345/eb243592-faa2-4ba2-a551q-1afdf565c889"))
		})
	})

	When("generating a url from a config with all optional config fields", func() {
		It("should generate a url with query parameters", func() {
			config := Config{
				Host:        "api.opsgenie.com",
				ApiKey:      "eb243592-faa2-4ba2-a551q-1afdf565c889",
				Alias:       "Life is too short for no alias",
				Description: "Every alert needs a description",
				Responders: []Entity{
					{Type: "user", Username: "Test"},
					{Type: "team", Name: "NOC"},
				},
				VisibleTo: []Entity{
					{Type: "user", Username: "A User"},
				},
				Actions:  []string{"action1", "action2"},
				Tags:     []string{"tag1", "tag2"},
				Details:  "these are details",
				Entity:   "An example entity",
				Source:   "The source",
				Priority: "P1",
				User:     "Dracula",
				Note:     "Here is a note",
			}

			url := config.GetURL()
			fmt.Println(url.String())
			//&responders=user:Test,team:NOC&visibleTo=user:A+User
			Expect(url.String()).To(Equal(`opsgenie://api.opsgenie.com/eb243592-faa2-4ba2-a551q-1afdf565c889?actions=action1,action2&alias=Life is too short for no alias&description=Every alert needs a description&details=these are details&entity=An example entity&note=Here is a note&priority=P1&source=The source&tags=tag1,tag2&user=Dracula`))
		})
	})
})
