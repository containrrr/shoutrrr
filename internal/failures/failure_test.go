package failures_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/nicholas-fedor/shoutrrr/internal/failures"
	"github.com/nicholas-fedor/shoutrrr/internal/testutils"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
)

// TestFailures runs the Ginkgo test suite for the failures package.
func TestFailures(t *testing.T) {
	format.CharactersAroundMismatchToInclude = 20 // Show more context in failure output

	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Failure Suite")
}

var _ = ginkgo.Describe("the failure package", func() {
	// Common test fixtures
	var (
		testID      failures.FailureID = 42                             // Consistent ID for testing
		testMessage string             = "test failure occurred"        // Sample error message
		wrappedErr  error              = errors.New("underlying error") // Sample wrapped error
	)

	ginkgo.Describe("Wrap function", func() {
		ginkgo.When("creating a basic failure", func() {
			ginkgo.It("returns a failure with the provided message and ID", func() {
				failure := failures.Wrap(testMessage, testID, nil)
				gomega.Expect(failure.Error()).To(gomega.Equal(testMessage))
				gomega.Expect(failure.ID()).To(gomega.Equal(testID))
				gomega.Expect(failure.Unwrap()).To(gomega.BeNil())
			})
		})

		ginkgo.When("wrapping an existing error", func() {
			ginkgo.It("combines the message and wrapped error", func() {
				failure := failures.Wrap(testMessage, testID, wrappedErr)
				expectedError := fmt.Sprintf("%s: %v", testMessage, wrappedErr)
				gomega.Expect(failure.Error()).To(gomega.Equal(expectedError))
				gomega.Expect(failure.ID()).To(gomega.Equal(testID))
				gomega.Expect(failure.Unwrap()).To(gomega.Equal(wrappedErr))
			})
		})

		ginkgo.When("using formatted message with arguments", func() {
			ginkgo.It("formats the message correctly", func() {
				formatMessage := "test failure %d"
				failure := failures.Wrap(formatMessage, testID, nil, 123)
				gomega.Expect(failure.Error()).To(gomega.Equal("test failure 123"))
				gomega.Expect(failure.ID()).To(gomega.Equal(testID))
			})
		})
	})

	ginkgo.Describe("Failure interface methods", func() {
		var failure failures.Failure

		// Setup a failure with a wrapped error before each test
		ginkgo.BeforeEach(func() {
			failure = failures.Wrap(testMessage, testID, wrappedErr)
		})

		ginkgo.Describe("Error method", func() {
			ginkgo.It("returns only the message when no wrapped error exists", func() {
				failureNoWrap := failures.Wrap(testMessage, testID, nil)
				gomega.Expect(failureNoWrap.Error()).To(gomega.Equal(testMessage))
			})
			ginkgo.It("combines message with wrapped error", func() {
				expected := fmt.Sprintf("%s: %v", testMessage, wrappedErr)
				gomega.Expect(failure.Error()).To(gomega.Equal(expected))
			})
		})

		ginkgo.Describe("ID method", func() {
			ginkgo.It("returns the assigned ID", func() {
				gomega.Expect(failure.ID()).To(gomega.Equal(testID))
			})
		})

		ginkgo.Describe("Unwrap method", func() {
			ginkgo.It("returns the wrapped error", func() {
				gomega.Expect(failure.Unwrap()).To(gomega.Equal(wrappedErr))
			})
			ginkgo.It("returns nil when no wrapped error exists", func() {
				failureNoWrap := failures.Wrap(testMessage, testID, nil)
				gomega.Expect(failureNoWrap.Unwrap()).To(gomega.BeNil())
			})
		})

		ginkgo.Describe("Is method", func() {
			ginkgo.It("returns true for failures with the same ID", func() {
				f1 := failures.Wrap("first", testID, nil)
				f2 := failures.Wrap("second", testID, nil)
				gomega.Expect(f1.Is(f2)).To(gomega.BeTrue())
				gomega.Expect(f2.Is(f1)).To(gomega.BeTrue())
			})
			ginkgo.It("returns false for failures with different IDs", func() {
				f1 := failures.Wrap("first", testID, nil)
				f2 := failures.Wrap("second", testID+1, nil)
				gomega.Expect(f1.Is(f2)).To(gomega.BeFalse())
				gomega.Expect(f2.Is(f1)).To(gomega.BeFalse())
			})
			ginkgo.It("returns false when comparing with a non-failure error", func() {
				f1 := failures.Wrap("first", testID, nil)
				gomega.Expect(f1.Is(wrappedErr)).To(gomega.BeFalse())
			})
		})
	})

	ginkgo.Describe("edge cases", func() {
		ginkgo.When("wrapping with an empty message", func() {
			ginkgo.It("handles an empty message gracefully", func() {
				failure := failures.Wrap("", testID, wrappedErr)
				gomega.Expect(failure.Error()).To(gomega.Equal(": " + wrappedErr.Error()))
				gomega.Expect(failure.ID()).To(gomega.Equal(testID))
				gomega.Expect(failure.Unwrap()).To(gomega.Equal(wrappedErr))
			})
		})

		ginkgo.When("wrapping with nil error and no args", func() {
			ginkgo.It("returns a valid failure with just message and ID", func() {
				failure := failures.Wrap(testMessage, testID, nil)
				gomega.Expect(failure.Error()).To(gomega.Equal(testMessage))
				gomega.Expect(failure.ID()).To(gomega.Equal(testID))
				gomega.Expect(failure.Unwrap()).To(gomega.BeNil())
			})
		})

		ginkgo.When("using multiple wrapped failures", func() {
			ginkgo.It("correctly chains and unwraps multiple errors", func() {
				innerErr := errors.New("inner error")
				middleErr := failures.Wrap("middle", testID+1, innerErr)
				outerErr := failures.Wrap("outer", testID, middleErr)
				gomega.Expect(outerErr.Error()).To(gomega.Equal("outer: middle: inner error"))
				gomega.Expect(outerErr.ID()).To(gomega.Equal(testID))
				gomega.Expect(outerErr.Unwrap()).To(gomega.Equal(middleErr))
				gomega.Expect(middleErr.Unwrap()).To(gomega.Equal(innerErr))
			})
		})
	})

	ginkgo.Describe("integration-like scenarios", func() {
		ginkgo.It("works with standard error wrapping utilities", func() {
			innerErr := errors.New("inner error")
			failure := failures.Wrap("wrapped failure", testID, innerErr)
			gomega.Expect(errors.Is(failure, innerErr)).To(gomega.BeTrue()) // Matches wrapped error
			gomega.Expect(errors.Unwrap(failure)).To(gomega.Equal(innerErr))
		})

		ginkgo.It("handles fmt.Errorf wrapping", func() {
			failure := failures.Wrap("failure", testID, nil)
			wrapped := fmt.Errorf("additional context: %w", failure)
			gomega.Expect(wrapped.Error()).To(gomega.Equal("additional context: failure"))
			gomega.Expect(errors.Unwrap(wrapped)).To(gomega.Equal(failure))
		})
	})

	ginkgo.Describe("testutils integration", func() {
		ginkgo.It("can use TestLogger for logging failures", func() {
			// Demonstrate compatibility with testutils logger
			failure := failures.Wrap("logged failure", testID, nil)
			logger := testutils.TestLogger()
			logger.Printf("Error occurred: %v", failure)
			// No assertion needed; ensures no panic during logging
		})
	})
})
