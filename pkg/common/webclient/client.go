package webclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// DefaultJsonClient is the singleton instance of WebClient using http.DefaultClient
var DefaultJsonClient = NewJSONClient()

// GetJson fetches url using GET and unmarshals into the passed response using DefaultJsonClient
func GetJson(url string, response interface{}) error {
	return DefaultJsonClient.Get(url, response)
}

// PostJson sends request as JSON and unmarshals the response JSON into the supplied struct using DefaultJsonClient
func PostJson(url string, request interface{}, response interface{}) error {
	return DefaultJsonClient.Post(url, request, response)
}

// WebClient is a JSON wrapper around http.WebClient
type client struct {
	headers    http.Header
	indent     string
	HttpClient http.Client
	parse      ParserFunc
	write      WriterFunc
}

// SetTransport overrides the http.RoundTripper for the web client, mainly used for testing
func (c *client) SetTransport(transport http.RoundTripper) {
	c.HttpClient.Transport = transport
}

// SetParser overrides the parser for the incoming response content
func (c *client) SetParser(parse ParserFunc) {
	c.parse = parse
}

// SetWriter overrides the writer for the outgoing request content
func (c *client) SetWriter(write WriterFunc) {
	c.write = write
}

// NewJSONClient returns a WebClient using the default http.Client and JSON serialization
func NewJSONClient() WebClient {
	var c client
	c = client{
		headers: http.Header{
			"Content-Type": []string{JsonContentType},
		},
		parse: json.Unmarshal,
		write: func(v interface{}) ([]byte, error) {
			return json.MarshalIndent(v, "", c.indent)
		},
	}
	return &c
}

// Headers return the default headers for requests
func (c *client) Headers() http.Header {
	return c.headers
}

// Get fetches url using GET and unmarshals into the passed response
func (c *client) Get(url string, response interface{}) error {
	res, err := c.HttpClient.Get(url)
	if err != nil {
		return err
	}

	return c.parseResponse(res, response)
}

// Post sends a serialized representation of request and deserializes the result into response
func (c *client) Post(url string, request interface{}, response interface{}) error {
	var err error
	var body []byte

	body, err = c.write(request)
	if err != nil {
		return fmt.Errorf("error creating payload: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	for key, val := range c.headers {
		req.Header.Set(key, val[0])
	}

	var res *http.Response
	res, err = c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending payload: %v", err)
	}

	return c.parseResponse(res, response)
}

// ErrorResponse tries to deserialize any response body into the supplied struct, returning whether successful or not
func (c *client) ErrorResponse(err error, response interface{}) bool {
	jerr, isWebError := err.(ClientError)
	if !isWebError {
		return false
	}

	return c.parse([]byte(jerr.Body), response) == nil
}

func (c *client) parseResponse(res *http.Response, response interface{}) error {
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if res.StatusCode >= 400 {
		err = fmt.Errorf("got HTTP %v", res.Status)
	}

	if err == nil {
		err = c.parse(body, response)
	}

	if err != nil {
		if body == nil {
			body = []byte{}
		}
		return ClientError{
			StatusCode: res.StatusCode,
			Body:       string(body),
			err:        err,
		}
	}

	return nil
}
