package util

import (
	"io/ioutil"
	"log"

	"github.com/onsi/ginkgo"
)

// Min returns the smallest of a and b
func Min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max returns the largest of a and b
func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

// TestLogger returns a log.Logger that writes to ginkgo.GinkgoWriter for use in tests
func TestLogger() *log.Logger {
	return log.New(ginkgo.GinkgoWriter, "Test", log.LstdFlags)
}

// DiscardLogger is a logger that discards any output written to it
var DiscardLogger = log.New(ioutil.Discard, "", 0)
