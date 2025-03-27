package jsonclient

import "net/http"

type Client interface {
	Get(url string, response any) error
	Post(url string, request any, response any) error
	Headers() http.Header
	ErrorResponse(err error, response any) bool
}
