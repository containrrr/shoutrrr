package testutils

import (
	"log"

	"github.com/onsi/ginkgo"
)

// TestLogger returns a log.Logger that writes to ginkgo.GinkgoWriter for use in tests
func TestLogger() *log.Logger {
	return log.New(ginkgo.GinkgoWriter, "Test", log.LstdFlags)
}
