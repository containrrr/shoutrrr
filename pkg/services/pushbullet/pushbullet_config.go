package pushbullet

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
)

// Config ...
type Config struct {
	standard.EnumlessConfig
	Targets []string
	Token   string
}

// GetURL returns a URL representation of it's current field values
func (config *Config) GetURL() *url.URL {
	return &url.URL{
		Host:       config.Token,
		Scheme:     Scheme,
		ForceQuery: false,
	}
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *Config) SetURL(url *url.URL) error {
	splitBySlash := func(c rune) bool {
		return c == '/'
	}

	path := strings.FieldsFunc(url.Path, splitBySlash)
	if url.Fragment != "" {
		path = append(path, fmt.Sprintf("#%s", url.Fragment))
	}
	if len(path) == 0 {
		path = []string{""}
	}

	config.Token = url.Host
	config.Targets = path[0:]

	if err := validateToken(config.Token); err != nil {
		return err
	}

	return nil
}

func validateToken(token string) error {
	if err := tokenHasCorrectSize(token); err != nil {
		return err
	}
	return nil
}

func tokenHasCorrectSize(token string) error {
	if len(token) != 34 {
		return errors.New(string(TokenIncorrectSize))
	}
	return nil
}

//ErrorMessage for error events within the pushbullet service
type ErrorMessage string

const (
	serviceURL = "https://api.pushbullet.com/v2/pushes"
	//Scheme is the scheme part of the service configuration URL
	Scheme = "pushbullet"
	//TokenIncorrectSize for the serviceURL
	TokenIncorrectSize ErrorMessage = "Token has incorrect size"
)
