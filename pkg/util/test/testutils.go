package test

import (
	"log"
	"net/url"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

// TestLogger returns a log.Logger that writes to ginkgo.GinkgoWriter for use in tests
func TestLogger() *log.Logger {
	return log.New(ginkgo.GinkgoWriter, "Test", log.LstdFlags)
}

func URLMust(rawURL string) *url.URL {
	parsed, err := url.Parse(rawURL)
	gomega.ExpectWithOffset(1, err).NotTo(gomega.HaveOccurred())
	return parsed
}
