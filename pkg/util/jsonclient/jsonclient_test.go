package jsonclient_test

import (
	"errors"
	"net"
	"net/http"
	"testing"

	"github.com/nicholas-fedor/shoutrrr/pkg/util/jsonclient"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

func TestJSONClient(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "JSONClient Suite")
}

var _ = ginkgo.Describe("JSONClient", func() {
	var server *ghttp.Server
	var client jsonclient.Client

	ginkgo.BeforeEach(func() {
		server = ghttp.NewServer()
		client = jsonclient.NewClient()
	})

	ginkgo.When("the server returns an invalid JSON response", func() {
		ginkgo.It("should return an error", func() {
			server.AppendHandlers(ghttp.RespondWith(http.StatusOK, "invalid json"))
			res := &mockResponse{}
			err := client.Get(server.URL(), res)
			gomega.Expect(server.ReceivedRequests()).Should(gomega.HaveLen(1))
			gomega.Expect(err).To(gomega.MatchError("invalid character 'i' looking for beginning of value"))
			gomega.Expect(res.Status).To(gomega.BeEmpty())
		})
	})

	ginkgo.When("the server returns an empty response", func() {
		ginkgo.It("should return an error", func() {
			server.AppendHandlers(ghttp.RespondWith(http.StatusOK, nil))
			res := &mockResponse{}
			err := client.Get(server.URL(), res)
			gomega.Expect(server.ReceivedRequests()).Should(gomega.HaveLen(1))
			gomega.Expect(err).To(gomega.MatchError("unexpected end of JSON input"))
			gomega.Expect(res.Status).To(gomega.BeEmpty())
		})
	})

	ginkgo.It("should deserialize GET response", func() {
		server.AppendHandlers(ghttp.RespondWithJSONEncoded(http.StatusOK, mockResponse{Status: "OK"}))
		res := &mockResponse{}
		err := client.Get(server.URL(), res)
		gomega.Expect(server.ReceivedRequests()).Should(gomega.HaveLen(1))
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		gomega.Expect(res.Status).To(gomega.Equal("OK"))
	})

	ginkgo.Describe("Top-level Functions", func() {
		ginkgo.It("should handle GET via DefaultClient", func() {
			server.AppendHandlers(ghttp.RespondWithJSONEncoded(http.StatusOK, mockResponse{Status: "Default OK"}))
			res := &mockResponse{}
			err := jsonclient.Get(server.URL(), res)
			gomega.Expect(server.ReceivedRequests()).Should(gomega.HaveLen(1))
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(res.Status).To(gomega.Equal("Default OK"))
		})

		ginkgo.It("should handle POST via DefaultClient", func() {
			req := &mockRequest{Number: 10}
			res := &mockResponse{}
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/"),
				ghttp.VerifyJSONRepresenting(&req),
				ghttp.RespondWithJSONEncoded(http.StatusOK, mockResponse{Status: "Default POST"})),
			)
			err := jsonclient.Post(server.URL(), req, res)
			gomega.Expect(server.ReceivedRequests()).Should(gomega.HaveLen(1))
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(res.Status).To(gomega.Equal("Default POST"))
		})
	})

	ginkgo.Describe("POST", func() {
		ginkgo.It("should de-/serialize request and response", func() {
			req := &mockRequest{Number: 5}
			res := &mockResponse{}
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/"),
				ghttp.VerifyJSONRepresenting(&req),
				ghttp.RespondWithJSONEncoded(http.StatusOK, &mockResponse{Status: "That's Numberwang!"})),
			)
			err := client.Post(server.URL(), req, res)
			gomega.Expect(server.ReceivedRequests()).Should(gomega.HaveLen(1))
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(res.Status).To(gomega.Equal("That's Numberwang!"))
		})

		ginkgo.It("should return error on error status responses", func() {
			server.AppendHandlers(ghttp.RespondWith(http.StatusNotFound, "Not found!"))
			err := client.Post(server.URL(), &mockRequest{}, &mockResponse{})
			gomega.Expect(server.ReceivedRequests()).Should(gomega.HaveLen(1))
			gomega.Expect(err).To(gomega.MatchError("got HTTP 404 Not Found"))
		})

		ginkgo.It("should return error on invalid request", func() {
			server.AppendHandlers(ghttp.VerifyRequest("POST", "/"))
			err := client.Post(server.URL(), func() {}, &mockResponse{})
			gomega.Expect(server.ReceivedRequests()).Should(gomega.HaveLen(0))
			gomega.Expect(err).To(gomega.MatchError("error creating payload: json: unsupported type: func()"))
		})

		ginkgo.It("should return error on invalid response type", func() {
			res := &mockResponse{Status: "cool skirt"}
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/"),
				ghttp.RespondWithJSONEncoded(http.StatusOK, res)),
			)
			err := client.Post(server.URL(), nil, &[]bool{})
			gomega.Expect(server.ReceivedRequests()).Should(gomega.HaveLen(1))
			gomega.Expect(err).To(gomega.MatchError("json: cannot unmarshal object into Go value of type []bool"))
			gomega.Expect(jsonclient.ErrorBody(err)).To(gomega.MatchJSON(`{"Status":"cool skirt"}`))
		})

		ginkgo.It("should handle string request without marshaling", func() {
			rawJSON := `{"Number": 42}`
			res := &mockResponse{}
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/"),
				ghttp.VerifyBody([]byte(rawJSON)),
				ghttp.RespondWithJSONEncoded(http.StatusOK, mockResponse{Status: "String Worked"})),
			)
			err := client.Post(server.URL(), rawJSON, res)
			gomega.Expect(server.ReceivedRequests()).Should(gomega.HaveLen(1))
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(res.Status).To(gomega.Equal("String Worked"))
		})

		ginkgo.It("should return error when NewRequest fails", func() {
			err := client.Post("://invalid-url", &mockRequest{}, &mockResponse{})
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("error creating request"))
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("parse"))
		})

		ginkgo.It("should return error when http.Client.Do fails", func() {
			brokenClient := jsonclient.NewWithHTTPClient(&http.Client{
				Transport: &http.Transport{
					Dial: func(network, addr string) (net.Conn, error) {
						return nil, errors.New("forced network error")
					},
				},
			})
			err := brokenClient.Post(server.URL(), &mockRequest{}, &mockResponse{})
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("error sending payload"))
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("forced network error"))
		})

		ginkgo.It("should set multiple custom headers in request", func() {
			customClient := jsonclient.NewWithHTTPClient(&http.Client{})
			headers := customClient.Headers()
			headers.Set("X-Custom-Header", "CustomValue")
			headers.Set("X-Another-Header", "AnotherValue")

			req := &mockRequest{Number: 99}
			res := &mockResponse{}
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/"),
				ghttp.VerifyHeader(http.Header{
					"Content-Type":     []string{jsonclient.ContentType},
					"X-Custom-Header":  []string{"CustomValue"},
					"X-Another-Header": []string{"AnotherValue"},
				}),
				ghttp.VerifyJSONRepresenting(&req),
				ghttp.RespondWithJSONEncoded(http.StatusOK, mockResponse{Status: "Headers Worked"})),
			)
			err := customClient.Post(server.URL(), req, res)
			gomega.Expect(server.ReceivedRequests()).Should(gomega.HaveLen(1))
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(res.Status).To(gomega.Equal("Headers Worked"))
		})
	})

	ginkgo.Describe("Headers", func() {
		ginkgo.It("should return default headers with Content-Type", func() {
			headers := client.Headers()
			gomega.Expect(headers.Get("Content-Type")).To(gomega.Equal(jsonclient.ContentType))
		})
	})

	ginkgo.Describe("ErrorResponse", func() {
		ginkgo.It("should return false for non-jsonclient.Error", func() {
			res := &mockResponse{}
			result := client.ErrorResponse(errors.New("generic error"), res)
			gomega.Expect(result).To(gomega.BeFalse())
			gomega.Expect(res.Status).To(gomega.BeEmpty())
		})

		ginkgo.It("should populate response from jsonclient.Error body", func() {
			res := &mockResponse{}
			jsonErr := jsonclient.Error{
				StatusCode: http.StatusBadRequest,
				Body:       `{"Status": "Bad Request"}`,
			}
			result := client.ErrorResponse(jsonErr, res)
			gomega.Expect(result).To(gomega.BeTrue())
			gomega.Expect(res.Status).To(gomega.Equal("Bad Request"))
		})

		ginkgo.It("should return false for invalid JSON in error body", func() {
			res := &mockResponse{}
			jsonErr := jsonclient.Error{
				StatusCode: http.StatusBadRequest,
				Body:       "not json",
			}
			result := client.ErrorResponse(jsonErr, res)
			gomega.Expect(result).To(gomega.BeFalse())
			gomega.Expect(res.Status).To(gomega.BeEmpty())
		})
	})

	ginkgo.Describe("Edge Cases", func() {
		ginkgo.It("should handle network failure in Get", func() {
			res := &mockResponse{}
			err := client.Get("http://127.0.0.1:54321", res)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("dial tcp"))
			gomega.Expect(res.Status).To(gomega.BeEmpty())
		})

		ginkgo.It("should handle invalid JSON with success status", func() {
			server.AppendHandlers(ghttp.RespondWith(http.StatusOK, "bad json"))
			res := &mockResponse{}
			err := client.Get(server.URL(), res)
			gomega.Expect(server.ReceivedRequests()).Should(gomega.HaveLen(1))
			gomega.Expect(err).To(gomega.MatchError("invalid character 'b' looking for beginning of value"))
			gomega.Expect(res.Status).To(gomega.BeEmpty())
		})

		ginkgo.It("should handle nil body in error response", func() {
			brokenClient := jsonclient.NewWithHTTPClient(&http.Client{
				Transport: &mockTransport{
					response: &http.Response{
						StatusCode: http.StatusBadRequest,
						Status:     "400 Bad Request",
						Body:       &failingReader{},
						Header:     make(http.Header),
					},
				},
			})
			res := &mockResponse{}
			err := brokenClient.Get(server.URL(), res)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("got HTTP 400"))
			gomega.Expect(jsonclient.ErrorBody(err)).To(gomega.Equal(""))
		})
	})

	ginkgo.AfterEach(func() {
		server.Close()
	})
})

var _ = ginkgo.Describe("Error", func() {
	ginkgo.When("no internal error has been set", func() {
		ginkgo.It("should return a generic message with status code", func() {
			errorWithNoError := jsonclient.Error{StatusCode: http.StatusEarlyHints}
			gomega.Expect(errorWithNoError.String()).To(gomega.Equal("unknown error (HTTP 103)"))
		})
	})

	ginkgo.Describe("ErrorBody", func() {
		ginkgo.When("passed a non-json error", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(jsonclient.ErrorBody(errors.New("unrelated error"))).To(gomega.BeEmpty())
			})
		})

		ginkgo.When("passed a jsonclient.Error", func() {
			ginkgo.It("should return the request body from that error", func() {
				errorBody := `{"error": "bad user"}`
				jsonError := jsonclient.Error{Body: errorBody}
				gomega.Expect(jsonclient.ErrorBody(jsonError)).To(gomega.MatchJSON(errorBody))
			})
		})
	})
})

type mockResponse struct {
	Status string
}

type mockRequest struct {
	Number int
}

// brokenConn simulates a connection that fails on ReadAll.
type brokenConn struct {
	net.Conn
}

func (bc *brokenConn) Read(b []byte) (n int, err error) {
	return 0, errors.New("simulated read failure")
}

// mockTransport returns a predefined response.
type mockTransport struct {
	response *http.Response
}

func (mt *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return mt.response, nil
}

// failingReader simulates an io.Reader that fails on Read.
type failingReader struct{}

func (fr *failingReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("simulated read failure")
}

func (fr *failingReader) Close() error {
	return nil
}
