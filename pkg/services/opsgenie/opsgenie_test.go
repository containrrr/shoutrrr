package opsgenie

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/nicholas-fedor/shoutrrr/pkg/types"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

const (
	mockAPIKey = "eb243592-faa2-4ba2-a551q-1afdf565c889"
	mockHost   = "api.opsgenie.com"
)

func TestOpsGenie(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Shoutrrr OpsGenie Suite")
}

var _ = ginkgo.Describe("the OpsGenie service", func() {
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

	ginkgo.BeforeEach(func() {
		// Initialize a mock http server
		httpHandler := func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			gomega.Expect(err).To(gomega.BeNil())
			defer r.Body.Close()

			checkRequest(string(body), r.Header)
		}
		mockServer = httptest.NewTLSServer(http.HandlerFunc(httpHandler))

		// Our mock server doesn't have a valid cert
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		// Determine the host of our mock http server
		mockServerURL, err := url.Parse(mockServer.URL)
		gomega.Expect(err).To(gomega.BeNil())
		mockHost = mockServerURL.Host

		// Initialize a mock logger
		var buf bytes.Buffer
		mockLogger = log.New(&buf, "", 0)
	})

	ginkgo.AfterEach(func() {
		mockServer.Close()
	})

	ginkgo.Context("without query parameters", func() {
		ginkgo.BeforeEach(func() {
			// Initialize service
			serviceURL, err := url.Parse(fmt.Sprintf("opsgenie://%s/%s", mockHost, mockAPIKey))
			gomega.Expect(err).To(gomega.BeNil())

			service = &Service{}
			err = service.Initialize(serviceURL, mockLogger)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.When("sending a simple alert", func() {
			ginkgo.It("should send a request to our mock OpsGenie server", func() {
				checkRequest = func(body string, header http.Header) {
					gomega.Expect(header["Authorization"][0]).To(gomega.Equal("GenieKey " + mockAPIKey))
					gomega.Expect(header["Content-Type"][0]).To(gomega.Equal("application/json"))
					gomega.Expect(body).To(gomega.Equal(`{"message":"hello world"}`))
				}

				err := service.Send("hello world", &types.Params{})
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.When("sending an alert with runtime parameters", func() {
			ginkgo.It("should send a request to our mock OpsGenie server with all fields populated from runtime parameters", func() {
				checkRequest = func(body string, header http.Header) {
					gomega.Expect(header["Authorization"][0]).To(gomega.Equal("GenieKey " + mockAPIKey))
					gomega.Expect(header["Content-Type"][0]).To(gomega.Equal("application/json"))
					gomega.Expect(body).To(gomega.Equal(`{"` +
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
					"details":     "key1:value1,key2:value2",
					"entity":      "An example entity",
					"source":      "The source",
					"priority":    "P1",
					"user":        "Dracula",
					"note":        "Here is a note",
				})
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Context("with query parameters", func() {
		ginkgo.BeforeEach(func() {
			// Initialize service
			serviceURL, err := url.Parse(
				fmt.Sprintf(`opsgenie://%s/%s?alias=query-alias&description=query-description&responders=team:query_team&visibleTo=user:query_user&actions=queryAction1,queryAction2&tags=queryTag1,queryTag2&details=queryKey1:queryValue1,queryKey2:queryValue2&entity=query-entity&source=query-source&priority=P2&user=query-user&note=query-note`,
					mockHost,
					mockAPIKey,
				),
			)
			gomega.Expect(err).To(gomega.BeNil())

			service = &Service{}
			err = service.Initialize(serviceURL, mockLogger)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.When("sending a simple alert", func() {
			ginkgo.It("should send a request to our mock OpsGenie server with all fields populated from query parameters", func() {
				checkRequest = func(body string, header http.Header) {
					gomega.Expect(header["Authorization"][0]).To(gomega.Equal("GenieKey " + mockAPIKey))
					gomega.Expect(header["Content-Type"][0]).To(gomega.Equal("application/json"))
					gomega.Expect(body).To(gomega.Equal(`{` +
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
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.When("sending two alerts", func() {
			ginkgo.It("should not mix-up the runtime parameters and the query parameters", func() {
				// Internally the opsgenie service copies runtime parameters into the config struct
				// before generating the alert payload. This test ensures that none of the parameters
				// from alert 1 remain in the config struct when sending alert 2
				// In short: This tests if we clone the config struct

				checkRequest = func(body string, header http.Header) {
					gomega.Expect(header["Authorization"][0]).To(gomega.Equal("GenieKey " + mockAPIKey))
					gomega.Expect(header["Content-Type"][0]).To(gomega.Equal("application/json"))
					gomega.Expect(body).To(gomega.Equal(`{"` +
						`message":"1",` +
						`"alias":"1",` +
						`"description":"1",` +
						`"responders":[{"type":"team","name":"1"}],` +
						`"visibleTo":[{"type":"team","name":"1"}],` +
						`"actions":["action1","action2"],` +
						`"tags":["tag1","tag2"],` +
						`"details":{"key1":"value1","key2":"value2"},` +
						`"entity":"1",` +
						`"source":"1",` +
						`"priority":"P1",` +
						`"user":"1",` +
						`"note":"1"` +
						`}`))
				}

				err := service.Send("1", &types.Params{
					"alias":       "1",
					"description": "1",
					"responders":  "team:1",
					"visibleTo":   "team:1",
					"actions":     "action1,action2",
					"tags":        "tag1,tag2",
					"details":     "key1:value1,key2:value2",
					"entity":      "1",
					"source":      "1",
					"priority":    "P1",
					"user":        "1",
					"note":        "1",
				})
				gomega.Expect(err).To(gomega.BeNil())

				checkRequest = func(body string, header http.Header) {
					gomega.Expect(header["Authorization"][0]).To(gomega.Equal("GenieKey " + mockAPIKey))
					gomega.Expect(header["Content-Type"][0]).To(gomega.Equal("application/json"))
					gomega.Expect(body).To(gomega.Equal(`{` +
						`"message":"2",` +
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

				err = service.Send("2", nil)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.It("should return the correct service ID", func() {
		service := &Service{}
		gomega.Expect(service.GetID()).To(gomega.Equal("opsgenie"))
	})
})

var _ = ginkgo.Describe("the OpsGenie Config struct", func() {
	ginkgo.When("generating a config from a simple URL", func() {
		ginkgo.It("should populate the config with host and apikey", func() {
			url, err := url.Parse(fmt.Sprintf("opsgenie://%s/%s", mockHost, mockAPIKey))
			gomega.Expect(err).To(gomega.BeNil())

			config := Config{}
			err = config.SetURL(url)
			gomega.Expect(err).To(gomega.BeNil())

			gomega.Expect(config.APIKey).To(gomega.Equal(mockAPIKey))
			gomega.Expect(config.Host).To(gomega.Equal(mockHost))
			gomega.Expect(config.Port).To(gomega.Equal(uint16(443)))
		})
	})

	ginkgo.When("generating a config from a url with port", func() {
		ginkgo.It("should populate the port field", func() {
			url, err := url.Parse(fmt.Sprintf("opsgenie://%s/%s", net.JoinHostPort(mockHost, "12345"), mockAPIKey))
			gomega.Expect(err).To(gomega.BeNil())

			config := Config{}
			err = config.SetURL(url)
			gomega.Expect(err).To(gomega.BeNil())

			gomega.Expect(config.Port).To(gomega.Equal(uint16(12345)))
		})
	})

	ginkgo.When("generating a config from a url with query parameters", func() {
		ginkgo.It("should populate the config fields with the query parameter values", func() {
			queryParams := `alias=Life+is+too+short+for+no+alias&description=Every+alert+needs+a+description&actions=An+action&tags=tag1,tag2&details=key:value,key2:value2&entity=An+example+entity&source=The+source&priority=P1&user=Dracula&note=Here+is+a+note&responders=user:Test,team:NOC&visibleTo=user:A+User`
			url, err := url.Parse(fmt.Sprintf("opsgenie://%s/%s?%s", net.JoinHostPort(mockHost, "12345"), mockAPIKey, queryParams))
			gomega.Expect(err).To(gomega.BeNil())

			config := Config{}
			err = config.SetURL(url)
			gomega.Expect(err).To(gomega.BeNil())

			gomega.Expect(config.Alias).To(gomega.Equal("Life is too short for no alias"))
			gomega.Expect(config.Description).To(gomega.Equal("Every alert needs a description"))
			gomega.Expect(config.Responders).To(gomega.Equal([]Entity{
				{Type: "user", Username: "Test"},
				{Type: "team", Name: "NOC"},
			}))
			gomega.Expect(config.VisibleTo).To(gomega.Equal([]Entity{
				{Type: "user", Username: "A User"},
			}))
			gomega.Expect(config.Actions).To(gomega.Equal([]string{"An action"}))
			gomega.Expect(config.Tags).To(gomega.Equal([]string{"tag1", "tag2"}))
			gomega.Expect(config.Details).To(gomega.Equal(map[string]string{"key": "value", "key2": "value2"}))
			gomega.Expect(config.Entity).To(gomega.Equal("An example entity"))
			gomega.Expect(config.Source).To(gomega.Equal("The source"))
			gomega.Expect(config.Priority).To(gomega.Equal("P1"))
			gomega.Expect(config.User).To(gomega.Equal("Dracula"))
			gomega.Expect(config.Note).To(gomega.Equal("Here is a note"))
		})
	})

	ginkgo.When("generating a config from a url with differently escaped spaces", func() {
		ginkgo.It("should parse the escaped spaces correctly", func() {
			// Use: '%20', '+' and a normal space
			queryParams := `alias=Life is+too%20short+for+no+alias`
			url, err := url.Parse(fmt.Sprintf("opsgenie://%s/%s?%s", net.JoinHostPort(mockHost, "12345"), mockAPIKey, queryParams))
			gomega.Expect(err).To(gomega.BeNil())

			config := Config{}
			err = config.SetURL(url)
			gomega.Expect(err).To(gomega.BeNil())

			gomega.Expect(config.Alias).To(gomega.Equal("Life is too short for no alias"))
		})
	})

	ginkgo.When("generating a url from a simple config", func() {
		ginkgo.It("should generate a url", func() {
			config := Config{
				Host:   "api.opsgenie.com",
				APIKey: "eb243592-faa2-4ba2-a551q-1afdf565c889",
			}

			url := config.GetURL()

			gomega.Expect(url.String()).To(gomega.Equal("opsgenie://api.opsgenie.com/eb243592-faa2-4ba2-a551q-1afdf565c889"))
		})
	})

	ginkgo.When("generating a url from a config with a port", func() {
		ginkgo.It("should generate a url with port", func() {
			config := Config{
				Host:   "api.opsgenie.com",
				APIKey: "eb243592-faa2-4ba2-a551q-1afdf565c889",
				Port:   12345,
			}

			url := config.GetURL()

			gomega.Expect(url.String()).To(gomega.Equal("opsgenie://api.opsgenie.com:12345/eb243592-faa2-4ba2-a551q-1afdf565c889"))
		})
	})

	ginkgo.When("generating a url from a config with all optional config fields", func() {
		ginkgo.It("should generate a url with query parameters", func() {
			config := Config{
				Host:        "api.opsgenie.com",
				APIKey:      "eb243592-faa2-4ba2-a551q-1afdf565c889",
				Alias:       "Life is too short for no alias",
				Description: "Every alert needs a description",
				Responders: []Entity{
					{Type: "user", Username: "Test"},
					{Type: "team", Name: "NOC"},
					{Type: "team", ID: "4513b7ea-3b91-438f-b7e4-e3e54af9147c"},
				},
				VisibleTo: []Entity{
					{Type: "user", Username: "A User"},
				},
				Actions:  []string{"action1", "action2"},
				Tags:     []string{"tag1", "tag2"},
				Details:  map[string]string{"key": "value"},
				Entity:   "An example entity",
				Source:   "The source",
				Priority: "P1",
				User:     "Dracula",
				Note:     "Here is a note",
			}

			url := config.GetURL()
			gomega.Expect(url.String()).To(gomega.Equal(`opsgenie://api.opsgenie.com/eb243592-faa2-4ba2-a551q-1afdf565c889?actions=action1%2Caction2&alias=Life+is+too+short+for+no+alias&description=Every+alert+needs+a+description&details=key%3Avalue&entity=An+example+entity&note=Here+is+a+note&priority=P1&responders=user%3ATest%2Cteam%3ANOC%2Cteam%3A4513b7ea-3b91-438f-b7e4-e3e54af9147c&source=The+source&tags=tag1%2Ctag2&user=Dracula&visibleto=user%3AA+User`))
		})
	})
})
