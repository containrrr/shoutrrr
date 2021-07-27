package webclient

import "fmt"

// ClientError contains additional http/JSON details
type ClientError struct {
	StatusCode int
	Body       string
	err        error
}

func (je ClientError) Error() string {
	return je.String()
}

func (je ClientError) String() string {
	if je.err == nil {
		return fmt.Sprintf("unknown error (HTTP %v)", je.StatusCode)
	}
	return je.err.Error()
}

// ErrorBody returns the request body from a ClientError
func ErrorBody(e error) string {
	if jsonError, ok := e.(ClientError); ok {
		return jsonError.Body
	}
	return ""
}
