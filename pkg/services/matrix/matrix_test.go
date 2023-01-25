package matrix

import (
	"fmt"
	"net/url"

	"github.com/containrrr/shoutrrr/internal/testutils"
	"github.com/jarcoal/httpmock"

	"log"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMatrix(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Matrix Suite")
}

var _ = Describe("the matrix service", func() {
	var service *Service
	logger := log.New(GinkgoWriter, "Test", log.LstdFlags)
	envMatrixURL := os.Getenv("SHOUTRRR_MATRIX_URL")

	BeforeEach(func() {
		service = &Service{}
	})

	When("running integration tests", func() {
		It("should not error out", func() {
			if envMatrixURL == "" {
				return
			}
			serviceURL, err := url.Parse(envMatrixURL)
			Expect(err).NotTo(HaveOccurred())
			err = service.Initialize(serviceURL, logger)
			Expect(err).NotTo(HaveOccurred())
			err = service.Send("This is an integration test message", nil)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("creating configurations", func() {
		When("given an url with title prop", func() {
			It("should not throw an error", func() {
				serviceURL := testutils.URLMust(`matrix://user:pass@mockserver?rooms=room1&title=Better%20Off%20Alone`)
				Expect((&Config{}).SetURL(serviceURL)).To(Succeed())
			})
		})

		When("given an url with the prop `room`", func() {
			It("should treat is as an alias for `rooms`", func() {
				serviceURL := testutils.URLMust(`matrix://user:pass@mockserver?room=room1`)
				config := Config{}
				Expect(config.SetURL(serviceURL)).To(Succeed())
				Expect(config.Rooms).To(ContainElement("#room1"))
			})
		})
		When("given an url with invalid props", func() {
			It("should return an error", func() {
				serviceURL := testutils.URLMust(`matrix://user:pass@mockserver?channels=room1,room2`)
				Expect((&Config{}).SetURL(serviceURL)).To(HaveOccurred())
			})
		})
		When("parsing the configuration URL", func() {
			It("should be identical after de-/serialization", func() {
				testURL := "matrix://user:pass@mockserver?rooms=%23room1%2C%23room2"

				url, err := url.Parse(testURL)
				Expect(err).NotTo(HaveOccurred(), "parsing")

				config := &Config{}
				err = config.SetURL(url)
				Expect(err).NotTo(HaveOccurred(), "verifying")

				outputURL := config.GetURL()

				Expect(outputURL.String()).To(Equal(testURL))

			})
		})
	})

	Describe("the matrix client", func() {

		BeforeEach(func() {
			httpmock.Activate()
		})

		When("not providing a logger", func() {
			It("should not crash", func() {
				setupMockResponders()
				serviceURL := testutils.URLMust("matrix://user:pass@mockserver")
				Expect(service.Initialize(serviceURL, nil)).To(Succeed())
			})
		})

		When("sending a message", func() {
			It("should not report any errors", func() {
				setupMockResponders()
				serviceURL, _ := url.Parse("matrix://user:pass@mockserver")
				err := service.Initialize(serviceURL, logger)
				Expect(err).NotTo(HaveOccurred())

				err = service.Send("Test message", nil)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("sending a message to explicit rooms", func() {
			It("should not report any errors", func() {
				setupMockResponders()
				serviceURL, _ := url.Parse("matrix://user:pass@mockserver?rooms=room1,room2")
				err := service.Initialize(serviceURL, logger)
				Expect(err).NotTo(HaveOccurred())

				err = service.Send("Test message", nil)
				Expect(err).NotTo(HaveOccurred())
			})
			When("sending to one room fails", func() {
				It("should report one error", func() {
					setupMockResponders()
					serviceURL, _ := url.Parse("matrix://user:pass@mockserver?rooms=secret,room2")
					err := service.Initialize(serviceURL, logger)
					Expect(err).NotTo(HaveOccurred())

					err = service.Send("Test message", nil)
					Expect(err).To(HaveOccurred())

				})
			})

		})

		AfterEach(func() {
			httpmock.DeactivateAndReset()
		})
	})

	It("should implement basic service API methods correctly", func() {
		testutils.TestConfigGetInvalidQueryValue(&Config{})
		testutils.TestConfigSetInvalidQueryValue(&Config{}, "matrix://user:pass@host/?foo=bar")

		testutils.TestConfigGetEnumsCount(&Config{}, 0)
		testutils.TestConfigGetFieldsCount(&Config{}, 4)
	})
})

func setupMockResponders() {
	const mockServer = "https://mockserver"

	httpmock.RegisterResponder(
		"GET",
		mockServer+apiLogin,
		httpmock.NewStringResponder(200, `{"flows": [ { "type": "m.login.password" } ] }`))

	httpmock.RegisterResponder(
		"POST",
		mockServer+apiLogin,
		httpmock.NewStringResponder(200, `{ "access_token": "TOKEN", "home_server": "mockserver", "user_id": "test:mockerserver" }`))

	httpmock.RegisterResponder(
		"GET",
		mockServer+apiJoinedRooms,
		httpmock.NewStringResponder(200, `{ "joined_rooms": [ "!room:mockserver" ] }`))

	httpmock.RegisterResponder("POST", mockServer+fmt.Sprintf(apiSendMessage, "%21room:mockserver"),
		httpmock.NewJsonResponderOrPanic(200, apiResEvent{EventID: "7"}))

	httpmock.RegisterResponder("POST", mockServer+fmt.Sprintf(apiSendMessage, "1"),
		httpmock.NewJsonResponderOrPanic(200, apiResEvent{EventID: "8"}))

	httpmock.RegisterResponder("POST", mockServer+fmt.Sprintf(apiSendMessage, "2"),
		httpmock.NewJsonResponderOrPanic(200, apiResEvent{EventID: "9"}))

	httpmock.RegisterResponder("POST", mockServer+fmt.Sprintf(apiRoomJoin, "%23room1"),
		httpmock.NewJsonResponderOrPanic(200, apiResRoom{RoomID: "1"}))

	httpmock.RegisterResponder("POST", mockServer+fmt.Sprintf(apiRoomJoin, "%23room2"),
		httpmock.NewJsonResponderOrPanic(200, apiResRoom{RoomID: "2"}))

	httpmock.RegisterResponder("POST", mockServer+fmt.Sprintf(apiRoomJoin, "%23secret"),
		httpmock.NewJsonResponderOrPanic(403, apiResError{
			Code:    "M_FORBIDDEN",
			Message: "You are not invited to this room.",
		}))

}
