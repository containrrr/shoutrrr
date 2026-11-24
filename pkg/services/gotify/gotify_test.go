package gotify

import (
	"log"
	"net/http"
	"net/url"
	"testing"

	"github.com/containrrr/shoutrrr/internal/testutils"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGotify(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Gotify Suite")
}

// Synthetic examples shared by validation, configuration, and delivery tests.
var supportedTokens = []struct {
	name  string
	token string
}{
	{"legacy", "AzyoeNS.D4iJLVa"},
	{"modern", "gtfya.ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmno-_"},
	{"signed modern", "gtfya.ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmno-_.1700000000.ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmno-_"},
}

var logger *log.Logger

var _ = Describe("the Gotify plugin URL building and token validation functions", func() {
	It("should build a valid gotify URL", func() {
		config := Config{
			Token: "Aaa.bbb.ccc.ddd",
			Host:  "my.gotify.tld",
		}
		url, err := buildURL(&config)
		Expect(err).To(BeNil())
		expectedURL := "https://my.gotify.tld/message?token=Aaa.bbb.ccc.ddd"
		Expect(url).To(Equal(expectedURL))
	})

	When("TLS is disabled", func() {
		It("should use http schema", func() {
			config := Config{
				Token:      "Aaa.bbb.ccc.ddd",
				Host:       "my.gotify.tld",
				DisableTLS: true,
			}
			url, err := buildURL(&config)
			Expect(err).To(BeNil())
			expectedURL := "http://my.gotify.tld/message?token=Aaa.bbb.ccc.ddd"
			Expect(url).To(Equal(expectedURL))
		})
	})

	When("a custom path is provided", func() {
		It("should add it to the URL", func() {
			config := Config{
				Token: "Aaa.bbb.ccc.ddd",
				Host:  "my.gotify.tld",
				Path:  "/gotify",
			}
			url, err := buildURL(&config)
			Expect(err).To(BeNil())
			expectedURL := "https://my.gotify.tld/gotify/message?token=Aaa.bbb.ccc.ddd"
			Expect(url).To(Equal(expectedURL))
		})
	})

	Describe("supported application token formats", func() {
		for _, tokenCase := range supportedTokens {
			It("should validate a "+tokenCase.name+" token", func() {
				Expect(isTokenValid(tokenCase.token)).To(BeTrue())
			})

			for _, pathCase := range []struct {
				name         string
				subpath      string
				endpointPath string
			}{
				{"without subpath", "", "/message"},
				// Preserve the existing URL builder's separator after a parsed subpath.
				{"with subpath", "/path/to/gotify", "/path/to/gotify//message"},
			} {
				for _, trailingSlash := range []string{"", "/"} {
					name := tokenCase.name + " " + pathCase.name + " trailing slash=" + trailingSlash
					canonicalURL := "gotify://my.gotify.tld" + pathCase.subpath + "/" + tokenCase.token + "?title=Test+title"
					inputURL := "gotify://my.gotify.tld" + pathCase.subpath + "/" + tokenCase.token + trailingSlash + "?title=Test+title"

					It("should preserve the token and canonicalize configuration for "+name, func() {
						config := &Config{}
						Expect(config.SetURL(testutils.URLMust(inputURL))).To(Succeed())
						Expect(config.Token).To(Equal(tokenCase.token))
						Expect(config.GetURL().String()).To(Equal(canonicalURL))

						roundTrip := &Config{}
						Expect(roundTrip.SetURL(config.GetURL())).To(Succeed())
						Expect(roundTrip.Token).To(Equal(tokenCase.token))
					})

					It("should initialize and deliver the complete token for "+name, func() {
						service := &Service{}
						Expect(service.Initialize(testutils.URLMust(inputURL), logger)).To(Succeed())
						httpmock.ActivateNonDefault(service.GetHTTPClient())
						DeferCleanup(httpmock.DeactivateAndReset)

						calls := 0
						httpmock.RegisterResponder("POST", "https://my.gotify.tld"+pathCase.endpointPath+"?token="+tokenCase.token,
							func(request *http.Request) (*http.Response, error) {
								calls++
								Expect(request.URL.Path).To(Equal(pathCase.endpointPath))
								Expect(request.URL.Query().Get("token")).To(Equal(tokenCase.token))
								Expect(request.URL.Query()).To(HaveLen(1))
								return testutils.JSONRespondMust(200, messageResponse{})(request)
							})

						Expect(service.Send("Message", nil)).To(Succeed())
						Expect(calls).To(Equal(1))
					})
				}
			}
		}
	})

	DescribeTable("checking modern token syntax without validating its key", func(token string) {
		Expect(isTokenValid(token)).To(BeTrue())
	},
		Entry("nonempty suffix without a fixed length", "gtfya.a"),
		Entry("all supported character classes", "gtfya.aAZ09.-_"),
	)

	DescribeTable("rejecting invalid application tokens", func(token string) {
		Expect(isTokenValid(token)).To(BeFalse())
		_, err := buildURL(&Config{Token: token, Host: "my.gotify.tld"})
		Expect(err).To(HaveOccurred())
	},
		Entry("empty", ""),
		Entry("bare modern prefix", "gtfya."),
		Entry("client prefix", "gtfyc.ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmno-_"),
		Entry("wrong prefix", "gtfyx.ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmno-_"),
		Entry("uppercase prefix", "GTFYA.ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmno-_"),
		Entry("short legacy", "AzyoeNS.D4iJLV"),
		Entry("long legacy", "AzyoeNS.D4iJLVaX"),
		Entry("legacy whitespace", "AzyoeNS D4iJLVa"),
		Entry("modern whitespace", "gtfya.abc def"),
		Entry("newline", "gtfya.abc\n"),
		Entry("Unicode", "gtfya.abcé"),
		Entry("query separator", "gtfya.abc&priority=1"),
		Entry("query start", "gtfya.abc?priority=1"),
		Entry("query assignment", "gtfya.abc=def"),
		Entry("fragment", "gtfya.abc#def"),
		Entry("path separator", "gtfya.abc/def"),
	)

	When("provided a valid token", func() {
		It("should return true", func() {
			token := "Ahwbsdyhwwgarxd"
			Expect(isTokenValid(token)).To(BeTrue())
		})
	})
	When("provided a token with an invalid prefix", func() {
		It("should return false", func() {
			token := "Chwbsdyhwwgarxd"
			Expect(isTokenValid(token)).To(BeFalse())
		})
	})
	When("provided a token with an invalid length", func() {
		It("should return false", func() {
			token := "Chwbsdyhwwga"
			Expect(isTokenValid(token)).To(BeFalse())
		})
	})
	Describe("creating the API URL", func() {
		When("the token is invalid", func() {
			It("should return an error", func() {
				config := Config{
					Token: "invalid",
				}
				_, err := buildURL(&config)
				Expect(err).To(HaveOccurred())
			})
		})
	})
	Describe("creating a config", func() {
		When("parsing the configuration URL", func() {
			It("should be identical after de-/serialization (with path)", func() {
				testURL := "gotify://my.gotify.tld/gotify/Aaa.bbb.ccc.ddd?title=Test+title"

				config := &Config{}
				Expect(config.SetURL(testutils.URLMust(testURL))).To(Succeed())
				Expect(config.GetURL().String()).To(Equal(testURL))
			})
			It("should be identical after de-/serialization (without path)", func() {
				testURL := "gotify://my.gotify.tld/Aaa.bbb.ccc.ddd?disabletls=Yes&priority=1&title=Test+title"

				config := &Config{}
				Expect(config.SetURL(testutils.URLMust(testURL))).To(Succeed())
				Expect(config.GetURL().String()).To(Equal(testURL))

			})
			It("should allow slash at the end of the token", func() {
				url := testutils.URLMust("gotify://my.gotify.tld/Aaa.bbb.ccc.ddd/")

				config := &Config{}
				Expect(config.SetURL(url)).To(Succeed())
				Expect(config.Token).To(Equal("Aaa.bbb.ccc.ddd"))
			})
			It("should allow slash at the end of the token, with additional path", func() {
				url := testutils.URLMust("gotify://my.gotify.tld/path/to/gotify/Aaa.bbb.ccc.ddd/")

				config := &Config{}
				Expect(config.SetURL(url)).To(Succeed())
				Expect(config.Token).To(Equal("Aaa.bbb.ccc.ddd"))
			})
			It("should not crash on empty token or path slash at the end of the token", func() {
				config := &Config{}
				Expect(config.SetURL(testutils.URLMust("gotify://my.gotify.tld//"))).To(Succeed())
				Expect(config.SetURL(testutils.URLMust("gotify://my.gotify.tld/"))).To(Succeed())
			})
		})
	})

	Describe("sending the payload", func() {
		var err error
		var service Service
		AfterEach(func() {
			httpmock.DeactivateAndReset()
		})
		It("should not report an error if the server accepts the payload", func() {
			serviceURL, _ := url.Parse("gotify://my.gotify.tld/Aaa.bbb.ccc.ddd")
			err = service.Initialize(serviceURL, logger)
			httpmock.ActivateNonDefault(service.GetHTTPClient())
			Expect(err).NotTo(HaveOccurred())

			targetURL := "https://my.gotify.tld/message?token=Aaa.bbb.ccc.ddd"
			httpmock.RegisterResponder("POST", targetURL, testutils.JSONRespondMust(200, messageResponse{}))

			err = service.Send("Message", nil)
			Expect(err).NotTo(HaveOccurred())
		})
		It("should not panic if an error occurs when sending the payload", func() {
			serviceURL, _ := url.Parse("gotify://my.gotify.tld/Aaa.bbb.ccc.ddd")
			err = service.Initialize(serviceURL, logger)
			httpmock.ActivateNonDefault(service.GetHTTPClient())
			Expect(err).NotTo(HaveOccurred())

			targetURL := "https://my.gotify.tld/message?token=Aaa.bbb.ccc.ddd"
			httpmock.RegisterResponder("POST", targetURL, testutils.JSONRespondMust(401, errorResponse{
				Name:        "Unauthorized",
				Code:        401,
				Description: "you need to provide a valid access token or user credentials to access this api",
			}))

			err = service.Send("Message", nil)
			Expect(err).To(HaveOccurred())
		})
	})
})
