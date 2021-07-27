package webclient

import (
	"net/http"
)

// WebClient ...
type WebClient interface {
	Get(url string, response interface{}) error
	Post(url string, request interface{}, response interface{}) error
	Headers() http.Header
	ErrorResponse(err error, response interface{}) bool
	SetTransport(http.RoundTripper)
	SetParser(ParserFunc)
	SetWriter(WriterFunc)
}
