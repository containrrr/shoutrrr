package failures

import "fmt"

// FailureID is a number to be used to identify a specific error
type FailureID int

type failure struct {
	message string
	id      FailureID
	wrapped error
}

// Failure is an extended error that also includes an ID to be used to identify a specific error
type Failure interface {
	error
	ID() FailureID
}

func (f *failure) Error() string {
	if f.wrapped == nil {
		return f.message
	}
	return fmt.Sprintf("%s: %v", f.message, f.wrapped)
}

func (f *failure) Unwrap() error {
	return f.wrapped
}

func (f *failure) ID() FailureID {
	return f.id
}

func (f *failure) Is(target error) bool {
	targetFailure, targetIsFailure := target.(*failure)
	return targetIsFailure && targetFailure.id == f.id
}

// Wrap returns a failure with the given message and id, saving the message of wrappedError for appending to Error()
func Wrap(message string, id FailureID, wrappedError error, v ...interface{}) Failure {

	if len(v) > 0 {
		message = fmt.Sprintf(message, v...)
	}

	return &failure{
		message: message,
		id:      id,
		wrapped: wrappedError,
	}
}

var _ error = &failure{}
