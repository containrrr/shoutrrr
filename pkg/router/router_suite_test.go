package router_test

import (
    "testing"
    . "github.com/containrrr/shoutrrr/pkg/router"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestRouter(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Router Suite")
}

var router ServiceRouter

var _ = Describe("the router suite", func() {
    BeforeEach(func() {
        router = ServiceRouter{}
    })

    When("extract service name is given a url", func() {
        It("should extract the protocol/service part", func() {
            url := "slack://invalid-part"
            serviceName, _, err := router.ExtractServiceName(url)
            Expect(err).ToNot(HaveOccurred())
            Expect(serviceName).To(Equal("slack"))
        })
        It("should return an error if the protocol/service part is missing", func() {
            url := "://invalid-part"
            serviceName, _, err := router.ExtractServiceName(url)
            Expect(err).To(HaveOccurred())
            Expect(serviceName).To(Equal(""))
        })
        It("should return an error if the protocol/service part is containing invalid letters", func() {
            url := "a d://invalid-part"
            serviceName, _, err := router.ExtractServiceName(url)
            Expect(err).To(HaveOccurred())
            Expect(serviceName).To(Equal(""))
        })
    })
})
