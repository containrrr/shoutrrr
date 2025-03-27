package telegram_test

import (
	"fmt"
	"io"
	"strings"

	"github.com/jarcoal/httpmock"
	"github.com/mattn/go-colorable"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"

	"github.com/nicholas-fedor/shoutrrr/pkg/services/telegram"
)

const (
	mockToken   = `0:MockToken`
	mockAPIBase = "https://api.telegram.org/bot" + mockToken + "/"
)

var (
	userOut    *gbytes.Buffer
	userIn     *gbytes.Buffer
	userInMono io.Writer
)

func mockTyped(a ...any) {
	_, _ = fmt.Fprint(userOut, a...)
	_, _ = fmt.Fprint(userOut, "\n")
}

func dumpBuffers() {
	for _, line := range strings.Split(string(userIn.Contents()), "\n") {
		println(">", line)
	}

	for _, line := range strings.Split(string(userOut.Contents()), "\n") {
		println("<", line)
	}
}

func mockAPI(endpoint string) string {
	return mockAPIBase + endpoint
}

var _ = ginkgo.Describe("TelegramGenerator", func() {
	ginkgo.BeforeEach(func() {
		userOut = gbytes.NewBuffer()
		userIn = gbytes.NewBuffer()
		userInMono = colorable.NewNonColorable(userIn)
		httpmock.Activate()
	})
	ginkgo.AfterEach(func() {
		httpmock.DeactivateAndReset()
	})
	ginkgo.It("should return the ", func() {
		gen := telegram.Generator{
			Reader: userOut,
			Writer: userInMono,
		}

		resultChannel := make(chan string, 1)

		httpmock.RegisterResponder("GET", mockAPI(`getMe`), httpmock.NewJsonResponderOrPanic(200, &struct {
			OK     bool
			Result *telegram.User
		}{
			true, &telegram.User{
				ID:       1,
				IsBot:    true,
				Username: "mockbot",
			},
		}))

		httpmock.RegisterResponder("POST", mockAPI(`getUpdates`), httpmock.NewJsonResponderOrPanic(200, &struct {
			OK     bool
			Result []telegram.Update
		}{
			true,
			[]telegram.Update{
				{
					ChatMemberUpdate: &telegram.ChatMemberUpdate{
						Chat:          &telegram.Chat{Type: `channel`, Title: `mockChannel`},
						OldChatMember: &telegram.ChatMember{Status: `kicked`},
						NewChatMember: &telegram.ChatMember{Status: `administrator`},
					},
				},
				{
					Message: &telegram.Message{
						Text: "hi!",
						From: &telegram.User{Username: `mockUser`},
						Chat: &telegram.Chat{Type: `private`, ID: 667, Username: `mockUser`},
					},
				},
			},
		}))

		go func() {
			defer ginkgo.GinkgoRecover()
			conf, err := gen.Generate(nil, nil, nil)

			gomega.Expect(conf).ToNot(gomega.BeNil())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			resultChannel <- conf.GetURL().String()
		}()

		defer dumpBuffers()

		mockTyped(mockToken)
		mockTyped(`no`)

		gomega.Eventually(userIn).Should(gbytes.Say(`Got a bot chat member update for mockChannel, status was changed from kicked to administrator`))
		gomega.Eventually(userIn).Should(gbytes.Say(`Got 1 chat ID\(s\) so far\. Want to add some more\?`))
		gomega.Eventually(userIn).Should(gbytes.Say(`Selected chats:`))
		gomega.Eventually(userIn).Should(gbytes.Say(`667 \(private\) @mockUser`))

		gomega.Eventually(resultChannel).Should(gomega.Receive(gomega.Equal(`telegram://0:MockToken@telegram?chats=667&preview=No`)))
	})
})
