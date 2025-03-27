package smtp

import (
	"log"
	"net/smtp"
	"net/url"
	"os"
	"reflect"
	"testing"
	"unsafe"

	"github.com/nicholas-fedor/shoutrrr/internal/failures"
	"github.com/nicholas-fedor/shoutrrr/internal/testutils"
	"github.com/nicholas-fedor/shoutrrr/pkg/format"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/standard"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	gomegaTypes "github.com/onsi/gomega/types"
)

var tt *testing.T

func TestSMTP(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)

	tt = t
	ginkgo.RunSpecs(t, "Shoutrrr SMTP Suite")
}

var (
	service    *Service
	envSMTPURL string
	logger     *log.Logger
	_          = ginkgo.BeforeSuite(func() {
		envSMTPURL = os.Getenv("SHOUTRRR_SMTP_URL")
		logger = testutils.TestLogger()
	})
	urlWithAllProps = "smtp://user:password@example.com:2225/?auth=None&clienthost=testhost&encryption=ExplicitTLS&fromaddress=sender%40example.com&fromname=Sender&subject=Subject&toaddresses=rec1%40example.com%2Crec2%40example.com&usehtml=Yes&usestarttls=No"
	// BaseNoAuthURL is a minimal SMTP config without authentication.
	BaseNoAuthURL = "smtp://example.com:2225/?useStartTLS=no&auth=none&fromAddress=sender@example.com&toAddresses=rec1@example.com&useHTML=no"
	// BaseAuthURL is a typical config with authentication.
	BaseAuthURL = "smtp://user:password@example.com:2225/?useStartTLS=no&fromAddress=sender@example.com&toAddresses=rec1@example.com,rec2@example.com&useHTML=yes"
	// BasePlusURL is a config with plus signs in email addresses.
	BasePlusURL = "smtp://user:password@example.com:2225/?useStartTLS=no&fromAddress=sender+tag@example.com&toAddresses=rec1+tag@example.com,rec2@example.com&useHTML=yes"
)

// modifyURL modifies a base URL by updating query parameters as specified.
func modifyURL(base string, params map[string]string) string {
	u := testutils.URLMust(base)

	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}

	u.RawQuery = q.Encode()

	return u.String()
}

var _ = ginkgo.Describe("the SMTP service", func() {
	ginkgo.BeforeEach(func() {
		service = &Service{}
	})
	ginkgo.When("parsing the configuration URL", func() {
		ginkgo.It("should be identical after de-/serialization", func() {
			url := testutils.URLMust(urlWithAllProps)
			config := &Config{}
			pkr := format.NewPropKeyResolver(config)
			err := config.setURL(&pkr, url)
			gomega.Expect(err).NotTo(gomega.HaveOccurred(), "verifying")

			outputURL := config.GetURL()
			ginkgo.GinkgoT().Logf("\n\n%s\n%s\n\n-", outputURL, urlWithAllProps)

			gomega.Expect(outputURL.String()).To(gomega.Equal(urlWithAllProps))
		})
		ginkgo.When("resolving client host", func() {
			ginkgo.When("clienthost is set to auto", func() {
				ginkgo.It("should return the os hostname", func() {
					hostname, err := os.Hostname()
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(service.resolveClientHost(&Config{ClientHost: "auto"})).To(gomega.Equal(hostname))
				})
			})
			ginkgo.When("clienthost is set to a custom value", func() {
				ginkgo.It("should return that value", func() {
					gomega.Expect(service.resolveClientHost(&Config{ClientHost: "computah"})).To(gomega.Equal("computah"))
				})
			})
		})
		ginkgo.When("fromAddress is missing", func() {
			ginkgo.It("should return an error", func() {
				testURL := testutils.URLMust("smtp://user:password@example.com:2225/?toAddresses=rec1@example.com,rec2@example.com")
				gomega.Expect((&Config{}).SetURL(testURL)).ToNot(gomega.Succeed())
			})
		})
		ginkgo.When("toAddresses are missing", func() {
			ginkgo.It("should return an error", func() {
				testURL := testutils.URLMust("smtp://user:password@example.com:2225/?fromAddress=sender@example.com")
				gomega.Expect((&Config{}).SetURL(testURL)).NotTo(gomega.Succeed())
			})
		})
	})
	ginkgo.Context("basic service API methods", func() {
		var config *Config
		ginkgo.BeforeEach(func() {
			config = &Config{}
		})
		ginkgo.It("should not allow getting invalid query values", func() {
			testutils.TestConfigGetInvalidQueryValue(config)
		})
		ginkgo.It("should not allow setting invalid query values", func() {
			testutils.TestConfigSetInvalidQueryValue(
				config,
				"smtp://example.com/?fromAddress=s@example.com&toAddresses=r@example.com&foo=bar",
			)
		})

		ginkgo.It("should have the expected number of fields and enums", func() {
			testutils.TestConfigGetEnumsCount(config, 2)
			testutils.TestConfigGetFieldsCount(config, 13)
		})
	})
	ginkgo.When("cloning a config", func() {
		ginkgo.It("should be identical to the original", func() {
			config := &Config{}
			gomega.Expect(config.SetURL(testutils.URLMust(urlWithAllProps))).To(gomega.Succeed())

			gomega.Expect(config.Clone()).To(gomega.Equal(*config))
		})
	})
	ginkgo.When("sending a message", func() {
		ginkgo.When("the service is not configured correctly", func() {
			ginkgo.It("should fail to send messages", func() {
				service := Service{Config: &Config{}}
				gomega.Expect(service.Send("test message", nil)).To(matchFailure(FailGetSMTPClient))

				service.Config.Encryption = EncMethods.ImplicitTLS
				gomega.Expect(service.Send("test message", nil)).To(matchFailure(FailGetSMTPClient))
			})
		})
		ginkgo.When("an invalid param is passed", func() {
			ginkgo.It("should fail to send messages", func() {
				service := Service{Config: &Config{}}
				gomega.Expect(service.Send("test message", &types.Params{"invalid": "value"})).To(matchFailure(FailApplySendParams))
			})
		})
	})

	ginkgo.When("the underlying stream stops working", func() {
		var service Service
		var message string
		ginkgo.BeforeEach(func() {
			service = Service{}
			message = ""
		})
		ginkgo.It("should fail when writing multipart plain header", func() {
			writer := testutils.CreateFailWriter(1)
			err := service.writeMultipartMessage(writer, message)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(err.ID()).To(gomega.Equal(FailPlainHeader))
		})

		ginkgo.It("should fail when writing multipart plain message", func() {
			writer := testutils.CreateFailWriter(2)
			err := service.writeMultipartMessage(writer, message)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(err.ID()).To(gomega.Equal(FailMessageRaw))
		})

		ginkgo.It("should fail when writing multipart HTML header", func() {
			writer := testutils.CreateFailWriter(4)
			err := service.writeMultipartMessage(writer, message)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(err.ID()).To(gomega.Equal(FailHTMLHeader))
		})

		ginkgo.It("should fail when writing multipart HTML message", func() {
			writer := testutils.CreateFailWriter(5)
			err := service.writeMultipartMessage(writer, message)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(err.ID()).To(gomega.Equal(FailMessageRaw))
		})

		ginkgo.It("should fail when writing multipart end header", func() {
			writer := testutils.CreateFailWriter(6)
			err := service.writeMultipartMessage(writer, message)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(err.ID()).To(gomega.Equal(FailMultiEndHeader))
		})

		ginkgo.It("should fail when writing message template", func() {
			writer := testutils.CreateFailWriter(0)
			e := service.SetTemplateString("dummy", "dummy template content")
			gomega.Expect(e).ToNot(gomega.HaveOccurred())

			err := service.writeMessagePart(writer, message, "dummy")
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(err.ID()).To(gomega.Equal(FailMessageTemplate))
		})
	})

	ginkgo.When("running E2E tests", func() {
		ginkgo.It("should work without errors", func() {
			if envSMTPURL == "" {
				ginkgo.Skip("environment not set up for E2E testing")

				return
			}

			serviceURL, err := url.Parse(envSMTPURL)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			err = service.Initialize(serviceURL, logger)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			err = service.Send("this is an integration test", nil)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
	})

	ginkgo.When("running integration tests", func() {
		ginkgo.When("given a typical usage case configuration URL", func() {
			ginkgo.It("should send notifications without any errors", func() {
				testURL := BaseAuthURL
				err := testIntegration(testURL, []string{
					"250-mx.google.com at your service",
					"250-SIZE 35651584",
					"250-AUTH LOGIN PLAIN",
					"250 8BITMIME",
					"235 Accepted",
					"250 Sender OK",
					"250 Receiver OK",
					"354 Go ahead",
					"250 Data OK",
					"250 Sender OK",
					"250 Receiver OK",
					"354 Go ahead",
					"250 Data OK",
					"221 OK",
				}, "<pre>{{ .message }}</pre>", "{{ .message }}")
				if msg, test := standard.IsTestSetupFailure(err); test {
					ginkgo.Skip(msg)

					return
				}
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
		})

		ginkgo.When("given e-mail addresses with pluses in the configuration URL", func() {
			ginkgo.It("should send notifications without any errors", func() {
				testURL := BasePlusURL
				err := testIntegration(
					testURL,
					[]string{
						"250-mx.google.com at your service",
						"250-SIZE 35651584",
						"250-AUTH LOGIN PLAIN",
						"250 8BITMIME",
						"235 Accepted",
						"250 Sender OK",
						"250 Receiver OK",
						"354 Go ahead",
						"250 Data OK",
						"250 Sender OK",
						"250 Receiver OK",
						"354 Go ahead",
						"250 Data OK",
						"221 OK",
					},
					"<pre>{{ .message }}</pre>", "{{ .message }}",
					"RCPT TO:<rec1+tag@example.com>",
					"To: rec1+tag@example.com",
					"From:  <sender+tag@example.com>")
				if msg, test := standard.IsTestSetupFailure(err); test {
					ginkgo.Skip(msg)

					return
				}
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
		})

		ginkgo.When("given a configuration URL with authentication disabled", func() {
			ginkgo.It("should send notifications without any errors", func() {
				testURL := BaseNoAuthURL
				err := testIntegration(testURL, []string{
					"250-mx.google.com at your service",
					"250-SIZE 35651584",
					"250-AUTH LOGIN PLAIN",
					"250 8BITMIME",
					"250 Sender OK",
					"250 Receiver OK",
					"354 Go ahead",
					"250 Data OK",
					"221 OK",
				}, "", "")
				if msg, test := standard.IsTestSetupFailure(err); test {
					ginkgo.Skip(msg)

					return
				}
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
		})

		ginkgo.When("given a configuration URL with StartTLS but it is not supported", func() {
			ginkgo.It("should send notifications without any errors", func() {
				testURL := modifyURL(BaseNoAuthURL, map[string]string{"useStartTLS": "yes"})
				err := testIntegration(testURL, []string{
					"250-mx.google.com at your service",
					"250-SIZE 35651584",
					"250-AUTH LOGIN PLAIN",
					"250 8BITMIME",
					"250 Sender OK",
					"250 Receiver OK",
					"354 Go ahead",
					"250 Data OK",
					"221 OK",
				}, "", "")
				if msg, test := standard.IsTestSetupFailure(err); test {
					ginkgo.Skip(msg)

					return
				}
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
		})

		ginkgo.When("server communication fails", func() {
			ginkgo.It("should fail when initial handshake is not accepted", func() {
				testURL := modifyURL(BaseNoAuthURL, map[string]string{"useStartTLS": "yes", "clienthost": "spammer"})
				err := testIntegration(testURL, []string{
					"421 4.7.0 Try again later, closing connection. (EHLO) r20-20020a50d694000000b004588af8956dsm771862edi.9 - gsmtp",
				}, "", "")
				if msg, test := standard.IsTestSetupFailure(err); test {
					ginkgo.Skip(msg)

					return
				}
				gomega.Expect(err).To(gomega.MatchError(fail(FailHandshake, nil)))
			})

			ginkgo.It("should fail when not being able to enable StartTLS", func() {
				testURL := modifyURL(BaseNoAuthURL, map[string]string{"useStartTLS": "yes"})
				err := testIntegration(testURL, []string{
					"250-mx.google.com at your service",
					"250-SIZE 35651584",
					"250-STARTTLS",
					"250-AUTH LOGIN PLAIN",
					"250 8BITMIME",
					"502 That's too hard",
				}, "", "")
				if msg, test := standard.IsTestSetupFailure(err); test {
					ginkgo.Skip(msg)

					return
				}
				gomega.Expect(err).To(matchFailure(FailEnableStartTLS))
			})

			ginkgo.It("should fail when authentication type is invalid", func() {
				testURL := modifyURL(BaseNoAuthURL, map[string]string{"auth": "bad"})
				err := testIntegration(testURL, []string{}, "", "")
				if msg, test := standard.IsTestSetupFailure(err); test {
					ginkgo.Skip(msg)

					return
				}
				gomega.Expect(err).To(matchFailure(standard.FailServiceInit))
			})

			ginkgo.It("should fail when not being able to use authentication type", func() {
				testURL := modifyURL(BaseNoAuthURL, map[string]string{"auth": "crammd5"})
				err := testIntegration(testURL, []string{
					"250-mx.google.com at your service",
					"250-SIZE 35651584",
					"250-AUTH LOGIN PLAIN",
					"250 8BITMIME",
					"504 Liar",
				}, "", "")
				if msg, test := standard.IsTestSetupFailure(err); test {
					ginkgo.Skip(msg)

					return
				}
				gomega.Expect(err).To(matchFailure(FailAuthenticating))
			})

			ginkgo.It("should fail when not being able to send to recipient", func() {
				testURL := BaseNoAuthURL
				err := testIntegration(testURL, []string{
					"250-mx.google.com at your service",
					"250-SIZE 35651584",
					"250-AUTH LOGIN PLAIN",
					"250 8BITMIME",
					"551 I don't know you",
				}, "", "")
				if msg, test := standard.IsTestSetupFailure(err); test {
					ginkgo.Skip(msg)

					return
				}
				gomega.Expect(err).To(matchFailure(FailSendRecipient))
			})

			ginkgo.It("should fail when the recipient is not accepted", func() {
				testURL := BaseNoAuthURL
				err := testSendRecipient(testURL, []string{
					"250 mx.google.com at your service",
					"250 Sender OK",
					"553 She doesn't want to be disturbed",
				})
				if msg, test := standard.IsTestSetupFailure(err); test {
					ginkgo.Skip(msg)

					return
				}
				gomega.Expect(err).To(matchFailure(FailSetRecipient))
			})

			ginkgo.It("should fail when the server does not accept the data stream", func() {
				testURL := BaseNoAuthURL
				err := testSendRecipient(testURL, []string{
					"250 mx.google.com at your service",
					"250 Sender OK",
					"250 Receiver OK",
					"554 Nah I'm fine thanks",
				})
				if msg, test := standard.IsTestSetupFailure(err); test {
					ginkgo.Skip(msg)

					return
				}
				gomega.Expect(err).To(matchFailure(FailOpenDataStream))
			})

			ginkgo.It("should fail when the server does not accept the data stream content", func() {
				testURL := BaseNoAuthURL
				err := testSendRecipient(testURL, []string{
					"250 mx.google.com at your service",
					"250 Sender OK",
					"250 Receiver OK",
					"354 Go ahead",
					"554 Such garbage!",
				})
				if msg, test := standard.IsTestSetupFailure(err); test {
					ginkgo.Skip(msg)

					return
				}
				gomega.Expect(err).To(matchFailure(FailCloseDataStream))
			})

			ginkgo.It("should fail when the server does not close the connection gracefully", func() {
				testURL := BaseNoAuthURL
				err := testIntegration(testURL, []string{
					"250-mx.google.com at your service",
					"250-SIZE 35651584",
					"250-AUTH LOGIN PLAIN",
					"250 8BITMIME",
					"250 Sender OK",
					"250 Receiver OK",
					"354 Go ahead",
					"250 Data OK",
					"502 You can't quit, you're fired!",
				}, "", "")
				if msg, test := standard.IsTestSetupFailure(err); test {
					ginkgo.Skip(msg)

					return
				}
				gomega.Expect(err).To(matchFailure(FailClosingSession))
			})
		})
	})
	ginkgo.When("writing headers and the output stream is closed", func() {
		ginkgo.When("it's closed during header content", func() {
			ginkgo.It("should fail with correct error", func() {
				fw := testutils.CreateFailWriter(0)
				gomega.Expect(writeHeaders(fw, map[string]string{"key": "value"})).To(matchFailure(FailWriteHeaders))
			})
		})
		ginkgo.When("it's closed after header content", func() {
			ginkgo.It("should fail with correct error", func() {
				fw := testutils.CreateFailWriter(1)
				gomega.Expect(writeHeaders(fw, map[string]string{"key": "value"})).To(matchFailure(FailWriteHeaders))
			})
		})
	})
	ginkgo.When("default port is not specified", func() {
		ginkgo.It("should use the default SMTP port when not specified", func() {
			testURL := "smtp://example.com/?fromAddress=sender@example.com&toAddresses=rec1@example.com"
			serviceURL := testutils.URLMust(testURL)
			err := service.Initialize(serviceURL, logger)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(service.Config.Port).To(gomega.Equal(uint16(DefaultSMTPPort)))
		})
	})
	ginkgo.It("returns the correct service identifier", func() {
		gomega.Expect(service.GetID()).To(gomega.Equal("smtp"))
	})
})

func testSendRecipient(testURL string, responses []string) failures.Failure {
	serviceURL, err := url.Parse(testURL)
	if err != nil {
		return standard.Failure(standard.FailParseURL, err)
	}

	err = service.Initialize(serviceURL, logger)
	if err != nil {
		return failures.Wrap("error parsing URL", standard.FailTestSetup, err)
	}

	if err := service.SetTemplateString("plain", "{{.message}}"); err != nil {
		return failures.Wrap("error setting plain template", standard.FailTestSetup, err)
	}

	textCon, tcfaker := testutils.CreateTextConFaker(responses, "\r\n")

	client := &smtp.Client{
		Text: textCon,
	}

	fakeTLSEnabled(client, serviceURL.Hostname())

	config := &Config{}
	message := "message body"

	ferr := service.sendToRecipient(client, "r@example.com", config, message)

	logger.Printf("\n%s", tcfaker.GetConversation(false))

	if ferr != nil {
		return ferr
	}

	return nil
}

func testIntegration(testURL string, responses []string, htmlTemplate string, plainTemplate string, expectRec ...string) failures.Failure {
	serviceURL, err := url.Parse(testURL)
	if err != nil {
		return standard.Failure(standard.FailParseURL, err)
	}

	if err = service.Initialize(serviceURL, logger); err != nil {
		return standard.Failure(standard.FailServiceInit, err)
	}

	if htmlTemplate != "" {
		if err := service.SetTemplateString("HTML", htmlTemplate); err != nil {
			return failures.Wrap("error setting HTML template", standard.FailTestSetup, err)
		}
	}

	if plainTemplate != "" {
		if err := service.SetTemplateString("plain", plainTemplate); err != nil {
			return failures.Wrap("error setting plain template", standard.FailTestSetup, err)
		}
	}

	textCon, tcfaker := testutils.CreateTextConFaker(responses, "\r\n")

	client := &smtp.Client{
		Text: textCon,
	}

	fakeTLSEnabled(client, serviceURL.Hostname())

	ferr := service.doSend(client, "Test message", service.Config)

	received := tcfaker.GetClientSentences()
	for _, expected := range expectRec {
		gomega.Expect(received).To(gomega.ContainElement(expected))
	}

	logger.Printf("\n%s", tcfaker.GetConversation(false))

	if ferr != nil {
		return ferr
	}

	return nil
}

// fakeTLSEnabled tricks a given client into believing that TLS is enabled even though it's not
// this is needed because the SMTP library won't allow plain authentication without TLS being turned on.
// having it turned on would of course mean that we cannot test the communication since it will be encrypted.
func fakeTLSEnabled(client *smtp.Client, hostname string) {
	// set the "tls" flag on the client which indicates that TLS encryption is enabled (even though it's not)
	cr := reflect.ValueOf(client).Elem().FieldByName("tls")
	cr = reflect.NewAt(cr.Type(), unsafe.Pointer(cr.UnsafeAddr())).Elem()
	cr.SetBool(true)

	// set the serverName field on the client which is used to identify the server and has to equal the hostname
	cr = reflect.ValueOf(client).Elem().FieldByName("serverName")
	cr = reflect.NewAt(cr.Type(), unsafe.Pointer(cr.UnsafeAddr())).Elem()
	cr.SetString(hostname)
}

// matchFailure is a simple wrapper around `fail` and `gomega.MatchError` to make it easier to use in tests.
func matchFailure(id failures.FailureID) gomegaTypes.GomegaMatcher {
	return gomega.MatchError(fail(id, nil))
}
