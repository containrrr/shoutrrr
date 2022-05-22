package jsonclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ContentType is the default mime type for JSON
const ContentType = "application/json"

// DefaultClient is the singleton instance of jsonclient using http.DefaultClient
var DefaultClient = NewClient()

// Get fetches url using GET and unmarshals into the passed response using DefaultClient
func Get(url string, response interface{}) error {
	return DefaultClient.Get(url, response)
}

// Post sends request as JSON and unmarshals the response JSON into the supplied struct using DefaultClient
func Post(url string, request interface{}, response interface{}) error {
	return DefaultClient.Post(url, request, response)
}

// Client is a JSON wrapper around http.Client
type client struct {
	httpClient *http.Client
	headers    http.Header
	indent     string
}

// NewClient returns a new JSON Client using the default http.Client
func NewClient() Client {
	return NewWithHTTPClient(http.DefaultClient)
}

// NewWithHTTPClient returns a new JSON Client using the specified http.Client
func NewWithHTTPClient(httpClient *http.Client) Client {
	return &client{
		httpClient: httpClient,
		headers: http.Header{
			"Content-Type": []string{ContentType},
		},
	}
}

// Headers return the default headers for requests
func (c *client) Headers() http.Header {
	return c.headers
}

// Get fetches url using GET and unmarshals into the passed response
func (c *client) Get(url string, response interface{}) error {
	res, err := c.httpClient.Get(url)
	if err != nil {
		return err
	}

	return parseResponse(res, response)
}

// Post sends request as JSON and unmarshals the response JSON into the supplied struct
func (c *client) Post(url string, request interface{}, response interface{}) error {
	var err error
	var body []byte

	body, err = json.MarshalIndent(request, "", c.indent)
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
	res, err = c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending payload: %v", err)
	}

	return parseResponse(res, response)
}

func (c *client) ErrorResponse(err error, response interface{}) bool {
	jerr, isJsonError := err.(Error)
	if !isJsonError {
		return false
	}

	return json.Unmarshal([]byte(jerr.Body), response) == nil
}

func parseResponse(res *http.Response, response interface{}) error {
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if res.StatusCode >= 400 {
		err = fmt.Errorf("got HTTP %v", res.Status)
	}

	if err == nil {
		err = json.Unmarshal(body, response)
	}

	if err != nil {
		if body == nil {
			body = []byte{}
		}
		return Error{
			StatusCode: res.StatusCode,
			Body:       string(body),
			err:        err,
		}
	}

	return nil
}
