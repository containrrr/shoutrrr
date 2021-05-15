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

// Error contains additional http/JSON details
type Error struct {
	StatusCode int
	Body       string
	err        error
}

func (je *Error) Error() string {
	return je.err.Error()
}

func (je *Error) String() string {
	return je.err.Error()
}

// ErrorBody returns the request body from an Error
func ErrorBody(e error) string {
	if jsonError, ok := e.(*Error); ok {
		return jsonError.Body
	}
	return ""
}

// Get fetches url using GET and unmarshals into the passed response
func Get(url string, response interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	return parseResponse(res, response)
}

// Post sends request as JSON and unmarshals the response JSON into the supplied struct
func Post(url string, request interface{}, response interface{}) error {

	var err error
	var body []byte

	body, err = json.Marshal(request)
	if err != nil {
		return fmt.Errorf("error creating payload: %v", err)
	}

	var res *http.Response
	res, err = http.Post(url, ContentType, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error sending payload: %v", err)
	}

	return parseResponse(res, response)
}

func parseResponse(res *http.Response, response interface{}) *Error {
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
		return &Error{
			StatusCode: res.StatusCode,
			Body:       string(body),
			err:        err,
		}
	}

	return nil
}
