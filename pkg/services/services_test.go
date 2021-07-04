package services_test

import (
	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/containrrr/shoutrrr/pkg/services/gotify"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/jarcoal/httpmock"
	"log"
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestServices(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Compliance Suite")
}

var serviceURLs = map[string]string{
	"discord":    "discord://token@id",
	"gotify":     "gotify://example.com/Aaa.bbb.ccc.ddd",
	"googlechat": "googlechat://chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=baz",
	"hangouts":   "hangouts://chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=baz",
	"ifttt":      "ifttt://key?events=event",
	"join":       "join://:apikey@join/?devices=device",
	"logger":     "logger://",
	"mattermost": "mattermost://user@example.com/token",
	"opsgenie":   "opsgenie://example.com/token?responders=user:dummy",
	"pushbullet": "pushbullet://tokentokentokentokentokentokentoke",
	"pushover":   "pushover://:token@user/?devices=device",
	"rocketchat": "rocketchat://example.com/token/channel",
	"slack":      "slack://AAAAAAAAA/BBBBBBBBB/123456789123456789123456",
	"smtp":       "smtp://host.tld:25/?fromAddress=from@host.tld&toAddresses=to@host.tld",
	"teams":      "teams://11111111-4444-4444-8444-cccccccccccc@22222222-4444-4444-8444-cccccccccccc/33333333012222222222333333333344/44444444-4444-4444-8444-cccccccccccc",
	"telegram":   "telegram://000000000:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA@telegram?channels=channel",
	"xmpp":       "xmpp://",
	"zulip":      "zulip://mail:key@example.com/?stream=foo&topic=bar",
}

var serviceResponses = map[string]string{
	"pushbullet": `{"created": 0}`,
}

var logger = log.New(GinkgoWriter, "Test", log.LstdFlags)

var _ = Describe("services", func() {

	BeforeEach(func() {

	})
	AfterEach(func() {

	})

	When("passed the a title param", func() {

		var serviceRouter *router.ServiceRouter

		AfterEach(func() {
			httpmock.DeactivateAndReset()
		})

		for key, configURL := range serviceURLs {

			key := key //necessary to ensure the correct value is passed to the closure
			configURL := configURL
			serviceRouter, _ = router.New(logger)

			It("should not throw an error for "+key, func() {

				if key == "smtp" {
					Skip("smtp does not use HTTP and needs a specific test")
				}
				if key == "xmpp" {
					Skip("not supported")
				}

				httpmock.Activate()
				// Always return an "OK" result, as the http request isn't what is under test
				respStatus := http.StatusOK
				if key == "discord" || key == "ifttt" {
					respStatus = http.StatusNoContent
				}
				httpmock.RegisterNoResponder(httpmock.NewStringResponder(respStatus, serviceResponses[key]))

				service, err := serviceRouter.Locate(configURL)
				Expect(err).NotTo(HaveOccurred())

				if key == "gotify" {
					gotifyService := service.(*gotify.Service)
					httpmock.ActivateNonDefault(gotifyService.Client)
				}

				err = service.Send("test", (*types.Params)(&map[string]string{
					"title": "test title",
				}))
				Expect(err).NotTo(HaveOccurred())
			})

		}
	})

})
