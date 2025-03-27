package router

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestRouter(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Router Suite")
}

var sr ServiceRouter

const (
	mockCustomURL = "teams+https://publicservice.webhook.office.com/webhookb2/11111111-4444-4444-8444-cccccccccccc@22222222-4444-4444-8444-cccccccccccc/IncomingWebhook/33333333012222222222333333333344/44444444-4444-4444-8444-cccccccccccc/V2ESyij_gAljSoUQHvZoZYzlpAoAXExyOl26dlf1xHEx05?host=publicservice.webhook.office.com"
)

var _ = ginkgo.Describe("the router suite", func() {
	ginkgo.BeforeEach(func() {
		sr = ServiceRouter{
			logger: log.New(ginkgo.GinkgoWriter, "Test", log.LstdFlags),
		}
	})

	ginkgo.When("extract service name is given a url", func() {
		ginkgo.It("should extract the protocol/service part", func() {
			url := "slack://rest/of/url"
			serviceName, _, err := sr.ExtractServiceName(url)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(serviceName).To(gomega.Equal("slack"))
		})
		ginkgo.It("should extract the service part when provided in custom form", func() {
			url := "teams+https://rest/of/url"
			serviceName, _, err := sr.ExtractServiceName(url)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(serviceName).To(gomega.Equal("teams"))
		})
		ginkgo.It("should return an error if the protocol/service part is missing", func() {
			url := "://rest/of/url"
			serviceName, _, err := sr.ExtractServiceName(url)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(serviceName).To(gomega.Equal(""))
		})
		ginkgo.It(
			"should return an error if the protocol/service part is containing invalid letters",
			func() {
				url := "a d://rest/of/url"
				serviceName, _, err := sr.ExtractServiceName(url)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(serviceName).To(gomega.Equal(""))
			},
		)
	})

	ginkgo.When("initializing a service with a custom URL", func() {
		ginkgo.It("should return an error if the service does not support it", func() {
			service, err := sr.initService("log+https://hybr.is")
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(service).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("the service map", func() {
		ginkgo.When("resolving implemented services", func() {
			services := (&ServiceRouter{}).ListServices()

			for _, scheme := range services {
				// copy ref to local closure
				serviceScheme := scheme

				ginkgo.It(fmt.Sprintf("should return a Service for '%s'", serviceScheme), func() {
					service, err := newService(serviceScheme)

					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(service).ToNot(gomega.BeNil())
				})
			}
		})
	})

	ginkgo.When("initializing a service with a custom URL", func() {
		ginkgo.It("should return an error if the service does not support it", func() {
			service, err := sr.initService("log+https://hybr.is")
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(service).To(gomega.BeNil())
		})
		ginkgo.It("should successfully init a service that does support it", func() {
			service, err := sr.initService(mockCustomURL)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(service).NotTo(gomega.BeNil())
		})
	})

	ginkgo.When("a message is enqueued", func() {
		ginkgo.It("should be added to the internal queue", func() {
			sr.Enqueue("message body")
			gomega.Expect(sr.queue).ToNot(gomega.BeNil())
			gomega.Expect(sr.queue).To(gomega.HaveLen(1))
		})
	})
	ginkgo.When("a formatted message is enqueued", func() {
		ginkgo.It("should be added with the specified format", func() {
			sr.Enqueue("message with number %d", 5)
			gomega.Expect(sr.queue).ToNot(gomega.BeNil())
			gomega.Expect(sr.queue[0]).To(gomega.Equal("message with number 5"))
		})
	})
	ginkgo.When("it leaves the scope after flush has been deferred", func() {
		ginkgo.When("it hasn't been assigned a sender", func() {
			ginkgo.It("should not cause a panic", func() {
				defer sr.Flush(nil)
				sr.Enqueue("message")
			})
		})
	})
	ginkgo.When("router has not been provided a logger", func() {
		ginkgo.It("should not crash when trying to log", func() {
			router := ServiceRouter{}
			_, err := router.initService(mockCustomURL)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
	})
})

func ExampleNew() {
	logger := log.New(os.Stdout, "", 0)

	sr, err := New(logger, "logger://")
	if err != nil {
		log.Fatalf("could not create router: %s", err)
	}

	sr.Send("hello", nil)
	// Output: hello
}

func ExampleServiceRouter_Enqueue() {
	logger := log.New(os.Stdout, "", 0)

	sr, err := New(logger, "logger://")
	if err != nil {
		log.Fatalf("could not create router: %s", err)
	}

	defer sr.Flush(nil)
	sr.Enqueue("hello")
	sr.Enqueue("world")
	// Output:
	// hello
	// world
}
