package standard

import (
	"fmt"

	"github.com/nicholas-fedor/shoutrrr/internal/failures"
)

const (
	// FailTestSetup is the FailureID used to represent an error that is part of the setup for tests.
	FailTestSetup failures.FailureID = -1
	// FailParseURL is the FailureID used to represent failing to parse the service URL.
	FailParseURL failures.FailureID = -2
	// FailServiceInit is the FailureID used to represent failure of a service.Initialize method.
	FailServiceInit failures.FailureID = -3
	// FailUnknown is the default FailureID.
	FailUnknown failures.FailureID = iota
)

// Failure creates a Failure instance corresponding to the provided failureID, wrapping the provided error.
func Failure(failureID failures.FailureID, err error, v ...any) failures.Failure {
	messages := map[int]string{
		int(FailTestSetup): "test setup failed",
		int(FailParseURL):  "error parsing Service URL",
		int(FailUnknown):   "an unknown error occurred",
	}

	msg := messages[int(failureID)]
	if msg == "" {
		msg = messages[int(FailUnknown)]
	}

	// If variadic arguments are provided, format them correctly
	if len(v) > 0 {
		if format, ok := v[0].(string); ok && len(v) > 1 {
			// Treat the first argument as a format string and the rest as its arguments
			extraMsg := fmt.Sprintf(format, v[1:]...)
			msg = fmt.Sprintf("%s %s", msg, extraMsg)
		} else {
			// If no format string is provided, just append the arguments as-is
			msg = fmt.Sprintf("%s %v", msg, v)
		}
	}

	return failures.Wrap(msg, failureID, err)
}

type failureLike interface {
	failures.Failure
}

// IsTestSetupFailure checks whether the given failure is due to the test setup being broken.
func IsTestSetupFailure(failure failureLike) (string, bool) {
	if failure != nil && failure.ID() == FailTestSetup {
		return fmt.Sprintf("test setup failed: %s", failure.Error()), true
	}

	return "", false
}
