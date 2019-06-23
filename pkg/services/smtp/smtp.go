package smtp

import (
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"net/smtp"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Service sends notifications to a given e-mail addresses via SMTP
type Service struct {
	standard.Standard
	standard.Templater
	config *Config
	multipartBoundry string
}

const (
	contentHTML = "text/html; charset=\"UTF-8\""
	contentPlain = "text/plain; charset=\"UTF-8\""
	contentMultipart = "multipart/alternative; boundary=%s"
)

// NewConfig returns an empty ServiceConfig for this Service
func (service *Service) NewConfig() types.ServiceConfig {
	return &Config{}
}

// Send a notification message to discord
func (service *Service) Send(message string, params *map[string]string) error {
	if params == nil {
		params = &map[string]string{}
	}
	return service.doSend(message, *params)
}

func (service *Service) doSend(message string, params map[string]string) error {
	config := service.config

	params["message"] = message

	if config.UseHTML {
		service.multipartBoundry = fmt.Sprintf("%x", rand.Int63())
	}

	client, err := smtp.Dial(fmt.Sprintf("%s:%d", config.Host, config.Port))
	if err != nil {
		return fmt.Errorf("error connecting to server: %s", err)
	}

	if config.UseStartTLS {
		if err := client.StartTLS(&tls.Config{
			ServerName: config.Host,
		}); err != nil {
			return fmt.Errorf("error enabling StartTLS message: %s", err)
		}
	}

	if auth, err := service.getAuth(); err != nil {
		return err
	} else if auth != nil {
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("error authenticating: %s", err)
		}
	}

	for _, toAddress := range config.ToAddresses {

		err := service.sendToRecipient(client, toAddress, &params)
		if err != nil {
			return fmt.Errorf("error sending message to recipient: %s", err)
		}

		service.Logf("Mail successfully sent to \"%s\"!\n", toAddress)
	}

	// Send the QUIT command and close the connection.
	err = client.Quit()
	if err != nil {
		return fmt.Errorf("error closing session: %s", err)
	}

	return nil
}

func (service *Service) getAuth() (smtp.Auth, error) {

		config := service.config

		switch config.Auth {
			case authTypes.None:
				return nil, nil
			case authTypes.Plain:
				return smtp.PlainAuth("", config.Username, config.Password, config.Host), nil
			case authTypes.CRAMMD5:
				return smtp.CRAMMD5Auth(config.Username, config.Password), nil
			default:
				return nil, fmt.Errorf("invalid authorization method '%s'", config.Auth.String())
		}

}

func (service *Service) sendToRecipient(client *smtp.Client, toAddress string, params *map[string]string) error {
	conf := service.config

	// Set the sender and recipient first
	if err := client.Mail(conf.FromAddress); err != nil {
		return fmt.Errorf("error creating new message: %s", err)
	}
	if err := client.Rcpt(toAddress); err != nil {
		return fmt.Errorf("error setting RCPT: %s", err)
	}

	// Send the email body.
	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("error creating message stream: %s", err)
	}

	// TODO: Move param override to shared service API
	subject, found := (*params)["subject"]
	if !found {
		subject = conf.Subject
	}

	if err := writeHeaders(&wc, service.getHeaders(toAddress, subject)); err != nil {
		return fmt.Errorf("error writing message headers: %s", err)
	}

	if conf.UseHTML {
		err = service.writeMultipartMessage(&wc, params)
	} else {
		err = writePlainMessage(&wc, (*params)["message"])
	}

	if err != nil {
		return err
	}

	if err = wc.Close(); err != nil {
		return fmt.Errorf("error closing message stream: %s", err)
	}

	return nil
}

func (service *Service) getHeaders(toAddress string, subject string) map[string]string {
	conf := service.config

	var contentType string
	if conf.UseHTML {
		contentType = fmt.Sprintf(contentMultipart, service.multipartBoundry)
	} else {
		contentType = contentPlain
	}

	return map[string]string {
		"Subject": subject,
		"To": toAddress,
		"From": fmt.Sprintf("%s <%s>", conf.FromName, conf.FromAddress),
		"MIME-version": "1.0;",
		"Content-Type": contentType,
	}
}

func (service *Service) writeMultipartMessage(wc *io.WriteCloser, params *map[string]string) error {

	message := (*params)["message"]

	if err := writeMultipartHeader(wc, service.multipartBoundry, contentPlain); err != nil {
		return fmt.Errorf("error writing message: %s", err)
	}

	if err := writePlainMessage(wc, message); err != nil {
		return err
	}

	if err := writeMultipartHeader(wc, service.multipartBoundry, contentHTML); err != nil {
		return fmt.Errorf("error writing message: %s", err)
	}

	if tpl, found := service.GetTemplate("message", ); found {
		if err := tpl.Execute(*wc, params); err != nil {
			return fmt.Errorf("error applying message template: %s", err)
		}
	} else {
		if _, err := fmt.Fprintf(*wc, message); err != nil {
			return fmt.Errorf("error writing message: %s", err)
		}
	}

	if err := writeMultipartHeader(wc, service.multipartBoundry, ""); err != nil {
		return fmt.Errorf("error writing message: %s", err)
	}

	return nil
}

func writePlainMessage(wc *io.WriteCloser, message string) error {
	if _, err := fmt.Fprintf(*wc, message); err != nil {
		return fmt.Errorf("error writing message: %s", err)
	}
	return nil
}

func writeMultipartHeader(wc *io.WriteCloser, boundry string, contentType string) error {
	suffix := "\n"
	if len(contentType) < 1 {
		suffix = "--"
	}

	if _, err := fmt.Fprintf(*wc, "\n\n--%s%s", boundry, suffix); err != nil {
		return err
	}

	if len(contentType) > 0 {
		if _, err := fmt.Fprintf(*wc, "Content-Type: %s\n\n", contentType); err != nil {
			return err
		}
	}

	return nil
}

func writeHeaders(wc *io.WriteCloser, headers map[string]string) error {
	for key, val := range headers {
		if _, err := fmt.Fprintf(*wc, "%s: %s\n", key, val); err != nil {
			return err
		}
	}

	_, err := fmt.Fprintln(*wc)
	return err
}