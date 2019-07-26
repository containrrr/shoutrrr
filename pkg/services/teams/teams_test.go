package teams

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTeams(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Teams Suite")
}

var _ = Describe("the teams plugin", func() {
	It("should work", func() {
		config := Config{
			Token: Token{
				"88888888-4444-333-333-cccccccccccc",
				"11111111112222222222333333333344",
				"88888888-4444-333-333-cccccccccccc",
			},
		}
		// testutils.TestConfigSetInvalidQueryValue(&Config{}, "teams://88888888-4444-333-333-cccccccccccc:11111111112222222222333333333344@88888888-4444-333-333-cccccccccccc/?foo=bar")
		apiURL := buildURL(&config)
		expectedURL := "https://outlook.office.com/webhook/88888888-4444-333-333-cccccccccccc/IncomingWebhook/11111111112222222222333333333344/88888888-4444-333-333-cccccccccccc"
		Expect(apiURL).To(Equal(expectedURL))
	})
})
