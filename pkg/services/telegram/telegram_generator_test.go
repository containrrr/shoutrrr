package telegram_test

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/jarcoal/httpmock"
	"github.com/mattn/go-colorable"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"

	"github.com/containrrr/shoutrrr/pkg/services/telegram"
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

func mockTyped(a ...interface{}) {
	_, _ = fmt.Fprint(userOut, a...)
	_, _ = fmt.Fprint(userOut, "\n")
}

func dumpBuffers() {
	logger := log.New(GinkgoWriter, "Test", log.LstdFlags)
	for _, line := range strings.Split(string(userIn.Contents()), "\n") {
		logger.Println(">", line)
	}
	logger.Println("----")
	for _, line := range strings.Split(string(userOut.Contents()), "\n") {
		logger.Println("<", line)
	}
}

func mockAPI(endpoint string) string {
	return mockAPIBase + endpoint
}

var _ = Describe("TelegramGenerator", func() {

	BeforeEach(func() {
		userOut = gbytes.NewBuffer()
		userIn = gbytes.NewBuffer()
		userInMono = colorable.NewNonColorable(userIn)
		httpmock.Activate()
	})
	AfterEach(func() {
		httpmock.DeactivateAndReset()
	})
	It("should return the ", func() {
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
					Message: &telegram.Message{
						Text: "hi!",
						From: &telegram.User{Username: `mockUser`},
						Chat: &telegram.Chat{Type: `private`, ID: 667, Username: `mockUser`},
					},
				},
			},
		}))

		go func() {
			defer GinkgoRecover()
			conf, err := gen.Generate(nil, nil, nil)

			Expect(conf).ToNot(BeNil())
			Expect(err).NotTo(HaveOccurred())
			resultChannel <- conf.GetURL().String()
		}()

		defer dumpBuffers()

		mockTyped(mockToken)
		mockTyped(`no`)

		Eventually(userIn).Should(gbytes.Say(`Got 1 chat ID\(s\) so far\. Want to add some more\?`))
		Eventually(userIn).Should(gbytes.Say(`Selected chats:`))
		Eventually(userIn).Should(gbytes.Say(`667 \(private\) @mockUser`))

		Eventually(resultChannel).Should(Receive(Equal(`telegram://0:MockToken@telegram?chats=667`)))
	})

})
