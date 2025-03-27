package testutils

import (
	"net/url"

	"github.com/jarcoal/httpmock"
	"github.com/onsi/gomega"
)

// URLMust creates a url.URL from the given rawURL and fails the test if it cannot be parsed.
func URLMust(rawURL string) *url.URL {
	parsed, err := url.Parse(rawURL)
	gomega.ExpectWithOffset(1, err).NotTo(gomega.HaveOccurred())

	return parsed
}

// JSONRespondMust creates a httpmock.Responder with the given response
// as the body, and fails the test if it cannot be created.
func JSONRespondMust(code int, response interface{}) httpmock.Responder {
	responder, err := httpmock.NewJsonResponder(code, response)
	gomega.ExpectWithOffset(1, err).NotTo(gomega.HaveOccurred(), "invalid test response struct")

	return responder
}
