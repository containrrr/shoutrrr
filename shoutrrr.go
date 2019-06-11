package shoutrrr

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/router"
	"log"
)

var routing = router.ServiceRouter{}

// SetLogger sets the logger that the services will use to write progress logs
func SetLogger(logger *log.Logger) {
	routing.SetLogger(logger)
}

// Send lets you send shoutrrr notifications using a supplied url and message
func Send(rawURL string, message string, params *map[string]string) error {
	if scheme, url, err := routing.ExtractServiceName(rawURL); err != nil {
		return err
	} else if plugin, err := routing.Locate(scheme); err != nil {
		return err
	} else {
		return plugin.Send(url, message, params)
	}
}

// Verify lets you verify that a configuration URL is valid and see what configuration it would map to
func Verify(rawURL string) error {

	routing := router.ServiceRouter{}

	svc, url, err := routing.ExtractServiceName(rawURL)
	if err != nil {
		return err
	}

	plugin, err := routing.Locate(svc)
	if err != nil {
		return err
	}

	config := plugin.GetConfig()
	if err := config.SetURL(url); err != nil {
		return err
	}

	configMap := format.GetConfigMap(config)
	for key, value := range configMap {
		fmt.Printf("%s: %s\n", key, value)
	}

	return nil
}