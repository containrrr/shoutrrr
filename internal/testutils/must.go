package testutils

import (
	"net/url"

	"github.com/jarcoal/httpmock"
	"github.com/onsi/gomega"
)

func URLMust(rawURL string) *url.URL {
	parsed, err := url.Parse(rawURL)
	gomega.ExpectWithOffset(1, err).NotTo(gomega.HaveOccurred())
	return parsed
}

func JSONRespondMust(code int, response interface{}) httpmock.Responder {
	responder, err := httpmock.NewJsonResponder(code, response)
	gomega.ExpectWithOffset(1, err).NotTo(gomega.HaveOccurred(), "invalid test response struct")
	return responder
}
