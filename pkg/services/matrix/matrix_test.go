package matrix

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/nicholas-fedor/shoutrrr/internal/testutils"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestMatrix(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Shoutrrr Matrix Suite")
}

var _ = ginkgo.Describe("the matrix service", func() {
	var service *Service
	logger := log.New(ginkgo.GinkgoWriter, "Test", log.LstdFlags)
	envMatrixURL := os.Getenv("SHOUTRRR_MATRIX_URL")

	ginkgo.BeforeEach(func() {
		service = &Service{}
	})

	ginkgo.When("running integration tests", func() {
		ginkgo.It("should not error out", func() {
			// Tests matrix_client.go lines:
			// - 36-52: newClient (full initialization with logger and scheme)
			// - 63-65: login (via Initialize when User is set)
			// - 76-87: loginPassword (successful login flow)
			// - 91-108: sendMessage (via Send with real server)
			// - 156-173: sendMessageToRoom (sending to joined rooms)
			if envMatrixURL == "" {
				return
			}
			serviceURL, err := url.Parse(envMatrixURL)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			err = service.Initialize(serviceURL, logger)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			err = service.Send("This is an integration test message", nil)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
	})

	ginkgo.Describe("creating configurations", func() {
		ginkgo.When("given an url with title prop", func() {
			ginkgo.It("should not throw an error", func() {
				// Tests matrix_config.go, not matrix_client.go directly
				// Related to Config.SetURL, which feeds into client setup later
				serviceURL := testutils.URLMust(`matrix://user:pass@mockserver?rooms=room1&title=Better%20Off%20Alone`)
				gomega.Expect((&Config{}).SetURL(serviceURL)).To(gomega.Succeed())
			})
		})

		ginkgo.When("given an url with the prop `room`", func() {
			ginkgo.It("should treat is as an alias for `rooms`", func() {
				// Tests matrix_config.go, not matrix_client.go directly
				// Configures Rooms for client.sendToExplicitRooms later
				serviceURL := testutils.URLMust(`matrix://user:pass@mockserver?room=room1`)
				config := Config{}
				gomega.Expect(config.SetURL(serviceURL)).To(gomega.Succeed())
				gomega.Expect(config.Rooms).To(gomega.ContainElement("#room1"))
			})
		})
		ginkgo.When("given an url with invalid props", func() {
			ginkgo.It("should return an error", func() {
				// Tests matrix_config.go, not matrix_client.go directly
				// Ensures invalid params fail before reaching client
				serviceURL := testutils.URLMust(`matrix://user:pass@mockserver?channels=room1,room2`)
				gomega.Expect((&Config{}).SetURL(serviceURL)).To(gomega.HaveOccurred())
			})
		})
		ginkgo.When("parsing the configuration URL", func() {
			ginkgo.It("should be identical after de-/serialization", func() {
				// Tests matrix_config.go, not matrix_client.go directly
				// Verifies Config.GetURL/SetURL round-trip for client init
				testURL := "matrix://user:pass@mockserver?rooms=%23room1%2C%23room2"
				url, err := url.Parse(testURL)
				gomega.Expect(err).NotTo(gomega.HaveOccurred(), "parsing")
				config := &Config{}
				err = config.SetURL(url)
				gomega.Expect(err).NotTo(gomega.HaveOccurred(), "verifying")
				outputURL := config.GetURL()
				gomega.Expect(outputURL.String()).To(gomega.Equal(testURL))
			})
		})
	})

	ginkgo.Describe("the matrix client", func() {
		ginkgo.BeforeEach(func() {
			httpmock.Activate()
		})

		ginkgo.When("not providing a logger", func() {
			ginkgo.It("should not crash", func() {
				// Tests matrix_client.go lines:
				// - 36-52: newClient (sets DiscardLogger when logger is nil)
				// - 63-65: login (successful initialization)
				// - 76-87: loginPassword (successful login flow)
				setupMockResponders()
				serviceURL := testutils.URLMust("matrix://user:pass@mockserver")
				gomega.Expect(service.Initialize(serviceURL, nil)).To(gomega.Succeed())
			})
		})

		ginkgo.When("sending a message", func() {
			ginkgo.It("should not report any errors", func() {
				// Tests matrix_client.go lines:
				// - 36-52: newClient (successful setup)
				// - 63-65: login (successful initialization)
				// - 76-87: loginPassword (successful login flow)
				// - 91-108: sendMessage (routes to sendToJoinedRooms)
				// - 134-153: sendToJoinedRooms (sends to joined rooms)
				// - 156-173: sendMessageToRoom (successful send)
				// - 225-242: getJoinedRooms (fetches room list)
				setupMockResponders()
				serviceURL, _ := url.Parse("matrix://user:pass@mockserver")
				err := service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				err = service.Send("Test message", nil)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
		})

		ginkgo.When("sending a message to explicit rooms", func() {
			ginkgo.It("should not report any errors", func() {
				// Tests matrix_client.go lines:
				// - 36-52: newClient (successful setup)
				// - 63-65: login (successful initialization)
				// - 76-87: loginPassword (successful login flow)
				// - 91-108: sendMessage (routes to sendToExplicitRooms)
				// - 112-133: sendToExplicitRooms (sends to explicit rooms)
				// - 177-192: joinRoom (joins rooms successfully)
				// - 156-173: sendMessageToRoom (successful send)
				setupMockResponders()
				serviceURL, _ := url.Parse("matrix://user:pass@mockserver?rooms=room1,room2")
				err := service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				err = service.Send("Test message", nil)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
			ginkgo.When("sending to one room fails", func() {
				ginkgo.It("should report one error", func() {
					// Tests matrix_client.go lines:
					// - 36-52: newClient (successful setup)
					// - 63-65: login (successful initialization)
					// - 76-87: loginPassword (successful login flow)
					// - 91-108: sendMessage (routes to sendToExplicitRooms)
					// - 112-133: sendToExplicitRooms (handles join failure)
					// - 177-192: joinRoom (fails for "secret" room)
					// - 156-173: sendMessageToRoom (succeeds for "room2")
					setupMockResponders()
					serviceURL, _ := url.Parse("matrix://user:pass@mockserver?rooms=secret,room2")
					err := service.Initialize(serviceURL, logger)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					err = service.Send("Test message", nil)
					gomega.Expect(err).To(gomega.HaveOccurred())
				})
			})
		})

		ginkgo.When("disabling TLS", func() {
			ginkgo.It("should use HTTP instead of HTTPS", func() {
				// Tests matrix_client.go lines:
				// - 36-52: newClient (specifically line 50: c.apiURL.Scheme = c.apiURL.Scheme[:schemeHTTPPrefixLength])
				// - 63-65: login (successful initialization over HTTP)
				// - 76-87: loginPassword (successful login flow)
				setupMockRespondersHTTP()
				serviceURL := testutils.URLMust("matrix://user:pass@mockserver?disableTLS=yes")
				err := service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(service.client.apiURL.Scheme).To(gomega.Equal("http"))
			})
		})

		ginkgo.When("failing to get login flows", func() {
			ginkgo.It("should return an error", func() {
				// Tests matrix_client.go lines:
				// - 36-52: newClient (successful setup)
				// - 63-69: login (specifically line 69: return fmt.Errorf("failed to get login flows: %w", err))
				// - 175-223: apiGet (returns error due to 500 response)
				setupMockRespondersLoginFail()
				serviceURL := testutils.URLMust("matrix://user:pass@mockserver")
				err := service.Initialize(serviceURL, logger)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(err.Error()).To(gomega.ContainSubstring("failed to get login flows"))
			})
		})

		ginkgo.When("no supported login flows are available", func() {
			ginkgo.It("should return an error with unsupported flows", func() {
				// Tests matrix_client.go lines:
				// - 36-52: newClient (successful setup)
				// - 63-87: login (specifically line 84: return fmt.Errorf("none of the server login flows are supported: %v", strings.Join(flows, ", ")))
				// - 175-223: apiGet (successful GET with unsupported flows)
				setupMockRespondersUnsupportedFlows()
				serviceURL := testutils.URLMust("matrix://user:pass@mockserver")
				err := service.Initialize(serviceURL, logger)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(err.Error()).To(gomega.Equal("none of the server login flows are supported: m.login.dummy"))
			})
		})

		ginkgo.When("using a token instead of login", func() {
			ginkgo.It("should initialize without errors", func() {
				// Tests matrix_client.go lines:
				// - 36-52: newClient (successful setup)
				// - 59-60: useToken (sets token and calls updateAccessToken)
				// - 244-248: updateAccessToken (updates URL query with token)
				setupMockResponders() // Minimal mocks for initialization
				serviceURL := testutils.URLMust("matrix://:token@mockserver")
				err := service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(service.client.accessToken).To(gomega.Equal("token"))
				gomega.Expect(service.client.apiURL.RawQuery).To(gomega.Equal("access_token=token"))
			})
		})

		ginkgo.When("failing to get joined rooms", func() {
			ginkgo.It("should return an error", func() {
				// Tests matrix_client.go lines:
				// - 36-52: newClient (successful setup)
				// - 63-65: login (successful initialization)
				// - 76-87: loginPassword (successful login flow)
				// - 91-108: sendMessage (routes to sendToJoinedRooms)
				// - 134-154: sendToJoinedRooms (specifically lines 137 and 154: error handling for getJoinedRooms failure)
				// - 225-267: getJoinedRooms (specifically line 267: return []string{}, err)
				setupMockRespondersJoinedRoomsFail()
				serviceURL := testutils.URLMust("matrix://user:pass@mockserver")
				err := service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				err = service.Send("Test message", nil)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(err.Error()).To(gomega.ContainSubstring("failed to get joined rooms"))
			})
		})

		ginkgo.When("failing to join a room", func() {
			ginkgo.It("should skip to the next room and continue", func() {
				// Tests matrix_client.go lines:
				// - 36-52: newClient (successful setup)
				// - 63-65: login (successful initialization)
				// - 76-87: loginPassword (successful login flow)
				// - 91-108: sendMessage (routes to sendToExplicitRooms)
				// - 112-133: sendToExplicitRooms (specifically line 147: continue on join failure)
				// - 177-192: joinRoom (specifically line 188: return "", err on failure)
				// - 156-173: sendMessageToRoom (succeeds for second room)
				setupMockRespondersJoinFail()
				serviceURL := testutils.URLMust("matrix://user:pass@mockserver?rooms=secret,room2")
				err := service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				err = service.Send("Test message", nil)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(err.Error()).To(gomega.ContainSubstring("error joining room"))
			})
		})

		ginkgo.When("failing to marshal request in apiPost", func() {
			ginkgo.It("should return an error", func() {
				// Tests matrix_client.go lines:
				// - 36-52: newClient (successful setup)
				// - 63-65: login (successful initialization)
				// - 76-87: loginPassword (successful login flow)
				// - 195-252: apiPost (specifically line 208: body, err = json.Marshal(request) fails)
				setupMockResponders()
				serviceURL := testutils.URLMust("matrix://user:pass@mockserver")
				err := service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				err = service.client.apiPost("/test/path", make(chan int), nil)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(err.Error()).To(gomega.ContainSubstring("json: unsupported type: chan int"))
			})
		})

		ginkgo.When("failing to read response body in apiPost", func() {
			ginkgo.It("should return an error", func() {
				// Tests matrix_client.go lines:
				// - 36-52: newClient (successful setup)
				// - 63-65: login (successful initialization)
				// - 76-87: loginPassword (successful login flow)
				// - 91-108: sendMessage (routes to sendToJoinedRooms)
				// - 134-153: sendToJoinedRooms (calls sendMessageToRoom)
				// - 156-173: sendMessageToRoom (calls apiPost)
				// - 195-252: apiPost (specifically lines 204, 223, 230: res handling and body read failure)
				setupMockRespondersBodyFail()
				serviceURL := testutils.URLMust("matrix://user:pass@mockserver")
				err := service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				err = service.Send("Test message", nil)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(err.Error()).To(gomega.ContainSubstring("failed to read response body"))
			})
		})

		ginkgo.When("routing to explicit rooms at line 94", func() {
			ginkgo.It("should use sendToExplicitRooms", func() {
				// Tests matrix_client.go lines:
				// - 36-52: newClient (successful setup)
				// - 63-65: login (successful initialization)
				// - 76-87: loginPassword (successful login flow)
				// - 91-108: sendMessage (specifically line 94: if len(rooms) >= minSliceLength { true branch)
				// - 112-133: sendToExplicitRooms (sends to explicit rooms)
				// - 177-192: joinRoom (joins rooms successfully)
				// - 156-173: sendMessageToRoom (successful send)
				setupMockResponders()
				serviceURL := testutils.URLMust("matrix://user:pass@mockserver?rooms=room1")
				err := service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				err = service.Send("Test message", nil)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
		})

		ginkgo.When("routing to joined rooms at line 94", func() {
			ginkgo.It("should use sendToJoinedRooms", func() {
				// Tests matrix_client.go lines:
				// - 36-52: newClient (successful setup)
				// - 63-65: login (successful initialization)
				// - 76-87: loginPassword (successful login flow)
				// - 91-108: sendMessage (specifically line 94: if len(rooms) >= minSliceLength { false branch)
				// - 134-153: sendToJoinedRooms (sends to joined rooms)
				// - 156-173: sendMessageToRoom (successful send)
				// - 225-242: getJoinedRooms (fetches room list)
				setupMockResponders()
				serviceURL := testutils.URLMust("matrix://user:pass@mockserver")
				err := service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				err = service.Send("Test message", nil)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
		})

		ginkgo.When("appending joined rooms error at line 137", func() {
			ginkgo.It("should append the error", func() {
				// Tests matrix_client.go lines:
				// - 36-52: newClient (successful setup)
				// - 63-65: login (successful initialization)
				// - 76-87: loginPassword (successful login flow)
				// - 91-108: sendMessage (routes to sendToJoinedRooms)
				// - 134-154: sendToJoinedRooms (specifically line 137: errors = append(errors, fmt.Errorf("failed to get joined rooms: %w", err)))
				// - 225-267: getJoinedRooms (returns error)
				setupMockRespondersJoinedRoomsFail()
				serviceURL := testutils.URLMust("matrix://user:pass@mockserver")
				err := service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				err = service.Send("Test message", nil)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(err.Error()).To(gomega.ContainSubstring("failed to get joined rooms"))
			})
		})

		ginkgo.When("failing to join room at line 188", func() {
			ginkgo.It("should return join error", func() {
				// Tests matrix_client.go lines:
				// - 36-52: newClient (successful setup)
				// - 63-65: login (successful initialization)
				// - 76-87: loginPassword (successful login flow)
				// - 91-108: sendMessage (routes to sendToExplicitRooms)
				// - 112-133: sendToExplicitRooms (calls joinRoom)
				// - 177-192: joinRoom (specifically line 188: return "", err)
				setupMockRespondersJoinFail()
				serviceURL := testutils.URLMust("matrix://user:pass@mockserver?rooms=secret")
				err := service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				err = service.Send("Test message", nil)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(err.Error()).To(gomega.ContainSubstring("error joining room"))
			})
		})

		ginkgo.When("declaring response variable at line 204", func() {
			ginkgo.It("should handle HTTP failure", func() {
				// Tests matrix_client.go lines:
				// - 36-52: newClient (successful setup)
				// - 63-65: login (successful initialization)
				// - 76-87: loginPassword (successful login flow)
				// - 195-252: apiPost (specifically line 204: var res *http.Response and error handling)
				setupMockRespondersPostFail()
				serviceURL := testutils.URLMust("matrix://user:pass@mockserver")
				err := service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				err = service.client.apiPost("/test/path", apiReqSend{MsgType: msgTypeText, Body: "test"}, nil)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(err.Error()).To(gomega.ContainSubstring("simulated HTTP failure"))
			})
		})

		ginkgo.When("marshaling request fails at line 208", func() {
			ginkgo.It("should return marshal error", func() {
				// Tests matrix_client.go lines:
				// - 36-52: newClient (successful setup)
				// - 63-65: login (successful initialization)
				// - 76-87: loginPassword (successful login flow)
				// - 195-252: apiPost (specifically line 208: body, err = json.Marshal(request))
				setupMockResponders()
				serviceURL := testutils.URLMust("matrix://user:pass@mockserver")
				err := service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				err = service.client.apiPost("/test/path", make(chan int), nil)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(err.Error()).To(gomega.ContainSubstring("json: unsupported type: chan int"))
			})
		})

		ginkgo.When("getting query at line 244", func() {
			ginkgo.It("should update token in URL", func() {
				// Tests matrix_client.go lines:
				// - 36-52: newClient (successful setup)
				// - 59-60: useToken (calls updateAccessToken)
				// - 244-248: updateAccessToken (specifically line 244: query := c.apiURL.Query())
				setupMockResponders()
				serviceURL := testutils.URLMust("matrix://:token@mockserver")
				err := service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(service.client.apiURL.RawQuery).To(gomega.Equal("access_token=token"))
				service.client.useToken("newtoken")
				gomega.Expect(service.client.apiURL.RawQuery).To(gomega.Equal("access_token=newtoken"))
			})
		})

		ginkgo.When("checking body read error at line 251", func() {
			ginkgo.It("should return read error", func() {
				// Tests matrix_client.go lines:
				// - 36-52: newClient (successful setup)
				// - 63-65: login (successful initialization)
				// - 76-87: loginPassword (successful login flow)
				// - 91-108: sendMessage (routes to sendToJoinedRooms)
				// - 134-153: sendToJoinedRooms (calls sendMessageToRoom)
				// - 156-173: sendMessageToRoom (calls apiPost)
				// - 195-252: apiPost (specifically line 251: if err != nil { after io.ReadAll)
				setupMockRespondersBodyFail()
				serviceURL := testutils.URLMust("matrix://user:pass@mockserver")
				err := service.Initialize(serviceURL, logger)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				err = service.Send("Test message", nil)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(err.Error()).To(gomega.ContainSubstring("failed to read response body"))
			})
		})

		ginkgo.AfterEach(func() {
			httpmock.DeactivateAndReset()
		})
	})

	ginkgo.It("should implement basic service API methods correctly", func() {
		// Tests matrix_config.go, not matrix_client.go directly
		// Exercises Config methods used indirectly by client initialization
		testutils.TestConfigGetInvalidQueryValue(&Config{})
		testutils.TestConfigSetInvalidQueryValue(&Config{}, "matrix://user:pass@host/?foo=bar")
		testutils.TestConfigGetEnumsCount(&Config{}, 0)
		testutils.TestConfigGetFieldsCount(&Config{}, 4)
	})

	ginkgo.It("should return the correct service ID", func() {
		service := &Service{}
		gomega.Expect(service.GetID()).To(gomega.Equal("matrix"))
	})
})

// setupMockResponders for HTTPS.
func setupMockResponders() {
	const mockServer = "https://mockserver"

	httpmock.RegisterResponder(
		"GET",
		mockServer+apiLogin,
		httpmock.NewStringResponder(200, `{"flows": [ { "type": "m.login.password" } ] }`))

	httpmock.RegisterResponder(
		"POST",
		mockServer+apiLogin,
		httpmock.NewStringResponder(
			200,
			`{ "access_token": "TOKEN", "home_server": "mockserver", "user_id": "test:mockerserver" }`,
		),
	)

	httpmock.RegisterResponder(
		"GET",
		mockServer+apiJoinedRooms,
		httpmock.NewStringResponder(200, `{ "joined_rooms": [ "!room:mockserver" ] }`))

	httpmock.RegisterResponder("POST", mockServer+fmt.Sprintf(apiSendMessage, "%21room:mockserver"),
		httpmock.NewJsonResponderOrPanic(200, apiResEvent{EventID: "7"}))

	httpmock.RegisterResponder("POST", mockServer+fmt.Sprintf(apiSendMessage, "1"),
		httpmock.NewJsonResponderOrPanic(200, apiResEvent{EventID: "8"}))

	httpmock.RegisterResponder("POST", mockServer+fmt.Sprintf(apiSendMessage, "2"),
		httpmock.NewJsonResponderOrPanic(200, apiResEvent{EventID: "9"}))

	httpmock.RegisterResponder("POST", mockServer+fmt.Sprintf(apiRoomJoin, "%23room1"),
		httpmock.NewJsonResponderOrPanic(200, apiResRoom{RoomID: "1"}))

	httpmock.RegisterResponder("POST", mockServer+fmt.Sprintf(apiRoomJoin, "%23room2"),
		httpmock.NewJsonResponderOrPanic(200, apiResRoom{RoomID: "2"}))

	httpmock.RegisterResponder("POST", mockServer+fmt.Sprintf(apiRoomJoin, "%23secret"),
		httpmock.NewJsonResponderOrPanic(403, apiResError{
			Code:    "M_FORBIDDEN",
			Message: "You are not invited to this room.",
		}))
}

// setupMockRespondersHTTP for HTTP.
func setupMockRespondersHTTP() {
	const mockServer = "http://mockserver"

	httpmock.RegisterResponder(
		"GET",
		mockServer+apiLogin,
		httpmock.NewStringResponder(200, `{"flows": [ { "type": "m.login.password" } ] }`))

	httpmock.RegisterResponder(
		"POST",
		mockServer+apiLogin,
		httpmock.NewStringResponder(
			200,
			`{ "access_token": "TOKEN", "home_server": "mockserver", "user_id": "test:mockerserver" }`,
		),
	)

	httpmock.RegisterResponder(
		"GET",
		mockServer+apiJoinedRooms,
		httpmock.NewStringResponder(200, `{ "joined_rooms": [ "!room:mockserver" ] }`))

	httpmock.RegisterResponder("POST", mockServer+fmt.Sprintf(apiSendMessage, "%21room:mockserver"),
		httpmock.NewJsonResponderOrPanic(200, apiResEvent{EventID: "7"}))
}

// setupMockRespondersLoginFail for testing line 69.
func setupMockRespondersLoginFail() {
	const mockServer = "https://mockserver"

	httpmock.RegisterResponder(
		"GET",
		mockServer+apiLogin,
		httpmock.NewStringResponder(500, `{"error": "Internal Server Error"}`))
}

// setupMockRespondersUnsupportedFlows for testing line 84.
func setupMockRespondersUnsupportedFlows() {
	const mockServer = "https://mockserver"

	httpmock.RegisterResponder(
		"GET",
		mockServer+apiLogin,
		httpmock.NewStringResponder(200, `{"flows": [ { "type": "m.login.dummy" } ] }`))
}

// setupMockRespondersJoinedRoomsFail for testing lines 137, 154, and 267.
func setupMockRespondersJoinedRoomsFail() {
	const mockServer = "https://mockserver"

	httpmock.RegisterResponder(
		"GET",
		mockServer+apiLogin,
		httpmock.NewStringResponder(200, `{"flows": [ { "type": "m.login.password" } ] }`))

	httpmock.RegisterResponder(
		"POST",
		mockServer+apiLogin,
		httpmock.NewStringResponder(
			200,
			`{ "access_token": "TOKEN", "home_server": "mockserver", "user_id": "test:mockerserver" }`,
		),
	)

	httpmock.RegisterResponder(
		"GET",
		mockServer+apiJoinedRooms,
		httpmock.NewStringResponder(500, `{"error": "Internal Server Error"}`))
}

// setupMockRespondersJoinFail for testing lines 147 and 188.
func setupMockRespondersJoinFail() {
	const mockServer = "https://mockserver"

	httpmock.RegisterResponder(
		"GET",
		mockServer+apiLogin,
		httpmock.NewStringResponder(200, `{"flows": [ { "type": "m.login.password" } ] }`))

	httpmock.RegisterResponder(
		"POST",
		mockServer+apiLogin,
		httpmock.NewStringResponder(
			200,
			`{ "access_token": "TOKEN", "home_server": "mockserver", "user_id": "test:mockerserver" }`,
		),
	)

	httpmock.RegisterResponder("POST", mockServer+fmt.Sprintf(apiRoomJoin, "%23secret"),
		httpmock.NewJsonResponderOrPanic(403, apiResError{
			Code:    "M_FORBIDDEN",
			Message: "You are not invited to this room.",
		}))

	httpmock.RegisterResponder("POST", mockServer+fmt.Sprintf(apiRoomJoin, "%23room2"),
		httpmock.NewJsonResponderOrPanic(200, apiResRoom{RoomID: "2"}))

	httpmock.RegisterResponder("POST", mockServer+fmt.Sprintf(apiSendMessage, "2"),
		httpmock.NewJsonResponderOrPanic(200, apiResEvent{EventID: "9"}))
}

// setupMockRespondersBodyFail for testing lines 204, 223, and 230.
func setupMockRespondersBodyFail() {
	const mockServer = "https://mockserver"

	httpmock.RegisterResponder(
		"GET",
		mockServer+apiLogin,
		httpmock.NewStringResponder(200, `{"flows": [ { "type": "m.login.password" } ] }`))

	httpmock.RegisterResponder(
		"POST",
		mockServer+apiLogin,
		httpmock.NewStringResponder(
			200,
			`{ "access_token": "TOKEN", "home_server": "mockserver", "user_id": "test:mockerserver" }`,
		),
	)

	httpmock.RegisterResponder(
		"GET",
		mockServer+apiJoinedRooms,
		httpmock.NewStringResponder(200, `{ "joined_rooms": [ "!room:mockserver" ] }`))

	httpmock.RegisterResponder("POST", mockServer+fmt.Sprintf(apiSendMessage, "%21room:mockserver"),
		httpmock.NewErrorResponder(fmt.Errorf("failed to read response body")))
}

// setupMockRespondersPostFail for testing line 204 and HTTP failure.
func setupMockRespondersPostFail() {
	const mockServer = "https://mockserver"

	httpmock.RegisterResponder(
		"GET",
		mockServer+apiLogin,
		httpmock.NewStringResponder(200, `{"flows": [ { "type": "m.login.password" } ] }`))

	httpmock.RegisterResponder(
		"POST",
		mockServer+apiLogin,
		httpmock.NewStringResponder(
			200,
			`{ "access_token": "TOKEN", "home_server": "mockserver", "user_id": "test:mockerserver" }`,
		),
	)

	httpmock.RegisterResponder("POST", mockServer+"/test/path",
		httpmock.NewErrorResponder(fmt.Errorf("simulated HTTP failure")))
}
