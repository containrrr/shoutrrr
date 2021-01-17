package jsonclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const ContentType = "application/json"

type Client struct{}

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

func WrapError(err error) *Error {
	return &Error{
		err: err,
	}
}

func Get(url string, response interface{}) *Error {
	//fmt.Printf("GET %v\n", url)

	res, err := http.Get(url)
	if err != nil {
		return WrapError(err)
	}

	return parseResponse(res, response)
}

func Post(url string, request interface{}, response interface{}) *Error {

	var err error
	var body []byte

	body, err = json.Marshal(request)
	if err != nil {
		return WrapError(fmt.Errorf("error creating payload: %v", err))
	}

	// fmt.Println(string(body))

	var res *http.Response
	res, err = http.Post(url, ContentType, bytes.NewReader(body))
	if err != nil {
		return WrapError(fmt.Errorf("error sending payload: %v", err))
	}

	return parseResponse(res, response)
}

func parseResponse(res *http.Response, response interface{}) *Error {
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err == nil {
		// fmt.Println(string(body))
	}

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
