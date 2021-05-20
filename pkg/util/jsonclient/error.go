package jsonclient

import "fmt"

// Error contains additional http/JSON details
type Error struct {
	StatusCode int
	Body       string
	err        error
}

func (je Error) Error() string {
	return je.String()
}

func (je Error) String() string {
	if je.err == nil {
		return fmt.Sprintf("unknown error (HTTP %v)", je.StatusCode)
	}
	return je.err.Error()
}

// ErrorBody returns the request body from an Error
func ErrorBody(e error) string {
	if jsonError, ok := e.(Error); ok {
		return jsonError.Body
	}
	return ""
}
