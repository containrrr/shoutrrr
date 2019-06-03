package smtp

import (
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"regexp"
	"strconv"
	"strings"
)

// Plugin sends notifications to a given e-mail addresses via SMTP
type Plugin struct {}

// Send a notification message to discord
func (plugin *Plugin) Send(url string, message string) error {
	config, err := plugin.CreateConfigFromURL(url)
	if err != nil {
		return err
	}

	return doSend(message, config)
}

// CreateAPIURLFromConfig takes a discord config object and creates a post url
func CreateAPIURLFromConfig(config Config) string {
	return fmt.Sprintf(
		"smtp://%s:%s@%s:%d",
		config.Username,
		config.Password,
		config.Host,
		config.Port)
}

// CreateConfigFromURL creates a Config struct given a valid discord notification url
func (plugin *Plugin) CreateConfigFromURL(url string) (Config, error) {
	args, err := ExtractArguments(url)
	if err != nil {
		return Config{}, err
	}
	if len(args) != 2 {
		return Config{}, errors.New("invalid SMTP configuration URL")
	}

	port, err := strconv.ParseUint(args[3], 10, 16)

	return Config{
		Username: args[0],
		Password: args[1],
		Host: args[2],
		Port: uint16(port),
		FromAddress: args[4],
		FromName: args[5],
		ToAddresses: strings.Split(args[6], "/"),
	}, nil
}

// ExtractArguments extracts the arguments from a notification url, i.e everything following the initial ://
func ExtractArguments(url string) ([]string, error) {
	regex, err := regexp.Compile("^smtp://([^:]+):([^@]+)@([^:]+):([1-9][0-9]*)/([a-z]+@[a-z|\\-|\\.])\\(([^\\)])\\)/(.*)$")
		if err != nil {
		return nil, errors.New("could not compile regex")
	}
	match := regex.FindStringSubmatch(url)
	println(match)

	if len(match[1]) <= 0 {
		return nil, errors.New("could not extract any arguments")
	}
	return match, nil
}

func doSend(message string, config Config) error {


	fmt.Println(config)

	for _, toAddress := range config.ToAddresses {

		client, err := smtp.Dial(fmt.Sprintf("%s:%d", config.Host, config.Port)
		if err != nil {
			log.Fatal(err)
		}

		// Set the sender and recipient first
		if err := client.Mail(config.FromAddress); err != nil {
			log.Fatal(err)
		}
		if err := client.Rcpt(toAddress); err != nil {
			log.Fatal(err)
		}

		// Send the email body.
		wc, err := client.Data()
		if err != nil {
			log.Fatal(err)
		}
		_, err = fmt.Fprintf(wc, message)
		if err != nil {
			log.Fatal(err)
		}
		err = wc.Close()
		if err != nil {
			log.Fatal(err)
		}

		// Send the QUIT command and close the connection.
		err = client.Quit()
		if err != nil {
			log.Fatal(err)
		}

	}


	return nil
}