package services_test

import (
	"log"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/nicholas-fedor/shoutrrr/internal/testutils"
	"github.com/nicholas-fedor/shoutrrr/pkg/router"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
)

func TestServices(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Service Compliance Suite")
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
	"teams":      "teams://11111111-4444-4444-8444-cccccccccccc@22222222-4444-4444-8444-cccccccccccc/33333333012222222222333333333344/44444444-4444-4444-8444-cccccccccccc/V2ESyij_gAljSoUQHvZoZYzlpAoAXExyOl26dlf1xHEx05?host=test.webhook.office.com",
	"telegram":   "telegram://000000000:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA@telegram?channels=channel",
	"xmpp":       "xmpp://",
	"zulip":      "zulip://mail:key@example.com/?stream=foo&topic=bar",
}

var serviceResponses = map[string]string{
	"discord":    "",
	"gotify":     `{"id": 0}`,
	"googlechat": "",
	"hangouts":   "",
	"ifttt":      "",
	"join":       "",
	"logger":     "",
	"mattermost": "",
	"opsgenie":   "",
	"pushbullet": `{"type": "note", "body": "test", "title": "test title", "active": true, "created": 0}`,
	"pushover":   "",
	"rocketchat": "",
	"slack":      "",
	"smtp":       "",
	"teams":      "",
	"telegram":   "",
	"xmpp":       "",
	"zulip":      "",
}

var logger = log.New(ginkgo.GinkgoWriter, "Test", log.LstdFlags)

var _ = ginkgo.Describe("services", func() {
	ginkgo.BeforeEach(func() {
	})
	ginkgo.AfterEach(func() {
	})

	ginkgo.When("passed the a title param", func() {
		var serviceRouter *router.ServiceRouter

		ginkgo.AfterEach(func() {
			httpmock.DeactivateAndReset()
		})

		for key, configURL := range serviceURLs {
			serviceRouter, _ = router.New(logger)

			ginkgo.It("should not throw an error for "+key, func() {
				if key == "smtp" {
					ginkgo.Skip("smtp does not use HTTP and needs a specific test")
				}
				if key == "xmpp" {
					ginkgo.Skip("not supported")
				}

				service, err := serviceRouter.Locate(configURL)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				httpmock.Activate()
				if mockService, ok := service.(testutils.MockClientService); ok {
					httpmock.ActivateNonDefault(mockService.GetHTTPClient())
				}

				respStatus := http.StatusOK
				if key == "discord" || key == "ifttt" {
					respStatus = http.StatusNoContent
				}
				if key == "mattermost" {
					httpmock.RegisterResponder(
						"POST",
						"https://example.com/hooks/token",
						httpmock.NewStringResponder(http.StatusOK, ""),
					)
				} else {
					httpmock.RegisterNoResponder(httpmock.NewStringResponder(respStatus, serviceResponses[key]))
				}

				err = service.Send("test", (*types.Params)(&map[string]string{
					"title": "test title",
				}))
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})

			if key == "mattermost" {
				ginkgo.It("should not throw an error for "+key+" with DisableTLS", func() {
					modifiedURL := configURL + "?disabletls=yes"
					service, err := serviceRouter.Locate(modifiedURL)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					httpmock.Activate()
					if mockService, ok := service.(testutils.MockClientService); ok {
						httpmock.ActivateNonDefault(mockService.GetHTTPClient())
					}
					httpmock.RegisterResponder(
						"POST",
						"http://example.com/hooks/token",
						httpmock.NewStringResponder(http.StatusOK, ""),
					)

					err = service.Send("test", (*types.Params)(&map[string]string{
						"title": "test title",
					}))
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
				})
			}
		}
	})
})
