package util

import (
	"github.com/onsi/gomega"
	"io/ioutil"
	"log"
	"math"
	"net/url"

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

// CeilDiv returns the quotient from dividing the dividend with the divisor, but rounded up to the nearest integer
func CeilDiv(dividend int, divisor int) int {
	return int(math.Ceil(float64(dividend) / float64(divisor)))
}

// TestLogger returns a log.Logger that writes to ginkgo.GinkgoWriter for use in tests
func TestLogger() *log.Logger {
	return log.New(ginkgo.GinkgoWriter, "Test", log.LstdFlags)
}

// DiscardLogger is a logger that discards any output written to it
var DiscardLogger = log.New(ioutil.Discard, "", 0)

// URLMust parses the specified URL and panics (with offset) if it fails
func URLMust(rawURL string) *url.URL {
	parsed, err := url.Parse(rawURL)
	gomega.ExpectWithOffset(1, err).NotTo(gomega.HaveOccurred())
	return parsed
}
