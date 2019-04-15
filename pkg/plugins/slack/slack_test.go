package slack

import (
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "testing"
)

func TestShoutrrr(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Shoutrrr Suite")
}

var plugin = &SlackPlugin{}
var _ = Describe("the slack plugin", func() {
    When("given a token with a malformed part", func () {
      It("should return an error if part A is not 9 letters", func() {
          expectErrorMessageGivenUrl(
              TokenAMalformed,
              "slack://lol/12345678/123456789/123456789123456789123456")
      })
        It("should return an error if part B is not 9 letters", func() {
            expectErrorMessageGivenUrl(
                TokenBMalformed,
                "slack://lol/123456789/12345678/123456789123456789123456")
        })
        It("should return an error if part C is not 24 letters", func() {
            expectErrorMessageGivenUrl(
                TokenCMalformed,
                "slack://123456789/123456789/12345678912345678912345")
        })
    })
    When("given a token missing a part", func () {
        It("should return an error if the missing part is A", func() {
            expectErrorMessageGivenUrl(
                TokenAMissing,
                "slack://lol//123456789/123456789123456789123456")
        })
        It("should return an error if the missing part is B", func() {
            expectErrorMessageGivenUrl(
                TokenBMissing,
                "slack://lol/123456789//123456789")

        })
        It("should return an error if the missing part is C", func() {
            expectErrorMessageGivenUrl(
                TokenCMissing,
                "slack://lol/123456789/123456789/")
        })
    })
})
func expectErrorMessageGivenUrl(msg ErrorMessage, url string){
    err := plugin.Send(url, "Hello")
    Expect(err).To(HaveOccurred())
    Expect(err.Error()).To(Equal(string(msg)))
}

