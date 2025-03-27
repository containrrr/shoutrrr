package failures

import "fmt"

// FailureID is a unique identifier for a specific error type.
type FailureID int

// failure is the concrete implementation of the Failure interface.
// It wraps an error with a message and an ID for categorization.
type failure struct {
	message string    // Descriptive message for the error
	id      FailureID // Unique identifier for this error type
	wrapped error     // Underlying error, if any, for chaining
}

// Failure extends the error interface with an ID and methods for unwrapping and comparison.
// It allows errors to be identified by a unique ID and supports Go’s error wrapping conventions.
type Failure interface {
	error
	ID() FailureID        // Returns the unique identifier for this failure
	Unwrap() error        // Returns the wrapped error, if any
	Is(target error) bool // Checks if the target error matches this failure by ID
}

// Error returns the failure’s message, appending the wrapped error’s message if present.
func (f *failure) Error() string {
	if f.wrapped == nil {
		return f.message
	}

	return fmt.Sprintf("%s: %v", f.message, f.wrapped)
}

// Unwrap returns the underlying error wrapped by this failure, or nil if none exists.
func (f *failure) Unwrap() error {
	return f.wrapped
}

// ID returns the unique identifier assigned to this failure.
func (f *failure) ID() FailureID {
	return f.id
}

// Is reports whether the target error is a failure with the same ID.
// It only returns true for failures of the same type with matching IDs.
func (f *failure) Is(target error) bool {
	targetFailure, ok := target.(*failure)

	return ok && targetFailure.id == f.id
}

// Wrap creates a new failure with the given message, ID, and optional wrapped error.
// If variadic arguments are provided, they are used to format the message using fmt.Sprintf.
// This supports Go’s error wrapping pattern while adding a unique ID for identification.
func Wrap(message string, id FailureID, wrappedError error, v ...any) Failure {
	if len(v) > 0 {
		message = fmt.Sprintf(message, v...)
	}

	return &failure{
		message: message,
		id:      id,
		wrapped: wrappedError,
	}
}

// Ensure failure implements the error interface at compile time.
var _ error = &failure{}
