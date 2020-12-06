package ifttt

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util"
	"github.com/jarcoal/httpmock"
)

func TestIFTTT(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr IFTTT Suite")
}

var (
	service    *Service
	logger     *log.Logger
	envTestURL string
)

var _ = Describe("the ifttt package", func() {
	BeforeSuite(func() {
		envTestURL = os.Getenv("SHOUTRRR_IFTTT_URL")
		logger = util.TestLogger()
	})
	BeforeEach(func() {
		service = &Service{}
	})
	When("running integration tests", func() {
		It("should work without errors", func() {
			if envTestURL == "" {
				return
			}

			serviceURL, err := url.Parse(envTestURL)
			Expect(err).NotTo(HaveOccurred())

			err = service.Initialize(serviceURL, logger)
			Expect(err).NotTo(HaveOccurred())

			err = service.Send(
				"this is an integration test",
				nil,
			)
			Expect(err).NotTo(HaveOccurred())
		})
	})
	When("creating a config", func() {
		When("given an url", func() {
			It("should return an error if no arguments where supplied", func() {
				serviceURL, _ := url.Parse("ifttt://")
				err := service.Initialize(serviceURL, logger)
				Expect(err).To(HaveOccurred())
			})
			It("should return an error if no webhook ID is given", func() {
				serviceURL, _ := url.Parse("ifttt:///?events=event1")
				err := service.Initialize(serviceURL, logger)
				Expect(err).To(HaveOccurred())
			})
			It("should return an error no events are given", func() {
				serviceURL, _ := url.Parse("ifttt://dummyID")
				err := service.Initialize(serviceURL, logger)
				Expect(err).To(HaveOccurred())
			})
			It("should return an error when an invalid query key is given", func() {
				serviceURL, _ := url.Parse("ifttt://dummyID/?events=event1&badquery=foo")
				err := service.Initialize(serviceURL, logger)
				Expect(err).To(HaveOccurred())
			})
			It("should return an error if message value is above 3", func() {
				serviceURL, _ := url.Parse("ifttt://dummyID/?events=event1&messagevalue=8")
				config := Config{}
				err := config.SetURL(serviceURL)
				Expect(err).To(HaveOccurred())
			})
			It("should not return an error if webhook ID and at least one event is given", func() {
				serviceURL, _ := url.Parse("ifttt://dummyID/?events=event1")
				err := service.Initialize(serviceURL, logger)
				Expect(err).NotTo(HaveOccurred())
			})
			It("should set value1, value2 and value3", func() {
				serviceURL, _ := url.Parse("ifttt://dummyID/?events=dummyevent&value3=three&value2=two&value1=one")
				config := Config{}
				err := config.SetURL(serviceURL)
				Expect(err).NotTo(HaveOccurred())

				Expect(config.Value1).To(Equal("one"))
				Expect(config.Value2).To(Equal("two"))
				Expect(config.Value3).To(Equal("three"))
			})
		})
	})
	When("serializing a config to URL", func() {
		When("given multiple events", func() {
			It("should return an URL with all the events comma-separated", func() {
				expectedURL := "ifttt://dummyID/?events=foo,bar,baz&messagevalue=0&value1=&value2=&value3="
				config := Config{
					Events:            []string{"foo", "bar", "baz"},
					WebHookID:         "dummyID",
					UseMessageAsValue: 0,
				}
				resultURL := config.GetURL().String()
				Expect(resultURL).To(Equal(expectedURL))
			})
		})
	})
	When("sending a message", func() {
		It("should error if the response code is not 204 no content", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			setupResponder("foo", "dummy", 404, "")

			URL, _ := url.Parse("ifttt://dummy/?events=foo")

			if err := service.Initialize(URL, logger); err != nil {
				Fail("errored during initialization")
			}

			err := service.Send("hello", nil)
			Expect(err).To(HaveOccurred())
		})
		It("should not error if the response code is 204", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			setupResponder("foo", "dummy", 204, "")

			URL, _ := url.Parse("ifttt://dummy/?events=foo")

			if err := service.Initialize(URL, logger); err != nil {
				Fail("errored during initialization")
			}

			err := service.Send("hello", nil)
			Expect(err).NotTo(HaveOccurred())
		})
	})
	When("creating a json payload", func() {
		When("given config values \"a\", \"b\" and \"c\"", func() {
			It("should return a valid jsonPayload string with values \"a\", \"b\" and \"c\"", func() {
				bytes, err := createJSONToSend(&Config{
					Value1:            "a",
					Value2:            "b",
					Value3:            "c",
					UseMessageAsValue: 0,
				}, "d", nil)
				Expect(err).ToNot(HaveOccurred())

				payload := jsonPayload{}
				err = json.Unmarshal(bytes, &payload)
				Expect(err).ToNot(HaveOccurred())

				Expect(payload.Value1).To(Equal("a"))
				Expect(payload.Value2).To(Equal("b"))
				Expect(payload.Value3).To(Equal("c"))
			})
		})
		When("message value is set to 3", func() {
			It("should return a jsonPayload string with value2 set to message", func() {
				config := &Config{
					Value1: "a",
					Value2: "b",
					Value3: "c",
				}

				for i := 1; i <= 3; i++ {
					config.UseMessageAsValue = uint8(i)
					bytes, err := createJSONToSend(config, "d", nil)
					Expect(err).ToNot(HaveOccurred())

					payload := jsonPayload{}
					err = json.Unmarshal(bytes, &payload)
					Expect(err).ToNot(HaveOccurred())

					if i == 1 {
						Expect(payload.Value1).To(Equal("d"))
					} else if i == 2 {
						Expect(payload.Value2).To(Equal("d"))
					} else if i == 3 {
						Expect(payload.Value3).To(Equal("d"))
					}

				}
			})
		})
		When("given a param overrides for value1, value2 and value3", func() {
			It("should return a jsonPayload string with value1, value2 and value3 overridden", func() {
				bytes, err := createJSONToSend(&Config{
					Value1:            "a",
					Value2:            "b",
					Value3:            "c",
					UseMessageAsValue: 0,
				}, "d", (*types.Params)(&map[string]string{
					"value1": "e",
					"value2": "f",
					"value3": "g",
				}))
				Expect(err).ToNot(HaveOccurred())

				payload := &jsonPayload{}
				err = json.Unmarshal(bytes, payload)
				Expect(err).ToNot(HaveOccurred())

				Expect(payload.Value1).To(Equal("e"))
				Expect(payload.Value2).To(Equal("f"))
				Expect(payload.Value3).To(Equal("g"))
			})
		})
	})
})

func setupResponder(event string, key string, code int, body string) {
	url := fmt.Sprintf("https://maker.ifttt.com/trigger/%s/with/key/%s", event, key)
	httpmock.RegisterResponder("POST", url, httpmock.NewStringResponder(code, body))
}
