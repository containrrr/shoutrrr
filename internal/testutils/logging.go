package testutils

import (
	"log"

	"github.com/onsi/ginkgo/v2"
)

// TestLogger returns a log.Logger that writes to ginkgo.GinkgoWriter for use in tests.
func TestLogger() *log.Logger {
	return log.New(ginkgo.GinkgoWriter, "[Test] ", 0)
}
