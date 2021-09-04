package smtp

import (
	"crypto/tls"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/format"
	"io"
	"math/rand"
	"net"
	"net/smtp"
	"net/url"
	"time"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Service sends notifications to a given e-mail addresses via SMTP
type Service struct {
	standard.Standard
	standard.Templater
	config            *Config
	multipartBoundary string
	propKeyResolver   format.PropKeyResolver
}

const (
	contentHTML      = "text/html; charset=\"UTF-8\""
	contentPlain     = "text/plain; charset=\"UTF-8\""
	contentMultipart = "multipart/alternative; boundary=%s"
)

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{
		Port:        25,
		ToAddresses: nil,
		Subject:     "",
		Auth:        AuthTypes.Unknown,
		UseStartTLS: true,
		UseHTML:     false,
		Encryption:  EncMethods.Auto,
	}

	pkr := format.NewPropKeyResolver(service.config)

	if err := service.config.setURL(&pkr, configURL); err != nil {
		return err
	}

	if service.config.Auth == AuthTypes.Unknown {
		if service.config.Username != "" {
			service.config.Auth = AuthTypes.Plain
		} else {
			service.config.Auth = AuthTypes.None
		}
	}

	service.propKeyResolver = pkr

	return nil
}

// Send a notification message to e-mail recipients
func (service *Service) Send(message string, params *types.Params) error {
	client, err := getClientConnection(service.config)
	if err != nil {
		return fail(FailGetSMTPClient, err)
	}

	config := service.config.Clone()
	if err := service.propKeyResolver.UpdateConfigFromParams(&config, params); err != nil {
		return fail(FailApplySendParams, err)
	}

	return service.doSend(client, message, &config)
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

func (service *Service) doSend(client *smtp.Client, message string, config *Config) failure {

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

	if auth, err := service.getAuth(config); err != nil {
		return err
	} else if auth != nil {
		if err := client.Auth(auth); err != nil {
			return fail(FailAuthenticating, err)
		}
	}

	for _, toAddress := range config.ToAddresses {

		err := service.sendToRecipient(client, toAddress, config, message)
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

func (service *Service) getAuth(config *Config) (smtp.Auth, failure) {

	switch config.Auth {
	case AuthTypes.None:
		return nil, nil
	case AuthTypes.Plain:
		return smtp.PlainAuth("", config.Username, config.Password, config.Host), nil
	case AuthTypes.CRAMMD5:
		return smtp.CRAMMD5Auth(config.Username, config.Password), nil
	case AuthTypes.OAuth2:
		return OAuth2Auth(config.Username, config.Password), nil
	default:
		return nil, fail(FailAuthType, nil, config.Auth.String())
	}

}

func (service *Service) sendToRecipient(client *smtp.Client, toAddress string, config *Config, message string) failure {

	// Set the sender and recipient first
	if err := client.Mail(config.FromAddress); err != nil {
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

	if err := writeHeaders(wc, service.getHeaders(toAddress, config.Subject)); err != nil {
		return fail(FailWriteHeaders, err)
	}

	var ferr failure
	if config.UseHTML {
		ferr = service.writeMultipartMessage(wc, message)
	} else {
		ferr = service.writeMessagePart(wc, message, "plain")
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
		"Date":         time.Now().Format(time.RFC1123Z),
		"To":           toAddress,
		"From":         fmt.Sprintf("%s <%s>", conf.FromName, conf.FromAddress),
		"MIME-version": "1.0",
		"Content-Type": contentType,
	}
}

func (service *Service) writeMultipartMessage(wc io.WriteCloser, message string) failure {

	if err := writeMultipartHeader(wc, service.multipartBoundary, contentPlain); err != nil {
		return fail(FailPlainHeader, err)
	}
	if err := service.writeMessagePart(wc, message, "plain"); err != nil {
		return err
	}

	if err := writeMultipartHeader(wc, service.multipartBoundary, contentHTML); err != nil {
		return fail(FailHTMLHeader, err)
	}
	if err := service.writeMessagePart(wc, message, "HTML"); err != nil {
		return err
	}

	if err := writeMultipartHeader(wc, service.multipartBoundary, ""); err != nil {
		return fail(FailMultiEndHeader, err)

	}

	return nil
}

func (service *Service) writeMessagePart(wc io.WriteCloser, message string, template string) failure {
	if tpl, found := service.GetTemplate(template); found {
		data := make(map[string]string)
		data["message"] = message
		if err := tpl.Execute(wc, data); err != nil {
			return fail(FailMessageTemplate, err)
		}
	} else {
		if _, err := fmt.Fprintf(wc, message); err != nil {
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
