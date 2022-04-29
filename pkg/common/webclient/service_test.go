package webclient_test

import (
	"github.com/containrrr/shoutrrr/pkg/common/webclient"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ClientService", func() {

	When("getting the web client from an empty service", func() {
		It("should return an initialized web client", func() {
			service := &webclient.ClientService{}
			Expect(service.WebClient()).ToNot(BeNil())
		})
	})

	When("getting the http client from an empty service", func() {
		It("should return an initialized http client", func() {
			service := &webclient.ClientService{}
			Expect(service.HTTPClient()).ToNot(BeNil())
		})
	})

	When("no certs have been added", func() {
		It("should use nil as the certificate pool", func() {
			service := &webclient.ClientService{}
			tp := service.HTTPClient().Transport.(*http.Transport)
			Expect(tp.TLSClientConfig.RootCAs).To(BeNil())
		})
	})

	When("a custom cert have been added", func() {
		It("should use a custom certificate pool", func() {
			service := &webclient.ClientService{}

			// Adding an empty cert should fail
			addedOk := service.AddTrustedRootCertificate([]byte{})
			Expect(addedOk).To(BeFalse())

			tp := service.HTTPClient().Transport.(*http.Transport)
			Expect(tp.TLSClientConfig.RootCAs).ToNot(BeNil())
		})
	})
})
