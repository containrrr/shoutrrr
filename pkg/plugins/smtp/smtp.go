package smtp

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/smtp"
)

// Plugin sends notifications to a given e-mail addresses via SMTP
type Plugin struct {}

const debugMode = false

// Send a notification message to discord
func (plugin *Plugin) Send(url string, message string) error {
	config, err := plugin.CreateConfigFromURL(url)
	if err != nil {
		return err
	}

	return doSend(message, config)
}


func doSend(message string, config Config) error {

	for _, toAddress := range config.ToAddresses {

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

		if auth := getAuth(config); auth != nil {
			if err := client.Auth(auth); err != nil {
				return fmt.Errorf("error authenticating: %s", err)
			}
		}


		// Set the sender and recipient first
		if err := client.Mail(config.FromAddress); err != nil {
			log.Fatalf("error creating new message: %s", err)
		}
		if err := client.Rcpt(toAddress); err != nil {
			log.Fatalf("error setting RCPT: %s", err)
		}

		// Send the email body.
		wc, err := client.Data()
		if err != nil {
			log.Fatalf("error creating message stream: %s", err)
		}

		if err := writeHeaders(&wc, map[string]string {
			"Subject": config.Subject,
			"To": toAddress,
			"From": fmt.Sprintf("%s <%s>", config.FromName, config.FromAddress),
			"MIME-version": "1.0;",
			"Content-Type": "text/html; charset=\"UTF-8\"",
		}); err != nil {
			return fmt.Errorf("error writing message headers: %s", err)
		}

		if _, err = fmt.Fprintf(wc, message); err != nil {
			return fmt.Errorf("error writing message: %s", err)
		}

		if err = wc.Close(); err != nil {
			return fmt.Errorf("error closing message stream: %s", err)
		}

		// Send the QUIT command and close the connection.
		err = client.Quit()
		if err != nil {
			return fmt.Errorf("error closing session: %s", err)
		}

		if debugMode {
			fmt.Printf("Mail successfully sent to \"%s\"!\n", toAddress)
		}
	}

	return nil
}

func getAuth(config Config) smtp.Auth {

		switch config.Auth {
			case Auth.None:
				return nil
			case Auth.Plain:
				return smtp.PlainAuth("", config.Username, config.Password, config.Host)
			case Auth.CRAMMD5:
				return smtp.CRAMMD5Auth(config.Username, config.Password)
			default:
				panic("invalid authorization method")
		}

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