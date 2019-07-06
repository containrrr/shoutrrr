package smtp

import (
	"fmt"
	"log"
	"net/smtp"
	"net/url"
	"os"
	"reflect"
	"testing"
	"unsafe"

	"github.com/containrrr/shoutrrr/pkg/failures"
	"github.com/containrrr/shoutrrr/pkg/testutils"
	"github.com/containrrr/shoutrrr/pkg/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSMTP(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr SMTP Suite")
}

var (
	service    *Service
	envSMTPURL string
	logger     *log.Logger
)

var _ = Describe("the SMTP service", func() {

	BeforeSuite(func() {

		envSMTPURL = os.Getenv("SHOUTRRR_SMTP_URL")
		logger = util.TestLogger()
	})
	BeforeEach(func() {
		service = &Service{}

	})
	When("parsing the configuration URL", func() {
		It("should be identical after de-/serialization", func() {
			testURL := "smtp://user:password@example.com:2225/?fromAddress=sender@example.com&fromName=Sender&toAddresses=rec1@example.com,rec2@example.com&auth=None&subject=Subject&startTls=No&useHTML=No"

			url, err := url.Parse(testURL)
			Expect(err).NotTo(HaveOccurred(),"parsing")

			config := &Config{}
			err = config.SetURL(url)
			Expect(err).NotTo(HaveOccurred(),"verifying")

			outputURL := config.GetURL()

			// fmt.Println(outputURL.String())

			Expect(outputURL.String()).To(Equal(testURL))

		})
		When("fromAddress is missing", func() {
			It("should return an error", func() {
				testURL := "smtp://user:password@example.com:2225/?toAddresses=rec1@example.com,rec2@example.com"

				url, err := url.Parse(testURL)
				Expect(err).NotTo(HaveOccurred(), "parsing")

				config := &Config{}
				err = config.SetURL(url)
				Expect(err).To(HaveOccurred(), "verifying")
			})
		})
		When("toAddresses are missing", func() {
			It("should return an error", func() {
				testURL := "smtp://user:password@example.com:2225/?fromAddress=sender@example.com"

				url, err := url.Parse(testURL)
				Expect(err).NotTo(HaveOccurred(), "parsing")

				config := &Config{}
				err = config.SetURL(url)
				Expect(err).To(HaveOccurred(), "verifying")
			})

		})
	})
	It("should implement basic service API methods correctly", func() {
		testutils.TestConfigGetInvalidQueryValue(&Config{})
		testutils.TestConfigSetInvalidQueryValue(&Config{}, "smtp://example.com/?fromAddress=s@example.com&toAddresses=r@example.com&foo=bar")

		testutils.TestConfigGetEnumsCount(&Config{}, 1)
		testutils.TestConfigGetFieldsCount(&Config{}, 7)
	})

	When("the service is not configured correctly", func() {
		It("should fail to send messages", func(){
			service := Service{
				config: &Config{},
			}
			err := service.Send("test message", nil)
			Expect(err).To(HaveOccurred())
		})
	})

	When("the underlying stream stops working", func() {
		var service Service
		var params map[string]string
		BeforeEach(func(){
			service = Service{}
			params = make(map[string]string, 0)
		})
		It("should fail when writing multipart plain header", func(){
			writer := testutils.CreateFailWriter(1)
			err := service.writeMultipartMessage(writer, &params)
			fmt.Printf("%+v\n", err)
			Expect(err).To(HaveOccurred())
			Expect(err.ID()).To(Equal(FailPlainHeader))
		})

		It("should fail when writing multipart plain message", func(){
			writer := testutils.CreateFailWriter(2)
			err := service.writeMultipartMessage(writer, &params)
			fmt.Printf("%+v\n", err)
			Expect(err).To(HaveOccurred())
			Expect(err.ID()).To(Equal(FailMessageRaw))
		})

		It("should fail when writing multipart HTML header", func(){
			writer := testutils.CreateFailWriter(4)
			err := service.writeMultipartMessage(writer, &params)
			fmt.Printf("%+v\n", err)
			Expect(err).To(HaveOccurred())
			Expect(err.ID()).To(Equal(FailHTMLHeader))
		})

		It("should fail when writing multipart HTML message", func(){
			writer := testutils.CreateFailWriter(5)
			err := service.writeMultipartMessage(writer, &params)
			fmt.Printf("%+v\n", err)
			Expect(err).To(HaveOccurred())
			Expect(err.ID()).To(Equal(FailMessageRaw))
		})

		It("should fail when writing multipart end header", func(){
			writer := testutils.CreateFailWriter(6)
			err := service.writeMultipartMessage(writer, &params)
			fmt.Printf("%+v\n", err)
			Expect(err).To(HaveOccurred())
			Expect(err.ID()).To(Equal(FailMultiEndHeader))
		})

		It("should fail when writing message template", func(){
			writer := testutils.CreateFailWriter(0)
			e := service.SetTemplateString("dummy", "dummy template content")
			Expect(e).ToNot(HaveOccurred())

			err := service.writeMessagePart(writer, &params, "dummy")
			fmt.Printf("%+v\n", err)
			Expect(err).To(HaveOccurred())
			Expect(err.ID()).To(Equal(FailMessageTemplate))
		})

	})


	When("running E2E tests", func() {

		It("should work without errors", func() {
			if envSMTPURL == "" {
				Skip("environment not set up for E2E testing")
				return
			}

			serviceURL, err := url.Parse(envSMTPURL)
			Expect(err).NotTo(HaveOccurred())

			err = service.Initialize(serviceURL, logger)
			Expect(err).NotTo(HaveOccurred())


			err = service.Send( "this is an integration test", nil)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	When("running integration tests", func() {

		When("given a typical usage case configuration URL", func() {

			It("should send notifications without any errors", func() {
				testURL := "smtp://user:password@example.com:2225/?startTLS=no&fromAddress=sender@example.com&toAddresses=rec1@example.com,rec2@example.com&useHTML=yes"
				err := testIntegration(testURL, []string {
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
				if msg, test := failures.IsTestSetupFailure(err); test {
					Skip(msg)
					return
				}
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("given a configuration URL with authentication disabled", func() {

			It("should send notifications without any errors", func() {
				testURL := "smtp://example.com:2225/?startTLS=no&fromAddress=sender@example.com&toAddresses=rec1@example.com&useHTML=no"
				err := testIntegration(testURL, []string {
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
				if msg, test := failures.IsTestSetupFailure(err); test {
					Skip(msg)
					return
				}
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("server communication fails", func() {

			It("should fail when not being able to enable StartTLS", func(){
				testURL := "smtp://example.com:2225/?startTLS=yes&auth=none&fromAddress=sender@example.com&toAddresses=rec1@example.com&useHTML=no"
				err := testIntegration(testURL, []string {
					"250-mx.google.com at your service",
					"250-SIZE 35651584",
					"250-AUTH LOGIN PLAIN",
					"250 8BITMIME",
					"502 That's too hard",
				}, "", "")
				if msg, test := failures.IsTestSetupFailure(err); test {
					Skip(msg)
					return
				}
				Expect(err).To(HaveOccurred())
				Expect(err.ID()).To(Equal(FailEnableStartTLS))
			})

			It("should fail when authentication type is invalid", func(){
				testURL := "smtp://example.com:2225/?startTLS=no&auth=bad&fromAddress=sender@example.com&toAddresses=rec1@example.com&useHTML=no"
				err := testIntegration(testURL, []string {}, "", "")
				if msg, test := failures.IsTestSetupFailure(err); test {
					Skip(msg)
					return
				}
				Expect(err).To(HaveOccurred())
				Expect(err.ID()).To(Equal(FailAuthType))
			})

			It("should fail when not being able to use authentication type", func(){
				testURL := "smtp://example.com:2225/?startTLS=no&auth=crammd5&fromAddress=sender@example.com&toAddresses=rec1@example.com&useHTML=no"
				err := testIntegration(testURL, []string {
					"250-mx.google.com at your service",
					"250-SIZE 35651584",
					"250-AUTH LOGIN PLAIN",
					"250 8BITMIME",
					"504 Liar",
				}, "", "")
				if msg, test := failures.IsTestSetupFailure(err); test {
					Skip(msg)
					return
				}
				Expect(err).To(HaveOccurred())
				Expect(err.ID()).To(Equal(FailAuthenticating))
			})

			It("should fail when not being able to send to recipient", func(){
				testURL := "smtp://example.com:2225/?startTLS=no&auth=none&fromAddress=sender@example.com&toAddresses=rec1@example.com&useHTML=no"
				err := testIntegration(testURL, []string {
					"250-mx.google.com at your service",
					"250-SIZE 35651584",
					"250-AUTH LOGIN PLAIN",
					"250 8BITMIME",
					"551 I don't know you",
				}, "", "")
				if msg, test := failures.IsTestSetupFailure(err); test {
					Skip(msg)
					return
				}
				Expect(err).To(HaveOccurred())
				Expect(err.ID()).To(Equal(FailSendRecipient))
			})

			It("should fail when the recipient is not accepted", func(){
				testURL := "smtp://example.com:2225/?startTLS=no&auth=none&fromAddress=sender@example.com&toAddresses=rec1@example.com&useHTML=no"
				err := testSendRecipient(testURL, []string {
					"250 mx.google.com at your service",
					"250 Sender OK",
					"553 She doesn't want to be disturbed",
				})
				if msg, test := failures.IsTestSetupFailure(err); test {
					Skip(msg)
					return
				}
				Expect(err).To(HaveOccurred())
				Expect(err.ID()).To(Equal(FailSetRecipient))
			})

			It("should fail when the server does not accept the data stream", func(){
				testURL := "smtp://example.com:2225/?startTLS=no&auth=none&fromAddress=sender@example.com&toAddresses=rec1@example.com&useHTML=no"
				err := testSendRecipient(testURL, []string {
					"250 mx.google.com at your service",
					"250 Sender OK",
					"250 Receiver OK",
					"554 Nah I'm fine thanks",
				})
				if msg, test := failures.IsTestSetupFailure(err); test {
					Skip(msg)
					return
				}
				Expect(err).To(HaveOccurred())
				Expect(err.ID()).To(Equal(FailOpenDataStream))
			})

			It("should fail when the server does not accept the data stream content", func(){
				testURL := "smtp://example.com:2225/?startTLS=no&auth=none&fromAddress=sender@example.com&toAddresses=rec1@example.com&useHTML=no"
				err := testSendRecipient(testURL, []string {
					"250 mx.google.com at your service",
					"250 Sender OK",
					"250 Receiver OK",
					"354 Go ahead",
					"554 Such garbage!",
				})
				if msg, test := failures.IsTestSetupFailure(err); test {
					Skip(msg)
					return
				}
				Expect(err).To(HaveOccurred())
				Expect(err.ID()).To(Equal(FailCloseDataStream))
			})

			It("should fail when the server does not close the connection gracefully", func(){
				testURL := "smtp://example.com:2225/?startTLS=no&auth=none&fromAddress=sender@example.com&toAddresses=rec1@example.com&useHTML=no"
				err := testIntegration(testURL, []string {
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
				if msg, test := failures.IsTestSetupFailure(err); test {
					Skip(msg)
					return
				}
				Expect(err).To(HaveOccurred())
				Expect(err.ID()).To(Equal(FailClosingSession))
			})

		})
	})
})

func testSendRecipient(testURL string, responses []string) failure {
	serviceURL, err := url.Parse(testURL)
	if err != nil {
		return failures.Wrap("error parsing URL", failures.FailTestSetup, err)
	}

	err = service.Initialize(serviceURL, logger)
	if err != nil {
		return failures.Wrap("error parsing URL", failures.FailTestSetup, err)
	}

	if err := service.SetTemplateString("plain", "{{.message}}"); err != nil {
		return failures.Wrap("error setting plain template", failures.FailTestSetup, err)
	}


	textCon, tcfaker := testutils.CreateTextConFaker(responses, "\r\n")

	client := &smtp.Client{
		Text: textCon,
	}

	fakeTLSEnabled(client, serviceURL.Hostname())

	params := &map[string]string{
		"message": "message body",
	}

	ferr := service.sendToRecipient(client, "r@example.com", params)

	logger.Printf("\n%s", tcfaker.GetConversation(false))
	if ferr != nil {
		return ferr
	}

	return nil
}

func testIntegration(testURL string, responses []string, htmlTemplate string, plainTemplate string) failure {

	serviceURL, err := url.Parse(testURL)
	if err != nil {
		return failures.Wrap("error parsing URL", failures.FailTestSetup, err)
	}

	err = service.Initialize(serviceURL, logger)
	if err != nil {
		return failures.Wrap("error parsing URL", failures.FailTestSetup, err)
	}

	if htmlTemplate != "" {
		if err := service.SetTemplateString("HTML", htmlTemplate); err != nil {
			return failures.Wrap("error setting HTML template", failures.FailTestSetup, err)
		}
	}
	if plainTemplate != "" {
		if err := service.SetTemplateString("plain", plainTemplate); err != nil {
			return failures.Wrap("error setting plain template", failures.FailTestSetup, err)
		}
	}

	textCon, tcfaker := testutils.CreateTextConFaker(responses, "\r\n")

	client := &smtp.Client{
		Text: textCon,
	}

	fakeTLSEnabled(client, serviceURL.Hostname())

	ferr := service.doSend(client, "Test message", map[string]string{})

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