package pinglet

import "fmt"

// pushPayload is the JSON body posted to a Pinglet topic
type pushPayload struct {
	Message  string            `json:"message"`
	Title    string            `json:"title,omitempty"`
	Priority string            `json:"priority"`
	Badges   map[string]string `json:"badges,omitempty"`
	Data     map[string]string `json:"data,omitempty"`
}

// apiResponse is the error body returned by the Pinglet API on a failure
type apiResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"error"`
}

func (e *apiResponse) Error() string {
	return fmt.Sprintf("server response: %v (%v)", e.Message, e.Code)
}
