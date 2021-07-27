package webclient

import (
	"encoding/json"
	"net/http"
)

// JsonContentType is the default mime type for JSON
const JsonContentType = "application/json"

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
