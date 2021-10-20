package webclient

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type webclient struct {
	httpClient    *http.Client
	headers       map[string]string
	ParseResponse func(res []byte, response interface{}) error
	CreateRequest func(request interface{}) (io.Reader, error)
}

// Headers return the default headers for requests
func (c *webclient) Headers() map[string]string {
	return c.headers
}

// Get fetches url using GET and unmarshals into the passed response
func (c *webclient) Get(url string, response interface{}) error {
	res, err := c.httpClient.Get(url)
	if err != nil {
		return err
	}

	return c.parseResponse(res, response)
}

// Post sends request as JSON and unmarshals the response JSON into the supplied struct
func (c *webclient) Post(url string, request interface{}, response interface{}) error {
	var err error

	body, err := c.CreateRequest(request)

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	for key, val := range c.headers {
		req.Header.Set(key, val)
	}

	var res *http.Response
	res, err = c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending payload: %v", err)
	}

	return c.parseResponse(res, response)
}

func (c *webclient) ErrorResponse(err error, response interface{}) bool {
	return false
}

func (c *webclient) parseResponse(res *http.Response, response interface{}) error {
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if res.StatusCode >= 400 {
		err = fmt.Errorf("got HTTP %v", res.Status)
	}

	if response == nil {
		return nil
	}

	if err == nil {
		err = c.ParseResponse(body, response)
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
