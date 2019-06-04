package smtp

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"errors"
)

// Config is the configuration needed to send discord notifications
type Config struct {
	Host        string
	Username    string
	Password    string
	Port        uint16
	FromAddress string
	FromName    string
	ToAddresses []string
	Subject     string
	Auth        AuthType
	UseStartTLS bool
}

// CreateAPIURLFromConfig takes a discord config object and creates a post url
func CreateAPIURLFromConfig(config Config) string {
	return fmt.Sprintf(
		URLFormat,
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.FromAddress,
		config.FromName,
		strings.Join(config.ToAddresses, ","),
	)
}

// CreateConfigFromURL creates a Config struct given a valid discord notification url
func (plugin *Plugin) CreateConfigFromURL(rawURL string) (Config, error) {
	// args, err := ExtractArguments(url)

	url, err := url.Parse(rawURL)


	if err != nil {
		return Config{}, err
	}

	hostParts := strings.Split(url.Host, ":")
	host := hostParts[0]
	port, err := strconv.ParseUint(hostParts[1], 10, 16)
	password, _ := url.User.Password()

	config := Config{
		Username: url.User.Username(),
		Password: password,
		Host: host,
		Port: uint16(port),
		FromAddress: "",
		FromName: "Shoutrrr",
		ToAddresses: make([]string, 0),
		Subject: "A message from Shoutrrr",
		UseStartTLS: true,
	}

	if err := config.UpdateFromValues(url.Query()); err != nil {
		return Config{}, err
	}

	if debugMode {
		log.Printf("Username: %s", config.Username)
		log.Printf("Password: %s", config.Password)
		log.Printf("Host: %s", config.Host)
		log.Printf("Auth: %s", String(config.Auth))
		log.Printf("UseStartTLS: %t", config.UseStartTLS)
		log.Printf("FromName: %s", config.FromName)
		log.Printf("FromAddress: %s", config.FromAddress)
		log.Printf("Subject: %s", config.Subject)
	}


	if len(config.FromAddress) < 1 {
		return Config{}, errors.New("fromAddress missing from config URL")
	}

	if len(config.ToAddresses) < 1 {
		return Config{}, errors.New("toAddress missing from config URL")
	}

	return config, nil
}

func (config *Config) UpdateFromValues(values url.Values) error {

	for key, vals := range values {

		if len(vals) > 1 {
			fmt.Printf("warning: %s additional value ignored!: %s\n", key, vals[1])
		}
		val := vals[0]

		if debugMode {
			fmt.Printf("Query \"%s\" => \"%s\"\n", key, val)
		}

		switch key {
		case "fromAddress":
			fallthrough
		case "from":
			config.FromAddress = val
		case "fromName":
			config.FromName = val
		case "toAddresses":
			fallthrough
		case "to":
			config.ToAddresses = strings.Split(val, ",")
		case "auth":
			switch strings.ToLower(val) {
			case "none":
				config.Auth = Auth.None
			case "plain":
				config.Auth = Auth.Plain
			case "crammd5":
				config.Auth = Auth.CRAMMD5
			}
		case "subject":
			config.Subject = val
		case "startTls":
			switch strings.ToLower(val) {
			case "true":
				fallthrough
			case "1":
				config.UseStartTLS = true
			case "false":
				fallthrough
			case "0":
				config.UseStartTLS = false
			default:

			}
		default:
			return fmt.Errorf("invalid query key \"%s\"", key)
		}
	}
	return nil
}

type urlParts struct {
	Host uint8
	Username uint8
	Password uint8
	Port uint8
	Query uint8
}

var URLPart = &urlParts{
	Username: 0,
	Password: 1,
	Host: 2,
	Port: 3,
	Query: 4,
}

type AuthType = uint8

type authType struct {
	None AuthType
	Plain AuthType
	CRAMMD5 AuthType
}

var Auth = &authType{
	None: 0,
	Plain : 1,
	CRAMMD5: 2,
}

func String(at AuthType) string {
	switch at {
		case Auth.None: return "None"
		case Auth.Plain : return "Plain"
		case Auth.CRAMMD5: return "CRAMMD5"
		default: return "Unknown"
	}
}

const URLPattern = "^smtp://([^:]+):([^@]+)@([^:]+):([1-9][0-9]*)/\\?(.*)$"
const URLFormat = "smtp://%s:%s@%s:%d/?fromAddress=%s&fromName=%s&toAddresses=%s"