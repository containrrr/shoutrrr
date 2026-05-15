package router

import (
	"fmt"
	"log"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRouter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Router Suite")
}

var sr ServiceRouter

const (
	mockCustomURL = "teams+https://publicservice.info/webhook/11111111-4444-4444-8444-cccccccccccc@22222222-4444-4444-8444-cccccccccccc/IncomingWebhook/33333333012222222222333333333344/44444444-4444-4444-8444-cccccccccccc"
)

var _ = Describe("the router suite", func() {
	BeforeEach(func() {
		sr = ServiceRouter{
			logger: log.New(GinkgoWriter, "Test", log.LstdFlags),
		}
	})

	When("extract service name is given a url", func() {
		It("should extract the protocol/service part", func() {
			url := "slack://rest/of/url"
			serviceName, _, err := sr.ExtractServiceName(url)
			Expect(err).ToNot(HaveOccurred())
			Expect(serviceName).To(Equal("slack"))
		})
		It("should extract the service part when provided in custom form", func() {
			url := "teams+https://rest/of/url"
			serviceName, _, err := sr.ExtractServiceName(url)
			Expect(err).ToNot(HaveOccurred())
			Expect(serviceName).To(Equal("teams"))
		})
		It("should return an error if the protocol/service part is missing", func() {
			url := "://rest/of/url"
			serviceName, _, err := sr.ExtractServiceName(url)
			Expect(err).To(HaveOccurred())
			Expect(serviceName).To(Equal(""))
		})
		It("should return an error if the protocol/service part is containing invalid letters", func() {
			url := "a d://rest/of/url"
			serviceName, _, err := sr.ExtractServiceName(url)
			Expect(err).To(HaveOccurred())
			Expect(serviceName).To(Equal(""))
		})
	})

	When("initializing a service with a custom URL", func() {
		It("should return an error if the service does not support it", func() {
			service, scheme, err := sr.initService("logger+https://hybr.is")
			Expect(err).To(HaveOccurred())
			Expect(scheme).To(Equal("logger"))
			Expect(service).To(BeNil())
		})
	})

	When("listing added services", func() {
		When("multiple instances of the same service have been added", func() {
			It("should return a list with unique identifiers for those services", func() {
				Expect(sr.AddService("logger://")).To(Succeed())
				Expect(sr.AddService("logger://")).To(Succeed())
				Expect(sr.AddService("logger://")).To(Succeed())
				Expect(sr.ListAddedServices()).To(ConsistOf("logger", "logger2", "logger3"))
			})
		})
	})

	Describe("the service map", func() {
		When("resolving implemented services", func() {
			services := (&ServiceRouter{}).ListServices()

			for _, scheme := range services {
				// copy ref to local closure
				serviceScheme := scheme

				It(fmt.Sprintf("should return a Service for '%s'", serviceScheme), func() {
					service, err := newService(serviceScheme)

					Expect(err).NotTo(HaveOccurred())
					Expect(service).ToNot(BeNil())
				})
			}
		})
	})

	When("initializing a service with a custom URL", func() {
		It("should return an error if the service does not support it", func() {
			service, scheme, err := sr.initService("logger+https://hybr.is")
			Expect(err).To(HaveOccurred())
			Expect(scheme).To(Equal("logger"))
			Expect(service).To(BeNil())
		})
		It("should successfully init a service that does support it", func() {
			service, scheme, err := sr.initService(mockCustomURL)
			Expect(err).NotTo(HaveOccurred())
			Expect(scheme).To(Equal("teams"))
			Expect(service).NotTo(BeNil())
		})
	})

	When("a message is enqueued", func() {
		It("should be added to the internal queue", func() {

			sr.Enqueue("message body")
			Expect(sr.queue).ToNot(BeNil())
			Expect(sr.queue).To(HaveLen(1))
		})
	})
	When("a formatted message is enqueued", func() {
		It("should be added with the specified format", func() {
			sr.Enqueue("message with number %d", 5)
			Expect(sr.queue).ToNot(BeNil())
			Expect(sr.queue[0]).To(Equal("message with number 5"))
		})
	})
	When("it leaves the scope after flush has been deferred", func() {
		When("it hasn't been assigned a sender", func() {
			It("should not cause a panic", func() {
				defer sr.Flush(nil)
				sr.Enqueue("message")
			})
		})
	})
	When("router has not been provided a logger", func() {
		It("should not crash when trying to log", func() {
			router := ServiceRouter{}
			_, _, err := router.initService(mockCustomURL)
			Expect(err).NotTo(HaveOccurred())
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
