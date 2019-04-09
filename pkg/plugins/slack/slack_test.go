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
          expectErrorMessageGivenToken(
              TokenAMalformed,
              SlackToken{
                  A: "12345678",
                  B: "123456789",
                  C: "123456789123456789123456",
              })
      })
        It("should return an error if part B is not 9 letters", func() {
            expectErrorMessageGivenToken(
                TokenBMalformed,
                SlackToken{
                    A: "123456789",
                    B: "12345678",
                    C: "123456789123456789123456",
                })
        })
        It("should return an error if part C is not 24 letters", func() {
            expectErrorMessageGivenToken(
                TokenCMalformed,
                SlackToken{
                    A: "123456789",
                    B: "123456789",
                    C: "12345678912345678912345",
                })
        })
    })
    When("given a token missing a part", func () {
        It("should return an error if the missing part is A", func() {
            expectErrorMessageGivenToken(
                TokenAMissing,
                SlackToken{
                    B: "123456789",
                    C: "123456789123456789123456",
                })
        })
        It("should return an error if the missing part is B", func() {
            expectErrorMessageGivenToken(
                TokenBMissing,
                SlackToken{
                    A: "123456789",
                    C: "123456789123456789123456",
                })
        })
        It("should return an error if the missing part is C", func() {
            expectErrorMessageGivenToken(
                TokenCMissing,
                SlackToken{
                    A: "123456789",
                    B: "123456789",
                })
        })
    })
})
func expectErrorMessageGivenToken(msg SlackErrorMessage, token SlackToken){
    config := SlackConfig{
        Token: token,
    }
    err := plugin.Send(config, "Hello")
    Expect(err).To(HaveOccurred())
    Expect(err.Error()).To(Equal(string(msg)))
}

