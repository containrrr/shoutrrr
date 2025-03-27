package rocketchat

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/nicholas-fedor/shoutrrr/internal/testutils"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var (
	service          *Service
	envRocketchatURL *url.URL
	_                = ginkgo.BeforeSuite(func() {
		service = &Service{}
		envRocketchatURL, _ = url.Parse(os.Getenv("SHOUTRRR_ROCKETCHAT_URL"))
	})
)

// Constants for repeated test values.
const (
	testTokenA = "tokenA"
	testTokenB = "tokenB"
)

func TestRocketchat(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Shoutrrr Rocketchat Suite")
}

var _ = ginkgo.Describe("the rocketchat service", func() {
	// Add tests for Initialize()
	ginkgo.Describe("Initialize method", func() {
		ginkgo.When("initializing with a valid URL", func() {
			ginkgo.It("should set logger and config without error", func() {
				service := &Service{}
				testURL, _ := url.Parse("rocketchat://testUser@rocketchat.my-domain.com:5055/" + testTokenA + "/" + testTokenB + "/#testChannel")
				err := service.Initialize(testURL, testutils.TestLogger())
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(service.Config).NotTo(gomega.BeNil())
				gomega.Expect(service.Config.Host).To(gomega.Equal("rocketchat.my-domain.com"))
				gomega.Expect(service.Config.Port).To(gomega.Equal("5055"))
				gomega.Expect(service.Config.UserName).To(gomega.Equal("testUser"))
				gomega.Expect(service.Config.TokenA).To(gomega.Equal(testTokenA))
				gomega.Expect(service.Config.TokenB).To(gomega.Equal(testTokenB))
				gomega.Expect(service.Config.Channel).To(gomega.Equal("#testChannel"))
			})
		})
		ginkgo.When("initializing with an invalid URL", func() {
			ginkgo.It("should return an error", func() {
				service := &Service{}
				testURL, _ := url.Parse("rocketchat://rocketchat.my-domain.com") // Missing tokens
				err := service.Initialize(testURL, testutils.TestLogger())
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(err.Error()).To(gomega.Equal(NotEnoughArguments))
			})
		})
	})

	// Add tests for Send()
	ginkgo.Describe("Send method", func() {
		var (
			mockServer *httptest.Server
			service    *Service
			client     *http.Client
		)

		ginkgo.BeforeEach(func() {
			// Create TLS server
			mockServer = httptest.NewTLSServer(nil) // Handler set in each test

			// Configure client to trust the mock server's certificate
			certPool := x509.NewCertPool()
			for _, cert := range mockServer.TLS.Certificates {
				certPool.AddCert(cert.Leaf)
			}
			client = &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						RootCAs:    certPool,
						MinVersion: tls.VersionTLS12, // Explicitly set minimum TLS version to 1.2
					},
				},
			}

			service = &Service{
				Config: &Config{},
				Client: client, // Assign the custom client here
			}
			service.Logger.SetLogger(testutils.TestLogger())
		})

		ginkgo.AfterEach(func() {
			if mockServer != nil {
				mockServer.Close()
			}
		})

		ginkgo.When("sending a message to a mock server with success", func() {
			ginkgo.It("should return no error", func() {
				mockServer.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				})
				mockURL, _ := url.Parse(mockServer.URL)
				service.Config.Host = mockURL.Hostname()
				service.Config.Port = mockURL.Port()
				service.Config.TokenA = testTokenA
				service.Config.TokenB = testTokenB

				err := service.Send("test message", nil)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
		})

		ginkgo.When("sending a message to a mock server with failure", func() {
			ginkgo.It("should return an error with status code and body", func() {
				mockServer.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("bad request"))
				})
				mockURL, _ := url.Parse(mockServer.URL)
				service.Config.Host = mockURL.Hostname()
				service.Config.Port = mockURL.Port()
				service.Config.TokenA = testTokenA
				service.Config.TokenB = testTokenB

				err := service.Send("test message", nil)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(err.Error()).To(gomega.ContainSubstring("notification failed: 400 bad request"))
			})
		})

		ginkgo.When("sending a message to an unreachable server", func() {
			ginkgo.It("should return a connection error", func() {
				service.Client = http.DefaultClient // Reset to default client for this test
				service.Config.Host = "nonexistent.domain"
				service.Config.TokenA = testTokenA
				service.Config.TokenB = testTokenB

				err := service.Send("test message", nil)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(err.Error()).To(gomega.ContainSubstring("error while posting to URL"))
			})
		})

		ginkgo.When("sending a message with params overriding username and channel", func() {
			ginkgo.It("should use params values in the payload", func() {
				mockServer.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				})
				mockURL, _ := url.Parse(mockServer.URL)
				service.Config.Host = mockURL.Hostname()
				service.Config.Port = mockURL.Port()
				service.Config.TokenA = testTokenA
				service.Config.TokenB = testTokenB
				service.Config.UserName = "defaultUser"
				service.Config.Channel = "#defaultChannel"

				params := types.Params{
					"username": "overrideUser",
					"channel":  "#overrideChannel",
				}
				err := service.Send("test message", &params)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				// Note: We can't directly inspect the payload here without mocking CreateJSONPayload,
				// but this ensures the params path is exercised.
			})
		})
	})

	// Add tests for GetURL() and SetURL()
	ginkgo.Describe("the rocketchat config", func() {
		ginkgo.When("generating a URL from a config with all fields", func() {
			ginkgo.It("should construct a correct URL", func() {
				config := &Config{
					Host:   "rocketchat.my-domain.com",
					Port:   "5055",
					TokenA: testTokenA,
					TokenB: testTokenB,
				}
				url := config.GetURL()
				gomega.Expect(url.String()).To(gomega.Equal("rocketchat://rocketchat.my-domain.com:5055/" + testTokenA + "/" + testTokenB))
			})
		})

		ginkgo.When("generating a URL from a config without port", func() {
			ginkgo.It("should construct a correct URL without port", func() {
				config := &Config{
					Host:   "rocketchat.my-domain.com",
					TokenA: testTokenA,
					TokenB: testTokenB,
				}
				url := config.GetURL()
				gomega.Expect(url.String()).To(gomega.Equal("rocketchat://rocketchat.my-domain.com/" + testTokenA + "/" + testTokenB))
			})
		})

		ginkgo.When("setting URL with a channel starting with @", func() {
			ginkgo.It("should set channel without adding #", func() {
				config := &Config{}
				testURL, _ := url.Parse("rocketchat://rocketchat.my-domain.com/" + testTokenA + "/" + testTokenB + "/@user")
				err := config.SetURL(testURL)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(config.Channel).To(gomega.Equal("@user"))
			})
		})

		ginkgo.When("setting URL with a regular channel without fragment", func() {
			ginkgo.It("should prepend # to the channel", func() {
				config := &Config{}
				testURL, _ := url.Parse("rocketchat://rocketchat.my-domain.com/" + testTokenA + "/" + testTokenB + "/general")
				err := config.SetURL(testURL)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(config.Channel).To(gomega.Equal("#general"))
			})
		})
	})

	// Add test for GetID()
	ginkgo.Describe("GetID method", func() {
		ginkgo.It("should return the correct scheme", func() {
			service := &Service{}
			id := service.GetID()
			gomega.Expect(id).To(gomega.Equal(Scheme))
		})
	})
})
