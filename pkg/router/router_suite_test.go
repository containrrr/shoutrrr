package router

import (
	"log"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRouter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Router Suite")
}

var sr ServiceRouter

var _ = Describe("the router suite", func() {
	BeforeEach(func() {
		sr = ServiceRouter{}
	})

	When("extract service name is given a url", func() {
		It("should extract the protocol/service part", func() {
			url := "slack://invalid-part"
			serviceName, _, err := sr.ExtractServiceName(url)
			Expect(err).ToNot(HaveOccurred())
			Expect(serviceName).To(Equal("slack"))
		})
		It("should return an error if the protocol/service part is missing", func() {
			url := "://invalid-part"
			serviceName, _, err := sr.ExtractServiceName(url)
			Expect(err).To(HaveOccurred())
			Expect(serviceName).To(Equal(""))
		})
		It("should return an error if the protocol/service part is containing invalid letters", func() {
			url := "a d://invalid-part"
			serviceName, _, err := sr.ExtractServiceName(url)
			Expect(err).To(HaveOccurred())
			Expect(serviceName).To(Equal(""))
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
