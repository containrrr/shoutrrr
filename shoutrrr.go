package shoutrrr

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/router"
	"log"
	"strings"
)

var routing = router.ServiceRouter{}

// SetLogger sets the logger that the services will use to write progress logs
func SetLogger(logger *log.Logger) {
	routing.SetLogger(logger)
}

// Send lets you send shoutrrr notifications using a supplied url and message
func Send(rawURL string, message string) error {
	service, err := routing.Locate(rawURL);
	if err != nil {
		return err
	}

	return service.Send(message, nil)
}

// Verify lets you verify that a configuration URL is valid and see what configuration it would map to
func Verify(rawURL string) error {

	config, err := routing.Parse(rawURL)
	if err != nil {
		return err
	}

	configMap, maxKeyLen := format.GetConfigMap(config)
	for key, value := range configMap {
		pad := strings.Repeat(" ", maxKeyLen -len(key))
		fmt.Printf("%s%s: %s\n", pad, key, value)
	}

	return nil
}