package failures

import "fmt"

// FailureID is a number to be used to identify a specific error
type FailureID int

const (
	// FailTestSetup is FailureID used to represent an error that is part of the setup for tests
	FailTestSetup FailureID = -1
)

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

// IsTestSetupFailure checks whether the given failure is due to the test setup being broken
func IsTestSetupFailure(f Failure) (string, bool) {
	if f != nil && f.ID() == FailTestSetup {
		return fmt.Sprintf("test setup failed: %s", f.Error()), true
	}
	return "", false
}

var _ error = &failure{}
