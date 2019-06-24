package queue

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("the standard queued sender implementation", func() {
	When("a message is enqueued", func() {
		It("should be added to the internal queue", func() {
			qs := &queuedSender{}
			qs.Enqueue("message body")
			Expect(qs.queue).ToNot(BeNil())
			Expect(qs.queue).To(HaveLen(1))
		})
	})
	When("a formatted message is enqueued", func() {
		It("should be added with the specified format", func() {
			qs := &queuedSender{}
			qs.Enqueuef("message with number %d", 5)
			Expect(qs.queue).ToNot(BeNil())
			Expect(qs.queue[0]).To(Equal("message with number 5"))
		})
	})
	When("it leaves the scope after flush has been deferred", func() {
		When("it hasn't been assigned a sender", func() {
			It("should not cause a panic", func() {
				qs := &queuedSender{}
				defer qs.Flush(nil)
				qs.Enqueue("message")
			})
		})
	})
})
