package failures

import "fmt"

// FailureID is a number to be used to identify a specific error
type FailureID int

type failure struct {
	message string
	id      FailureID
	stack   string
}

// Failure is an extended error that also includes an ID to be used to identify a specific error
type Failure interface {
	error
	ID() FailureID
}

func (f *failure) Error() string {
	return fmt.Sprintf("%s: %s", f.message, f.stack)
}

func (f *failure) ID() FailureID {
	return f.id
}

// Wrap returns a failure with the given message and id, saving the message of wrappedError for appending to Error()
func Wrap(message string, id FailureID, wrappedError error, v ...interface{}) Failure {
	var stack string
	if wrappedError != nil {
		stack = wrappedError.Error()
	}

	if len(v) > 0 {
		message = fmt.Sprintf(message, v...)
	}

	return &failure{
		message: message,
		id:      id,
		stack:   stack,
	}
}

var _ error = &failure{}
