package smtp

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

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
	}, nil
}

// ExtractArguments extracts the arguments from a notification url, i.e everything following the initial ://
func ExtractArguments(url string) ([]string, error) {
	regex, err := regexp.Compile("^smtp://([^:]+):([^@]+)@([^:]+):([1-9][0-9]*)$")
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

func doSend(payload string, config Config) error {


	fmt.Println(config)


	return errors.New("not implemented")
}