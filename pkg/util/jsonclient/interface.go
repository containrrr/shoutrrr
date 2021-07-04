package jsonclient

import "net/http"

type Client interface {
	Get(url string, response interface{}) error
	Post(url string, request interface{}, response interface{}) error
	Headers() http.Header
	ErrorResponse(err error, response interface{}) bool
}
