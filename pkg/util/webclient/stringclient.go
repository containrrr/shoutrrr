package webclient

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// DefaultClient is the singleton instance of jsonclient using http.DefaultClient
var DefaultStringClient = NewStringClient()

// GetString fetches url using GET and sets the response to the passed string pointer
func GetString(url string, response *string) error {
	return DefaultStringClient.Get(url, response)
}

// PostString sends request string and sets the response to the passed string pointer
func PostString(url string, request string, response *string) error {
	return DefaultStringClient.Post(url, request, response)
}

func PostUrl(url string, request url.Values, response *string) error {
	return DefaultStringClient.Post(url, request.Encode(), response)
}

func NewStringClient() Client {
	return &webclient{
		httpClient: http.DefaultClient,
		headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		ParseResponse: func(res []byte, response interface{}) error {
			resStr, ok := response.(*string)
			if !ok {
				return fmt.Errorf("response is not a string pointer")
			}
			*resStr = string(res)
			return nil
		},
		CreateRequest: func(request interface{}) (io.Reader, error) {
			reqStr, ok := request.(string)
			if !ok {
				return nil, fmt.Errorf("request is not a string")
			}
			return strings.NewReader(reqStr), nil
		},
	}
}
