package webclient

import (
	"encoding/json"
	"net/http"
)

// JSONContentType is the default mime type for JSON
const JSONContentType = "application/json"

// DefaultJSONClient is the singleton instance of WebClient using http.DefaultClient
var DefaultJSONClient = NewJSONClient()

// GetJSON fetches url using GET and unmarshals into the passed response using DefaultJSONClient
func GetJSON(url string, response interface{}) error {
	return DefaultJSONClient.Get(url, response)
}

// PostJSON sends request as JSON and unmarshals the response JSON into the supplied struct using DefaultJSONClient
func PostJSON(url string, request interface{}, response interface{}) error {
	return DefaultJSONClient.Post(url, request, response)
}

// NewJSONClient returns a WebClient using the default http.Client and JSON serialization
func NewJSONClient() WebClient {
	var c client
	c = client{
		headers: http.Header{
			"Content-Type": []string{JSONContentType},
		},
		parse: json.Unmarshal,
		write: func(v interface{}) ([]byte, error) {
			return json.MarshalIndent(v, "", c.indent)
		},
	}
	return &c
}
