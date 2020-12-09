package smtp

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/smtp"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Service sends notifications to a given e-mail addresses via SMTP
type Service struct {
	standard.Standard
	standard.Templater
	config            *Config
	multipartBoundary string
}

const (
	contentHTML      = "text/html; charset=\"UTF-8\""
	contentPlain     = "text/plain; charset=\"UTF-8\""
	contentMultipart = "multipart/alternative; boundary=%s"
)

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger *log.Logger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{
		Port:        25,
		ToAddresses: nil,
		Subject:     "",
		Auth:        authTypes.Unknown,
		UseStartTLS: true,
		UseHTML:     false,
		Encryption:  encMethods.Auto,
	}
	if err := service.config.SetURL(configURL); err != nil {
		return err
	}

	if service.config.Auth == authTypes.Unknown {
		if service.config.Username != "" {
			service.config.Auth = authTypes.Plain
		} else {
			service.config.Auth = authTypes.None
		}
	}

	return nil
}

// Send a notification message to e-mail recipients
func (service *Service) Send(message string, params *types.Params) error {
	if params == nil {
		params = &types.Params{}
	}
	client, err := getClientConnection(service.config)
	if err != nil {
		return fail(FailGetSMTPClient, err)
	}
	return service.doSend(client, message, *params)
}

func getClientConnection(config *Config) (*smtp.Client, error) {

	var conn net.Conn
	var err error

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	if useImplicitTLS(config.Encryption, config.Port) {
		conn, err = tls.Dial("tcp", addr, &tls.Config{
			ServerName: config.Host,
		})
	} else {
		conn, err = net.Dial("tcp", addr)
	}

	if err != nil {
		return nil, fail(FailConnectToServer, err)
	}

	client, err := smtp.NewClient(conn, config.Host)
	if err != nil {
		return nil, fail(FailCreateSMTPClient, err)
	}

	return client, nil
}

func (service *Service) doSend(client *smtp.Client, message string, params map[string]string) failure {
	config := service.config

	params["message"] = message

	if config.UseHTML {
		service.multipartBoundary = fmt.Sprintf("%x", rand.Int63())
	}

	if config.UseStartTLS && !useImplicitTLS(config.Encryption, config.Port) {
		if supported, _ := client.Extension("StartTLS"); !supported {
			service.Logf("Warning: StartTLS enabled, but server did not report support for it. Connection is NOT encrypted")
		} else {
			if err := client.StartTLS(&tls.Config{
				ServerName: config.Host,
			}); err != nil {
				return fail(FailEnableStartTLS, err)
			}
		}
	}

	if auth, err := service.getAuth(); err != nil {
		return err
	} else if auth != nil {
		if err := client.Auth(auth); err != nil {
			return fail(FailAuthenticating, err)
		}
	}

	for _, toAddress := range config.ToAddresses {

		err := service.sendToRecipient(client, toAddress, &params)
		if err != nil {
			return fail(FailSendRecipient, err)
		}

		service.Logf("Mail successfully sent to \"%s\"!\n", toAddress)
	}

	// Send the QUIT command and close the connection.
	err := client.Quit()
	if err != nil {
		return fail(FailClosingSession, err)
	}

	return nil
}

func (service *Service) getAuth() (smtp.Auth, failure) {

	config := service.config

	switch config.Auth {
	case authTypes.None:
		return nil, nil
	case authTypes.Plain:
		return smtp.PlainAuth("", config.Username, config.Password, config.Host), nil
	case authTypes.CRAMMD5:
		return smtp.CRAMMD5Auth(config.Username, config.Password), nil
	case authTypes.OAuth2:
		return OAuth2Auth(config.Username, config.Password), nil
	default:
		return nil, fail(FailAuthType, nil, config.Auth.String())
	}

}

func (service *Service) sendToRecipient(client *smtp.Client, toAddress string, params *map[string]string) failure {
	conf := service.config

	// Set the sender and recipient first
	if err := client.Mail(conf.FromAddress); err != nil {
		return fail(FailSetSender, err)
	}
	if err := client.Rcpt(toAddress); err != nil {
		return fail(FailSetRecipient, err)
	}

	// Send the email body.
	wc, err := client.Data()
	if err != nil {
		return fail(FailOpenDataStream, err)
	}

	// TODO: Move param override to shared service API
	subject, found := (*params)["subject"]
	if !found {
		subject = conf.Subject
	}

	if err := writeHeaders(wc, service.getHeaders(toAddress, subject)); err != nil {
		return fail(FailWriteHeaders, err)
	}

	var ferr failure
	if conf.UseHTML {
		ferr = service.writeMultipartMessage(wc, params)
	} else {
		ferr = service.writeMessagePart(wc, params, "plain")
	}

	if ferr != nil {
		return ferr
	}

	if err = wc.Close(); err != nil {
		return fail(FailCloseDataStream, err)
	}

	return nil
}

func (service *Service) getHeaders(toAddress string, subject string) map[string]string {
	conf := service.config

	var contentType string
	if conf.UseHTML {
		contentType = fmt.Sprintf(contentMultipart, service.multipartBoundary)
	} else {
		contentType = contentPlain
	}

	return map[string]string{
		"Subject":      subject,
		"To":           toAddress,
		"From":         fmt.Sprintf("%s <%s>", conf.FromName, conf.FromAddress),
		"MIME-version": "1.0;",
		"Content-Type": contentType,
	}
}

func (service *Service) writeMultipartMessage(wc io.WriteCloser, params *map[string]string) failure {

	if err := writeMultipartHeader(wc, service.multipartBoundary, contentPlain); err != nil {
		return fail(FailPlainHeader, err)
	}
	if err := service.writeMessagePart(wc, params, "plain"); err != nil {
		return err
	}

	if err := writeMultipartHeader(wc, service.multipartBoundary, contentHTML); err != nil {
		return fail(FailHTMLHeader, err)
	}
	if err := service.writeMessagePart(wc, params, "HTML"); err != nil {
		return err
	}

	if err := writeMultipartHeader(wc, service.multipartBoundary, ""); err != nil {
		return fail(FailMultiEndHeader, err)

	}

	return nil
}

func (service *Service) writeMessagePart(wc io.WriteCloser, params *map[string]string, template string) failure {
	if tpl, found := service.GetTemplate(template); found {
		if err := tpl.Execute(wc, params); err != nil {
			return fail(FailMessageTemplate, err)
		}
	} else {
		if _, err := fmt.Fprintf(wc, (*params)["message"]); err != nil {
			return fail(FailMessageRaw, err)
		}
	}
	return nil
}

func writeMultipartHeader(wc io.WriteCloser, boundary string, contentType string) error {
	suffix := "\n"
	if len(contentType) < 1 {
		suffix = "--"
	}

	if _, err := fmt.Fprintf(wc, "\n\n--%s%s", boundary, suffix); err != nil {
		return err
	}

	if len(contentType) > 0 {
		if _, err := fmt.Fprintf(wc, "Content-Type: %s\n\n", contentType); err != nil {
			return err
		}
	}

	return nil
}

func writeHeaders(wc io.WriteCloser, headers map[string]string) error {
	for key, val := range headers {
		if _, err := fmt.Fprintf(wc, "%s: %s\n", key, val); err != nil {
			return err
		}
	}

	_, err := fmt.Fprintln(wc)
	return err
}
