package telegram_test

import (
	"github.com/containrrr/shoutrrr/pkg/common/webclient"
	"github.com/containrrr/shoutrrr/pkg/services/telegram"
	"github.com/jarcoal/httpmock"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var client *telegram.Client

var _ = Describe("the telegram client", func() {
	BeforeEach(func() {
		client = &telegram.Client{WebClient: webclient.NewJSONClient(), Token: `Test`}
		httpmock.ActivateNonDefault(client.WebClient.HTTPClient())
	})
	AfterEach(func() {
		httpmock.DeactivateAndReset()
	})
	When("an error is returned from the API", func() {
		It("should return the error description", func() {

			errRes := httpmock.NewJsonResponderOrPanic(http.StatusNotAcceptable, telegram.ErrorResponse{
				OK:          false,
				Description: "no.",
			})
			httpmock.RegisterResponder("POST", `https://api.telegram.org/botTest/getUpdates`, errRes)
			httpmock.RegisterResponder("GET", `https://api.telegram.org/botTest/getMe`, errRes)
			httpmock.RegisterResponder("POST", `https://api.telegram.org/botTest/sendMessage`, errRes)

			_, err := client.GetUpdates(0, 1, 10, []string{})
			Expect(err).To(MatchError(`no.`))

			_, err = client.GetBotInfo()
			Expect(err).To(MatchError(`no.`))

			_, err = client.SendMessage(&telegram.SendMessagePayload{})
			Expect(err).To(MatchError(`no.`))
		})
	})

})
