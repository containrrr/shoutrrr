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
            serviceName, err := router.ExtractServiceName(url)
            Expect(err).ToNot(HaveOccurred())
            Expect(serviceName).To(Equal("slack"))
        })
        It("should return an error if the protocol/service part is missing", func() {
            url := "://invalid-part"
            serviceName, err := router.ExtractServiceName(url)
            Expect(err).To(HaveOccurred())
            Expect(serviceName).To(Equal(""))
        })
        It("should return an error if the protocol/service part is containing non-alphabetic letters", func() {
            url := "a12d://invalid-part"
            serviceName, err := router.ExtractServiceName(url)
            Expect(err).To(HaveOccurred())
            Expect(serviceName).To(Equal(""))
        })
    })

    When("extract arguments is given a url", func() {
        It("should return the arguments", func() {
            url := "slack://aaaa"
            arguments, err := router.ExtractArguments(url)
            Expect(err).NotTo(HaveOccurred())
            Expect(len(arguments)).To(Equal(1))
        })
        It("should return an error if no arguments could be found", func() {
            url := "slack://"
            arguments, err := router.ExtractArguments(url)
            Expect(err).To(HaveOccurred())
            Expect(len(arguments)).To(Equal(0))
        })
        It("should split the arguments by /", func() {
            url := "slack://aaaa/bbb/ccc"
            arguments, err := router.ExtractArguments(url)
            Expect(err).NotTo(HaveOccurred())
            Expect(len(arguments)).To(Equal(3))
            Expect(arguments[0]).To(Equal("aaaa"))
            Expect(arguments[1]).To(Equal("bbb"))
            Expect(arguments[2]).To(Equal("ccc"))
        })
    })
})
