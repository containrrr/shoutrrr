package standard

import (
	"fmt"
	f "github.com/containrrr/shoutrrr/internal/failures"
)

const (
	// FailTestSetup is the FailureID used to represent an error that is part of the setup for tests
	FailTestSetup f.FailureID = -1
	// FailParseURL is the FailureID used to represent failing to parse the service URL
	FailParseURL f.FailureID = -2
	// FailServiceInit is the FailureID used to represent failure of a service.Initialize method
	FailServiceInit f.FailureID = -3
	// FailUnknown is the default FailureID
	FailUnknown f.FailureID = iota
)

// Failure creates a Failure instance corresponding to the provided failureID, wrapping the provided error
func Failure(failureID f.FailureID, err error, v ...interface{}) f.Failure {
	messages := map[int]string{
		int(FailParseURL): "error parsing Service URL",
		int(FailUnknown):  "an unknown error occurred",
	}

	msg := messages[int(failureID)]
	if msg == "" {
		msg = messages[int(FailUnknown)]
	}

	return f.Wrap(msg, failureID, err, v...)
}

type failureLike interface {
	f.Failure
}

// IsTestSetupFailure checks whether the given failure is due to the test setup being broken
func IsTestSetupFailure(failure failureLike) (string, bool) {
	if failure != nil && failure.ID() == FailTestSetup {
		return fmt.Sprintf("test setup failed: %s", failure.Error()), true
	}
	return "", false
}
