package webclient_test

import (
	"errors"
	"github.com/containrrr/shoutrrr/pkg/common/webclient"
	"github.com/onsi/gomega/ghttp"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WebClient", func() {
	var server *ghttp.Server

	BeforeEach(func() {
		server = ghttp.NewServer()
	})

	When("the server returns an invalid JSON response", func() {
		It("should return an error", func() {
			server.AppendHandlers(ghttp.RespondWith(http.StatusOK, "invalid json"))
			res := &mockResponse{}
			err := webclient.GetJSON(server.URL(), &res)
			Expect(server.ReceivedRequests()).Should(HaveLen(1))
			Expect(err).To(MatchError("invalid character 'i' looking for beginning of value"))
			Expect(res.Status).To(BeEmpty())
		})
	})

	When("the server returns an empty response", func() {
		It("should return an error", func() {
			server.AppendHandlers(ghttp.RespondWith(http.StatusOK, nil))
			res := &mockResponse{}
			err := webclient.GetJSON(server.URL(), &res)
			Expect(server.ReceivedRequests()).Should(HaveLen(1))
			Expect(err).To(MatchError("unexpected end of JSON input"))
			Expect(res.Status).To(BeEmpty())
		})
	})

	It("should deserialize GET response", func() {
		server.AppendHandlers(ghttp.RespondWithJSONEncoded(http.StatusOK, mockResponse{Status: "OK"}))
		res := &mockResponse{}
		err := webclient.GetJSON(server.URL(), &res)
		Expect(server.ReceivedRequests()).Should(HaveLen(1))
		Expect(err).ToNot(HaveOccurred())
		Expect(res.Status).To(Equal("OK"))
	})

	It("should update the parser and writer", func() {
		client := webclient.NewJSONClient()
		client.SetParser(func(raw []byte, v interface{}) error {
			return errors.New(`mock parser`)
		})
		server.AppendHandlers(ghttp.RespondWithJSONEncoded(http.StatusOK, mockResponse{Status: "OK"}))
		err := client.Get(server.URL(), nil)
		Expect(err).To(MatchError(`mock parser`))

		client.SetWriter(func(v interface{}) ([]byte, error) {
			return nil, errors.New(`mock writer`)
		})
		err = client.Post(server.URL(), nil, nil)
		Expect(err).To(MatchError(`error creating payload: mock writer`))
	})

	It("should unwrap serialized error responses", func() {
		client := webclient.NewJSONClient()
		err := webclient.ClientError{Body: `{"Status": "BadStuff"}`}
		res := &mockResponse{}
		Expect(client.ErrorResponse(err, res)).To(BeTrue())
		Expect(res.Status).To(Equal(`BadStuff`))
	})

	It("should send any additional headers that has been added", func() {
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyHeaderKV(`Authentication`, `you don't need to see my identification`),
				ghttp.RespondWithJSONEncoded(http.StatusOK, mockResponse{Status: "OK"}),
			),
		)
		client := webclient.NewJSONClient()
		client.Headers().Set(`Authentication`, `you don't need to see my identification`)
		res := &mockResponse{}
		err := client.Get(server.URL(), &res)
		Expect(server.ReceivedRequests()).Should(HaveLen(1))
		Expect(err).ToNot(HaveOccurred())
		Expect(res.Status).To(Equal("OK"))
	})

	Describe("POST", func() {
		It("should de-/serialize request and response", func() {

			req := &mockRequest{Number: 5}
			res := &mockResponse{}

			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/"),
				ghttp.VerifyJSONRepresenting(&req),
				ghttp.RespondWithJSONEncoded(http.StatusOK, &mockResponse{Status: "That's Numberwang!"})),
			)

			err := webclient.PostJSON(server.URL(), &req, &res)
			Expect(server.ReceivedRequests()).Should(HaveLen(1))
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Status).To(Equal("That's Numberwang!"))
		})

		It("should return error on error status responses", func() {
			server.AppendHandlers(ghttp.RespondWith(404, "Not found!"))
			err := webclient.PostJSON(server.URL(), &mockRequest{}, &mockResponse{})
			Expect(server.ReceivedRequests()).Should(HaveLen(1))
			Expect(err).To(MatchError("got HTTP 404 Not Found"))
		})

		It("should return error on invalid request", func() {
			server.AppendHandlers(ghttp.VerifyRequest("POST", "/"))
			err := webclient.PostJSON(server.URL(), func() {}, &mockResponse{})
			Expect(server.ReceivedRequests()).Should(HaveLen(0))
			Expect(err).To(MatchError("error creating payload: json: unsupported type: func()"))
		})

		It("should return error on invalid response type", func() {
			res := &mockResponse{Status: "cool skirt"}
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/"),
				ghttp.RespondWithJSONEncoded(http.StatusOK, res)),
			)

			err := webclient.PostJSON(server.URL(), nil, &[]bool{})
			Expect(server.ReceivedRequests()).Should(HaveLen(1))
			Expect(err).To(MatchError("json: cannot unmarshal object into Go value of type []bool"))
			Expect(webclient.ErrorBody(err)).To(MatchJSON(`{"Status":"cool skirt"}`))
		})
	})

	AfterEach(func() {
		//shut down the server between tests
		server.Close()
	})
})

var _ = Describe("ClientError", func() {
	When("no internal error has been set", func() {
		It("should return a generic message with status code", func() {
			errorWithNoError := webclient.ClientError{StatusCode: http.StatusEarlyHints}
			Expect(errorWithNoError.String()).To(Equal("unknown error (HTTP 103)"))
		})
	})
	Describe("ErrorBody", func() {
		When("passed a non-json error", func() {
			It("should return an empty string", func() {
				Expect(webclient.ErrorBody(errors.New("unrelated error"))).To(BeEmpty())
			})
		})
		When("passed a jsonclient.ClientError", func() {
			It("should return the request body from that error", func() {
				errorBody := `{"error": "bad user"}`
				jsonError := webclient.ClientError{Body: errorBody}
				Expect(webclient.ErrorBody(jsonError)).To(MatchJSON(errorBody))
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
