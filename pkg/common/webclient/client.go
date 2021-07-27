package webclient

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// client is a wrapper around http.Client with common notification service functionality
type client struct {
	headers    http.Header
	indent     string
	HttpClient http.Client
	parse      ParserFunc
	write      WriterFunc
}

// SetParser overrides the parser for the incoming response content
func (c *client) SetParser(parse ParserFunc) {
	c.parse = parse
}

// SetWriter overrides the writer for the outgoing request content
func (c *client) SetWriter(write WriterFunc) {
	c.write = write
}

// Headers return the default headers for requests
func (c *client) Headers() http.Header {
	return c.headers
}

// Get fetches url using GET and unmarshals into the passed response
func (c *client) Get(url string, response interface{}) error {
	return c.request(url, response, nil)
}

// Post sends a serialized representation of request and deserializes the result into response
func (c *client) Post(url string, request interface{}, response interface{}) error {
	body, err := c.write(request)
	if err != nil {
		return fmt.Errorf("error creating payload: %v", err)
	}

	return c.request(url, response, bytes.NewReader(body))
}

// ErrorResponse tries to deserialize any response body into the supplied struct, returning whether successful or not
func (c *client) ErrorResponse(err error, response interface{}) bool {
	jerr, isWebError := err.(ClientError)
	if !isWebError {
		return false
	}

	return c.parse([]byte(jerr.Body), response) == nil
}

func (c *client) request(url string, response interface{}, payload io.Reader) error {
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return err
	}

	for key, val := range c.headers {
		req.Header.Set(key, val[0])
	}

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending payload: %v", err)
	}

	return c.parseResponse(res, response)
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
